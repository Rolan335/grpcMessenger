package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/repository"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/entities"

	"github.com/redis/go-redis/v9"
)

var (
	keyUser            = "users"        // users - set session_UUID
	keyActiveChats     = "active_chats" // active_chats - list chat_UUID
	keyPrefixChat      = "chat:"        // chat:{chat_UUID} - list Chat{...}
	keyPostfixMessages = ":messages"    // chat:{chat_UUID}:messages - list message{...}
)

type Message struct {
	MessageUUID string `json:"message_UUID"`
	SessionUUID string `json:"session_UUID"`
	Text        string `json:"text"`
}

type Chat struct {
	ChatUUID    string `json:"chat_UUID"`
	SessionUUID string `json:"session_UUID"`
	ReadOnly    bool   `json:"read_only"`
	TTL         int    `json:"ttl"`
}

type User struct {
	SessionUUID string `json:"session_UUID"`
}

type Storage struct {
	MaxChatSize int
	MaxChats    int
	client      *redis.Client
}

type Config struct {
	Addr     string
	Password string
	DB       int
	FlushAll bool
}

var rdb *redis.Client

func NewStorage(cfg Config, maxChatSize int, maxChats int) *Storage {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic("failed to connect to redis: " + err.Error())
	}
	if cfg.FlushAll {
		rdb.FlushAll(context.Background())
	}
	return &Storage{
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
		client:      rdb,
	}
}

func (r *Storage) AddSession(sessionUUID string) {
	r.client.SAdd(context.Background(), keyUser, sessionUUID)
}

func (r *Storage) AddChat(sessionUUID string, ttl int, readOnly bool, chatUUID string) error {
	ctx := context.Background()
	isPresent := r.client.SIsMember(ctx, keyUser, sessionUUID).Val()
	if !isPresent {
		return repository.ErrNotFound
	}
	r.client.RPush(ctx, keyActiveChats, chatUUID)
	//Удаление чата если больше maxChats (LRU)
	excessChats := r.client.LRange(ctx, keyActiveChats, 0, -int64(r.MaxChats)-1).Val()
	for range excessChats {
		chatDeleted := r.client.LPop(ctx, keyActiveChats).Val()
		r.client.Del(ctx, fmt.Sprintf("%s%s", keyPrefixChat, chatDeleted))
		r.client.Del(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chatDeleted, keyPostfixMessages))
	}

	chatJSON, _ := json.Marshal(Chat{
		SessionUUID: sessionUUID,
		ChatUUID:    chatUUID,
		TTL:         ttl,
		ReadOnly:    readOnly,
	})
	r.client.Set(ctx, fmt.Sprintf("%s%s", keyPrefixChat, chatUUID), chatJSON, 0)
	return nil
}
func (r *Storage) DeleteChat(sessionUUID string, chatUUID string) error {
	ctx := context.Background()
	chat, err := getChatFromKey(ctx, r, chatUUID)
	if errors.Is(err, redis.Nil) {
		return repository.ErrNotFound
	}
	if chat.SessionUUID != sessionUUID {
		return repository.ErrProhibited
	}

	r.client.LRem(ctx, keyActiveChats, 0, chat.ChatUUID)
	r.client.Del(ctx, fmt.Sprintf("%s%s", keyPrefixChat, chat.ChatUUID))
	r.client.Del(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chat.ChatUUID, keyPostfixMessages))

	return nil
}

func (r *Storage) AddMessage(sessionUUID string, chatUUID string, messageUUID string, message string) error {
	ctx := context.Background()
	chat, err := getChatFromKey(ctx, r, chatUUID)
	if errors.Is(err, redis.Nil) {
		return repository.ErrNotFound
	}
	if !r.client.SIsMember(ctx, keyUser, sessionUUID).Val() {
		return repository.ErrUserDoesntExist
	}
	if chat.ReadOnly && chat.SessionUUID != sessionUUID {
		return repository.ErrProhibited
	}

	messageJSON, _ := json.Marshal(Message{
		MessageUUID: messageUUID,
		SessionUUID: sessionUUID,
		Text:        message,
	})

	r.client.RPush(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chat.ChatUUID, keyPostfixMessages), messageJSON)
	//Удаление Сообщения если больше maxChatSize (LRU)
	r.client.LTrim(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chat.ChatUUID, keyPostfixMessages), -int64(r.MaxChatSize), -1)

	return nil
}

func getChatFromKey(ctx context.Context, r *Storage, chatUUID string) (Chat, error) {
	chat, err := r.client.Get(ctx, fmt.Sprintf("%s%s", keyPrefixChat, chatUUID)).Result()
	if err != nil {
		return Chat{}, fmt.Errorf("redis: %w", err)
	}
	var chatUnmarshalled Chat
	err = json.Unmarshal([]byte(chat), &chatUnmarshalled)
	if err != nil {
		return Chat{}, fmt.Errorf("redis: %w", err)
	}
	return chatUnmarshalled, nil
}

func (r *Storage) GetHistory(chatUUID string) (history []entities.Message, err error) {
	ctx := context.Background()
	messages, err := r.client.LRange(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chatUUID, keyPostfixMessages), 0, -1).Result()
	if errors.Is(err, redis.Nil) {
		return nil, repository.ErrNotFound
	}
	for _, v := range messages {
		var message Message
		err := json.Unmarshal([]byte(v), &message)
		if err != nil {
			return nil, fmt.Errorf("redis: %w", err)
		}
		history = append(history, entities.Message{
			SessionUUID: message.SessionUUID,
			MessageUUID: message.MessageUUID,
			Text:        message.Text,
		})
	}
	return
}
func (r *Storage) GetActiveChats() (chats []entities.Chat) {
	ctx := context.Background()
	chatsUUID := r.client.LRange(ctx, keyActiveChats, 0, -1).Val()
	for _, v := range chatsUUID {
		chat, _ := getChatFromKey(ctx, r, v)
		chats = append(chats, entities.Chat{
			SessionUUID: chat.SessionUUID,
			ChatUUID:    chat.ChatUUID,
			ReadOnly:    chat.ReadOnly,
			TTL:         chat.TTL,
		})
	}
	return
}

func GracefulStop() {
	if rdb == nil {
		return
	}
	rdb.Close()
	fmt.Println("redis closed successfully")
}

// nolint
func Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return rdb.Ping(ctx).Err()
}
