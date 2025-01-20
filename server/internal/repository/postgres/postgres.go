package postgres

import (
	"errors"
	"fmt"
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	Host       string
	User       string
	Password   string
	Dbname     string
	Port       int
	FreshStart bool
}

type Storage struct {
	MaxChats    int
	MaxChatSize int
	Db          *gorm.DB
}

type Message struct {
	MessageUUID string `gorm:"type:UUID;primaryKey;uniqueIndex"`
	SessionUUID string `gorm:"type:UUID"`
	ChatUUID    string `gorm:"type:UUID;index:idx_chat_uuid_created_at"`
	Text        string
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_chat_uuid_created_at"`
}

type Chat struct {
	ChatUUID    string `gorm:"type:UUID;primaryKey;uniqueIndex"`
	SessionUUID string `gorm:"type:UUID"`
	ReadOnly    bool
	TTL         int
	Message     []Message `gorm:"foreignKey:ChatUUID;references:ChatUUID"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type User struct {
	SessionUUID string `gorm:"type:UUID;primaryKey;uniqueIndex"`
	Chat        []Chat `gorm:"foreignKey:SessionUUID;references:SessionUUID"`
}

var db *gorm.DB

func NewStorage(cfg Config, maxChats int, maxChatSize int) *Storage {
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
	return &Storage{
		MaxChats:    maxChats,
		MaxChatSize: maxChatSize,
		Db:          db,
	}
}

func (p *Storage) AddSession(sessionUUID string) {
	p.Db.Create(&User{SessionUUID: sessionUUID})
}

func (p *Storage) AddChat(sessionUUID string, ttl int, readOnly bool, chatUUID string) error {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.First(&User{}, "session_UUID = ?", sessionUUID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return repository.ErrNotFound
	}
	if err := tx.Create(&Chat{
		ChatUUID:    chatUUID,
		SessionUUID: sessionUUID,
		TTL:         ttl,
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

func (p *Storage) DeleteLeastChats(tx *gorm.DB, chatsCount int64) {
	var chats []Chat
	tx.
		Order("created_at ASC").
		Limit(int(chatsCount - int64(p.MaxChatSize))).
		Find(&chats)
	for _, v := range chats {
		tx.Delete(&Message{}, "chat_UUID = ?", v.ChatUUID)
		tx.Delete(&v)
	}
}

func (p *Storage) DeleteChat(sessionUUID string, chatUUID string) error {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var chatToDelete Chat
	if err := tx.First(&chatToDelete, "chat_UUID = ?", chatUUID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return repository.ErrNotFound
	}
	if chatToDelete.SessionUUID != sessionUUID {
		tx.Rollback()
		return repository.ErrProhibited
	}
	tx.Delete(&Message{}, "chat_UUID = ?", chatToDelete.ChatUUID)
	tx.Delete(&chatToDelete)

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
func (p *Storage) AddMessage(sessionUUID string, chatUUID string, messageUUID string, message string) error {
	tx := p.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var chat Chat
	if err := tx.First(&chat, "chat_UUID = ?", chatUUID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return repository.ErrNotFound
	}
	if err := tx.First(&User{}, "session_UUID = ?", sessionUUID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return repository.ErrUserDoesntExist
	}
	if chat.ReadOnly && chat.SessionUUID != sessionUUID {
		tx.Rollback()
		return repository.ErrProhibited
	}
	if err := tx.Create(&Message{
		MessageUUID: messageUUID,
		SessionUUID: sessionUUID,
		ChatUUID:    chatUUID,
		Text:        message,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}
	var messageCount int64
	if err := tx.Model(&Message{}).Where("chat_UUID = ?", chatUUID).Count(&messageCount).Error; err != nil {
		tx.Rollback()
		return err
	}
	if messageCount > int64(p.MaxChatSize) {
		err := p.DeleteLeastMsg(tx, messageCount, chatUUID)
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

func (p *Storage) DeleteLeastMsg(tx *gorm.DB, messageCount int64, chatUUID string) error {
	var messages []Message
	if err := tx.Where("chat_UUID = ?", chatUUID).
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

func (p *Storage) GetHistory(chatUUID string) (history []repository.Message, err error) {
	if err := p.Db.First(&Chat{}, "chat_UUID = ?", chatUUID).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	var messages []Message
	p.Db.Select("session_uuid, message_uuid, text").
		Where("chat_UUID = ?", chatUUID).
		Order("created_at ASC").
		Find(&messages)
	for _, v := range messages {
		history = append(history, repository.Message{
			SessionUUID: v.SessionUUID,
			MessageUUID: v.MessageUUID,
			Text:        v.Text,
		})
	}
	return
}
func (p *Storage) GetActiveChats() (chats []repository.Chat) {
	var chatsToFind []Chat
	p.Db.Find(&chatsToFind)
	for _, v := range chatsToFind {
		chats = append(chats, repository.Chat{
			SessionUUID: v.SessionUUID,
			ChatUUID:    v.ChatUUID,
			TTL:         v.TTL,
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

// nolint
func Ping() error {
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDb.Ping()
}
