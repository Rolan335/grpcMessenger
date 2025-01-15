package messenger

import (
	"errors"
	"os"
	"strconv"

	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/repository"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/inmemory"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/postgres"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/redis"

	"github.com/google/uuid"
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
		postgresConfig := postgres.PostgresConfig{
			Host:       os.Getenv("POSTGRES_HOST"),
			User:       os.Getenv("POSTGRES_USER"),
			Password:   os.Getenv("POSTGRES_PASSWORD"),
			Dbname:     os.Getenv("POSTGRES_DBNAME"),
			Port:       port,
			FreshStart: freshstart,
		}
		db = postgres.NewPostgresStorage(postgresConfig, config.MaxChats, config.MaxChatSize)
	case "redis":
		freshstart, err := strconv.ParseBool(os.Getenv("REDIS_FLUSH"))
		if err != nil {
			panic("failed to parse .env REDIS_FLUSH: " + err.Error())
		}
		redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			panic("failed to parse .env REDIS_DB: " + err.Error())
		}
		redisConfig := redis.RedisConfig{
			Addr:     os.Getenv("REDIS_ADDRESS"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       redisDb,
			FlushAll: freshstart,
		}
		db = redis.NewRedisStorage(redisConfig, config.MaxChatSize, config.MaxChats)
	default:
		db = inmemory.NewInMemoryStorage(config.MaxChatSize, config.MaxChats)
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

func (m *Messenger) CreateChat(sessionUuid string, ttl int, readOnly bool) (string, error) {
	//If invalid sessionUuid provided - request cannot be completed, return invalidUuid error.
	if _, err := uuid.Parse(sessionUuid); err != nil {
		return "", ErrInvalidSessionUuid
	}

	//Creating uuid for chat
	id, _ := uuid.NewRandom()

	//Add new chat to server storage
	err := m.storage.AddChat(sessionUuid, ttl, readOnly, id.String())
	if err != nil {
		//If nonExistent session-uuid provided - returning error
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrUserDoesNotExist
		}
		return "", err
	}

	//If ttl is set, chat will be deleted after time elapsed
	if ttl > 0 {
		DeleteAfter(ttl, sessionUuid, id.String(), m.storage, m.logger)
	}
	return id.String(), nil
}

func (m *Messenger) SendMessage(sessionUuid string, chatUuid string, message string) error {
	//If invalid sessionUuid or chatUuid provided  - request cannot be completed, return invalidargs error.
	if _, err := uuid.Parse(sessionUuid); err != nil {
		return ErrInvalidSessionUuid
	}
	if _, err := uuid.Parse(chatUuid); err != nil {
		return ErrInvalidChatUuid
	}

	//Creating uuid for message
	id, _ := uuid.NewRandom()
	//Adding new message to storage and if failed - returns error
	err := m.storage.AddMessage(sessionUuid, chatUuid, id.String(), message)
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
		return err
	}

	return nil
}

func (m *Messenger) GetHistory(chatUuid string) ([]repository.Message, error) {
	// if invalid chatUuid provided - request cannot be completed, return error
	if _, err := uuid.Parse(chatUuid); err != nil {
		return nil, ErrInvalidChatUuid
	}

	//get history from storage with chatUuid provided
	history, err := m.storage.GetHistory(chatUuid)
	if err != nil {
		//If chat not found - returning error
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrChatNotFound
		}
		return nil, err
	}
	return history, nil
}

func (m *Messenger) GetActiveChats() []repository.Chat {
	chats := m.storage.GetActiveChats()
	return chats
}
