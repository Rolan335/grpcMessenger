package messenger

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Rolan335/grpcMessenger/server/internal/repository"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/entities"
)

// Storage interface with methods that we need to implement so our storage will be able to work in service
// Storage is responsible for checking violations, returning errors if violated (ex. trying to send message to readonly chat)
type Storage interface {
	AddSession(sessionUUID string)
	AddChat(sessionUUID string, ttl int, readOnly bool, chatUUID string) error
	DeleteChat(sessionUUID string, chatUUID string) error
	AddMessage(sessionUUID string, chatUUID string, messageUUID string, message string) error
	GetHistory(chatUUID string) (history []entities.Message, err error)
	GetActiveChats() (chats []entities.Chat)
}

type Messenger struct {
	storage Storage
}

func NewMessenger(storage Storage) *Messenger {
	return &Messenger{
		storage: storage,
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
		return "", fmt.Errorf("messenger: %w", err)
	}

	//If ttl is set, chat will be deleted after time elapsed
	if ttl > 0 {
		DeleteAfter(ttl, sessionUUID, id.String(), m.storage)
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
		return fmt.Errorf("messenger: %w", err)
	}

	return nil
}

func (m *Messenger) GetHistory(chatUUID string) ([]entities.Message, error) {
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
		return nil, fmt.Errorf("messenger: %w", err)
	}
	return history, nil
}

func (m *Messenger) GetActiveChats() []entities.Chat {
	chats := m.storage.GetActiveChats()
	return chats
}
