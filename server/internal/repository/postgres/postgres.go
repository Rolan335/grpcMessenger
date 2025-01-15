package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresConfig struct {
	Host       string
	User       string
	Password   string
	Dbname     string
	Port       int
	FreshStart bool
}

type PostgresStorage struct {
	MaxChats    int
	MaxChatSize int
	Db          *gorm.DB
}

type Message struct {
	MessageUuid string `gorm:"type:uuid;primaryKey;uniqueIndex"`
	SessionUuid string `gorm:"type:uuid"`
	ChatUuid    string `gorm:"type:uuid"`
	Text        string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type Chat struct {
	ChatUuid    string `gorm:"type:uuid;primaryKey;uniqueIndex"`
	SessionUuid string `gorm:"type:uuid"`
	ReadOnly    bool
	Ttl         int
	Message     []Message `gorm:"foreignKey:ChatUuid;references:ChatUuid"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type User struct {
	SessionUuid string `gorm:"type:uuid;primaryKey;uniqueIndex"`
	Chat        []Chat `gorm:"foreignKey:SessionUuid;references:SessionUuid"`
}

var db *gorm.DB

func NewPostgresStorage(cfg PostgresConfig, maxChats int, maxChatSize int) *PostgresStorage {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.Host, cfg.User, cfg.Password, cfg.Dbname, cfg.Port)
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if cfg.FreshStart {
		err := deleteAllTables(db)
		if err != nil {
			fmt.Println(err)
		}
	}
	if err != nil {
		panic("can't connect to postgres: " + err.Error())
	}
	err = db.AutoMigrate(&User{}, &Message{}, &Chat{})
	if err != nil {
		panic("can't do migrations: " + err.Error())
	}
	return &PostgresStorage{
		MaxChats:    maxChats,
		MaxChatSize: maxChatSize,
		Db:          db,
	}
}

func (p *PostgresStorage) AddSession(sessionUuid string) {
	p.Db.Create(&User{SessionUuid: sessionUuid})
}

func (p *PostgresStorage) AddChat(sessionUuid string, ttl int, readOnly bool, chatUuid string) error {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.First(&User{}, "session_uuid = ?", sessionUuid).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return repository.ErrNotFound
	}
	if err := tx.Create(&Chat{
		ChatUuid:    chatUuid,
		SessionUuid: sessionUuid,
		Ttl:         ttl,
		ReadOnly:    readOnly,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	var chatsCount int64
	if err := tx.Model(&Chat{}).Count(&chatsCount).Error; err != nil {
		tx.Rollback()
		return err
	}

	if chatsCount > int64(p.MaxChats) {
		p.DeleteLeastChats(tx, chatsCount)
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func deleteAllTables(db *gorm.DB) error {
	// Get the list of all tables
	var tables []string
	if err := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables).Error; err != nil {
		return fmt.Errorf("failed to get tables: %w", err)
	}

	// Drop each table
	for _, table := range tables {
		if err := db.Migrator().DropTable(table); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
		fmt.Printf("Dropped table: %s\n", table)
	}

	return nil
}

func (p *PostgresStorage) DeleteLeastChats(tx *gorm.DB, chatsCount int64) {
	var chats []Chat
	tx.
		Order("created_at ASC").
		Limit(int(chatsCount - int64(p.MaxChatSize))).
		Find(&chats)
	for _, v := range chats {
		tx.Delete(&Message{}, "chat_uuid = ?", v.ChatUuid)
		tx.Delete(&v)
	}
}

func (p *PostgresStorage) DeleteChat(sessionUuid string, chatUuid string) error {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var chatToDelete Chat
	if err := tx.First(&chatToDelete, "chat_uuid = ?", chatUuid).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return repository.ErrNotFound
	}
	if chatToDelete.SessionUuid != sessionUuid {
		tx.Rollback()
		return repository.ErrProhibited
	}
	tx.Delete(&Message{}, "chat_uuid = ?", chatToDelete.ChatUuid)
	tx.Delete(&chatToDelete)

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
func (p *PostgresStorage) AddMessage(sessionUuid string, chatUuid string, messageUuid string, message string) error {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var chat Chat
	if err := tx.First(&chat, "chat_uuid = ?", chatUuid).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return repository.ErrNotFound
	}
	if err := tx.First(&User{}, "session_uuid = ?", sessionUuid).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return repository.ErrUserDoesntExist
	}
	if chat.ReadOnly && chat.SessionUuid != sessionUuid {
		tx.Rollback()
		return repository.ErrProhibited
	}
	if err := tx.Create(&Message{
		MessageUuid: messageUuid,
		SessionUuid: sessionUuid,
		ChatUuid:    chatUuid,
		Text:        message,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	var messageCount int64
	if err := tx.Model(&Message{}).Where("chat_uuid = ?", chatUuid).Count(&messageCount).Error; err != nil {
		tx.Rollback()
		return err
	}
	if messageCount > int64(p.MaxChatSize) {
		err := p.DeleteLeastMsg(tx, messageCount, chatUuid)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (p *PostgresStorage) DeleteLeastMsg(tx *gorm.DB, messageCount int64, chatUuid string) error {
	var messages []Message
	if err := tx.Where("chat_uuid = ?", chatUuid).
		Order("created_at ASC").
		Limit(int(messageCount - int64(p.MaxChatSize))).
		Find(&messages).Error; err != nil {
		return err
	}
	for _, v := range messages {
		if err := tx.Delete(&v).Error; err != nil {
			return err
		}
	}
	return nil
}

func (p *PostgresStorage) GetHistory(chatUuid string) (history []repository.Message, err error) {
	if err := p.Db.First(&Chat{}, "chat_uuid = ?", chatUuid).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	var messages []Message
	p.Db.Where("chat_uuid = ?", chatUuid).Order("created_at ASC").Find(&messages)
	for _, v := range messages {
		history = append(history, repository.Message{
			SessionUuid: v.SessionUuid,
			MessageUuid: v.MessageUuid,
			Text:        v.Text,
		})
	}
	return
}
func (p *PostgresStorage) GetActiveChats() (chats []repository.Chat) {
	var chatsToFind []Chat
	p.Db.Find(&chatsToFind)
	for _, v := range chatsToFind {
		chats = append(chats, repository.Chat{
			SessionUuid: v.SessionUuid,
			ChatUuid:    v.ChatUuid,
			Ttl:         v.Ttl,
			ReadOnly:    v.ReadOnly,
		})
	}
	return
}

func GracefulStop() {
	if db == nil {
		return
	}
	sqlDb, _ := db.DB()
	sqlDb.Close()
	fmt.Println("postgres stopped successfully")
}
