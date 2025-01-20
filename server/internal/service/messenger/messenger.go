package messenger

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/google/uuid"

	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/repository"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/inmemory"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/postgres"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/redis"
)

type Messenger struct {
	storage repository.Storage
	logger  logger.Logger
}

func NewMessenger(config config.ServiceCfg, logger logger.Logger) *Messenger {
	var db repository.Storage
	switch config.StorageType {
	case "postgres":
		port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
		if err != nil {
			panic("failed to parse .env POSTGRES_PORT: " + err.Error())
		}
		freshstart, err := strconv.ParseBool(os.Getenv("POSTGRES_FLUSH"))
		if err != nil {
			panic("failed to parse .env POSTGRES_FLUSH: " + err.Error())
		}
		postgresConfig := postgres.Config{
			Host:       os.Getenv("POSTGRES_HOST"),
			User:       os.Getenv("POSTGRES_USER"),
			Password:   os.Getenv("POSTGRES_PASSWORD"),
			Dbname:     os.Getenv("POSTGRES_DBNAME"),
			Port:       port,
			FreshStart: freshstart,
		}
		db = postgres.NewStorage(postgresConfig, config.MaxChats, config.MaxChatSize)
	case "redis":
		freshstart, err := strconv.ParseBool(os.Getenv("REDIS_FLUSH"))
		if err != nil {
			panic("failed to parse .env REDIS_FLUSH: " + err.Error())
		}
		redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			panic("failed to parse .env REDIS_DB: " + err.Error())
		}
		redisConfig := redis.Config{
			Addr:     os.Getenv("REDIS_ADDRESS"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDb,
			FlushAll: freshstart,
		}
		db = redis.NewStorage(redisConfig, config.MaxChatSize, config.MaxChats)
	default:
		db = inmemory.NewStorage(config.MaxChatSize, config.MaxChats)
	}
	return &Messenger{
		storage: db,
		logger:  logger,
	}
}

func (m *Messenger) InitSession() string {
	//Creating uuid for user
	id, _ := uuid.NewRandom()

	//Add Session to server storage
	m.storage.AddSession(id.String())
	return id.String()
}

func (m *Messenger) CreateChat(sessionUUID string, ttl int, readOnly bool) (string, error) {
	//If invalid sessionUUID provided - request cannot be completed, return invalidUUID error.
	if _, err := uuid.Parse(sessionUUID); err != nil {
		return "", ErrInvalidSessionUUID
	}

	//Creating uuid for chat
	id, _ := uuid.NewRandom()

	//Add new chat to server storage
	err := m.storage.AddChat(sessionUUID, ttl, readOnly, id.String())
	if err != nil {
		//If nonExistent session-UUID provided - returning error
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrUserDoesNotExist
		}
		return "", fmt.Errorf("messenger: %w",err)
	}

	//If ttl is set, chat will be deleted after time elapsed
	if ttl > 0 {
		DeleteAfter(ttl, sessionUUID, id.String(), m.storage, m.logger)
	}
	return id.String(), nil
}

func (m *Messenger) SendMessage(sessionUUID string, chatUUID string, message string) error {
	//If invalid sessionUUID or chatUUID provided  - request cannot be completed, return invalidargs error.
	if _, err := uuid.Parse(sessionUUID); err != nil {
		return ErrInvalidSessionUUID
	}
	if _, err := uuid.Parse(chatUUID); err != nil {
		return ErrInvalidChatUUID
	}

	//Creating uuid for message
	id, _ := uuid.NewRandom()
	//Adding new message to storage and if failed - returns error
	err := m.storage.AddMessage(sessionUUID, chatUUID, id.String(), message)
	if err != nil {
		//Check if chat not found - send error
		if errors.Is(err, repository.ErrNotFound) {
			return ErrChatNotFound
		}
		if errors.Is(err, repository.ErrUserDoesntExist) {
			return ErrUserDoesNotExist
		}
		//Check if chat is readonly
		if errors.Is(err, repository.ErrProhibited) {
			return ErrProhibited
		}
		return fmt.Errorf("messenger: %w",err)
	}

	return nil
}

func (m *Messenger) GetHistory(chatUUID string) ([]repository.Message, error) {
	// if invalid chatUUID provided - request cannot be completed, return error
	if _, err := uuid.Parse(chatUUID); err != nil {
		return nil, ErrInvalidChatUUID
	}

	//get history from storage with chatUUID provided
	history, err := m.storage.GetHistory(chatUUID)
	if err != nil {
		//If chat not found - returning error
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, fmt.Errorf("messenger: %w",err)
	}
	return history, nil
}

func (m *Messenger) GetActiveChats() []repository.Chat {
	chats := m.storage.GetActiveChats()
	return chats
}
