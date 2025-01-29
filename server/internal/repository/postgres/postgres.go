package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/repository"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/entities"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" //import std drivers *sql.Db for goose migrations
	"github.com/pressly/goose/v3"
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
	Timeout     int
	Db          *pgxpool.Pool
}

type Message struct {
	MessageUUID uuid.UUID `pg:"message_uuid"`
	SessionUUID uuid.UUID `pg:"session_uuid"`
	ChatUUID    string    `pg:"chat_uuid"`
	Text        string    `pg:"text"`
	CreatedAt   time.Time `pg:"created_at"`
}

type Chat struct {
	ChatUUID    uuid.UUID `pg:"chat_uuid"`
	SessionUUID uuid.UUID `pg:"session_uuid"`
	ReadOnly    bool      `pg:"read_only"`
	TTL         int       `pg:"ttl"`
	CreatedAt   time.Time `pg:"created_at"`
}

type User struct {
	SessionUUID uuid.UUID `pg:"session_uuid"`
}

var conn *pgxpool.Pool

func NewStorage(cfg Config, maxChats int, maxChatSize int, timeout int, migrationsPath string) *Storage {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname)
	var err error
	conn, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic("can't connect to postgres: " + err.Error())
	}

	if cfg.FreshStart {
		if err := deleteAllTables(); err != nil {
			fmt.Println(err)
		}
	}

	MustDoMigrations(connStr, migrationsPath)

	return &Storage{
		MaxChats:    maxChats,
		MaxChatSize: maxChatSize,
		Timeout:     timeout,
		Db:          conn,
	}
}

func deleteAllTables() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Получение списка всех таблиц
	rows, err := conn.Query(ctx, "SELECT tablename FROM pg_tables WHERE schemaname = 'public'")
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer rows.Close()

	// Удаление таблиц
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("postgres: %w", err)
		}
		_, err := conn.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE;", tableName))
		if err != nil {
			return fmt.Errorf("postgres: %w", err)
		}
	}

	if rows.Err() != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	return nil
}

func MustDoMigrations(connStr string, migrationsPath string) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		panic("failed to do sql.Conn: " + err.Error())
	}
	defer db.Close()
	err = goose.Up(db, migrationsPath)
	if err != nil {
		panic("failed to do migrations: " + err.Error())
	}
}

func (p *Storage) AddSession(sessionUUID string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Timeout)*time.Second)
	defer cancel()
	p.Db.Exec(ctx, "INSERT INTO users (session_uuid) VALUES ($1)", sessionUUID)
}

func (p *Storage) AddChat(sessionUUID string, ttl int, readOnly bool, chatUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Timeout)*time.Second)
	defer cancel()
	tx, err := p.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
		}
	}()

	//Check if user exists
	if err := tx.QueryRow(ctx, "SELECT session_uuid FROM users WHERE session_uuid = $1 LIMIT 1", sessionUUID).Scan(nil); err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return repository.ErrUserDoesntExist
		}
		return fmt.Errorf("postgres: %w", err)
	}

	//Add new record to chats table
	if _, err := tx.Exec(ctx, "INSERT INTO chats (chat_uuid, session_uuid, read_only, ttl) VALUES ($1, $2, $3, $4)", chatUUID, sessionUUID, readOnly, ttl); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("postgres: %w", err)
	}

	//Check if any chats should be deleted
	var chatsCount int
	if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM chats").Scan(&chatsCount); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("postgres: %w", err)
	}

	if chatsCount > p.MaxChats {
		err := p.DeleteLeastChats(ctx, tx, chatsCount)
		if err != nil {
			return fmt.Errorf("postgres: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("postgres: %w", err)
	}

	return nil
}

func (p *Storage) DeleteLeastChats(ctx context.Context, tx pgx.Tx, chatsCount int) error {
	query := `
	DELETE FROM chats 
	USING (
        SELECT chat_uuid FROM chats
        ORDER BY created_at ASC
        LIMIT $1
	) AS to_delete
	WHERE chats.chat_uuid = to_delete.chat_uuid;
	`
	if _, err := tx.Exec(ctx, query, chatsCount-p.MaxChats); err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	return nil
}

func (p *Storage) DeleteChat(sessionUUID string, chatUUID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Timeout)*time.Second)
	defer cancel()

	tx, err := p.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
		}
	}()

	var recievedSessionUUID uuid.UUID
	if err := tx.QueryRow(ctx, "SELECT session_uuid FROM chats WHERE chat_uuid = $1 LIMIT 1", chatUUID).Scan(&recievedSessionUUID); err != nil {
		tx.Rollback(ctx)
		return repository.ErrUserDoesntExist
	}
	if parsedSessionUUID, _ := uuid.Parse(sessionUUID); parsedSessionUUID != recievedSessionUUID {
		tx.Rollback(ctx)
		return repository.ErrProhibited
	}
	if _, err := tx.Exec(ctx, "DELETE FROM chats WHERE chat_uuid = $1", chatUUID); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("postgres: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("postgres: %w", err)
	}

	return nil
}
func (p *Storage) AddMessage(sessionUUID string, chatUUID string, messageUUID string, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Timeout)*time.Second)
	defer cancel()
	tx, err := p.Db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("postgres: %w", err)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback(ctx)
		}
	}()

	//Check if user exists
	var userCount int
	if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM users WHERE session_uuid = $1", sessionUUID).Scan(&userCount); err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return repository.ErrUserDoesntExist
		}
		return fmt.Errorf("postgres: %w", err)
	}

	var chat Chat
	//проверить существует ли чат
	if err := tx.QueryRow(ctx, "SELECT chat_uuid, session_uuid, read_only, ttl FROM chats WHERE chat_uuid = $1 LIMIT 1", chatUUID).
		Scan(&chat.ChatUUID, &chat.SessionUUID, &chat.ReadOnly, &chat.TTL); err != nil {
		tx.Rollback(ctx)
		if err == pgx.ErrNoRows {
			return repository.ErrNotFound
		}
		return fmt.Errorf("postgres: %w", err)
	}

	//Проверить ридонли ли чат
	if chat.ReadOnly && sessionUUID != chat.SessionUUID.String() {
		tx.Rollback(ctx)
		return repository.ErrProhibited
	}

	//Добавить запись в чат
	if _, err := tx.Exec(ctx, "INSERT INTO messages (message_uuid, session_uuid, chat_uuid, text) VALUES ($1, $2, $3, $4)",
		messageUUID, sessionUUID, chatUUID, message,
	); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("postgres: %w", err)
	}

	//Добавить проверку логики lru
	var messageCount int
	if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM messages WHERE chat_uuid = $1", chatUUID).Scan(&messageCount); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("postgres: %w", err)
	}

	if messageCount > p.MaxChatSize {
		err := p.DeleteLeastMsg(ctx, tx, messageCount, chatUUID)
		if err != nil {
			return fmt.Errorf("postgres: %w", err)
		}
	}

	//Коммит
	if err := tx.Commit(ctx); err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("postgres: %w", err)
	}
	return nil
}

func (p *Storage) DeleteLeastMsg(ctx context.Context, tx pgx.Tx, messageCount int, chatUUID string) error {
	query := `
	DELETE FROM messages
	USING (
		SELECT message_uuid FROM messages
		WHERE chat_uuid = $1
		ORDER BY created_at ASC
		LIMIT $2
	) AS to_delete
	WHERE messages.message_uuid = to_delete.message_uuid;
	`
	if _, err := tx.Exec(ctx, query, chatUUID, messageCount-p.MaxChatSize); err != nil {
		fmt.Println(messageCount - p.MaxChatSize)
		return fmt.Errorf("postgres: %w", err)
	}
	return nil
}

func (p *Storage) GetHistory(chatUUID string) (history []entities.Message, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Timeout)*time.Second)
	defer cancel()
	//Check if chat exists
	if err := p.Db.QueryRow(ctx, "SELECT chat_uuid FROM chats WHERE chat_uuid = $1 LIMIT 1", chatUUID).Scan(nil); err != nil {
		if err == pgx.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: %w", err)
	}
	rows, err := p.Db.Query(ctx, "SELECT session_uuid, message_uuid, text FROM messages WHERE chat_uuid = $1", chatUUID)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.SessionUUID, &message.MessageUUID, &message.Text); err != nil {
			return nil, fmt.Errorf("postgres: %w", err)
		}
		history = append(history, entities.Message{
			SessionUUID: message.SessionUUID.String(),
			MessageUUID: message.MessageUUID.String(),
			Text:        message.Text,
		})
	}
	return
}
func (p *Storage) GetActiveChats() (chats []entities.Chat) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(p.Timeout)*time.Second)
	defer cancel()
	rows, err := p.Db.Query(ctx, "SELECT session_uuid, read_only, ttl, chat_uuid FROM chats")
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.SessionUUID, &chat.ReadOnly, &chat.TTL, &chat.ChatUUID); err != nil {
			return nil
		}
		chats = append(chats, entities.Chat{
			SessionUUID: chat.SessionUUID.String(),
			ReadOnly:    chat.ReadOnly,
			TTL:         chat.TTL,
			ChatUUID:    chat.ChatUUID.String(),
		})
	}
	return
}

func GracefulStop() {
	if conn == nil {
		return
	}
	conn.Close()
}

// nolint
func Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return conn.Ping(ctx)
}
