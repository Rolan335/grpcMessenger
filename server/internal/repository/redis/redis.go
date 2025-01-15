package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Rolan335/grpcMessenger/server/internal/repository"

	"github.com/redis/go-redis/v9"
)

var (
	keyUser            = "users"        // users - set session_uuid
	keyActiveChats     = "active_chats" // active_chats - list chat_uuid
	keyPrefixChat      = "chat:"        // chat:{chat_uuid} - list RedisChat{...}
	keyPostfixMessages = ":messages"    // chat:{chat_uuid}:messages - list message{...}
)

type RedisMessage struct {
	MessageUuid string `json:"message_uuid"`
	SessionUuid string `json:"session_uuid"`
	Text        string `json:"text"`
}

type RedisChat struct {
	ChatUuid    string `json:"chat_uuid"`
	SessionUuid string `json:"session_uuid"`
	ReadOnly    bool   `json:"read_only"`
	Ttl         int    `json:"ttl"`
}

type RedisUser struct {
	SessionUuid string `json:"session_uuid"`
}

type RedisStorage struct {
	MaxChatSize int
	MaxChats    int
	client      *redis.Client
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
	FlushAll bool
}

var rdb *redis.Client

func NewRedisStorage(cfg RedisConfig, maxChatSize int, maxChats int) *RedisStorage {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       0,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		panic("failed to connect to redis: " + err.Error())
	}
	rdb.FlushAll(context.Background())
	return &RedisStorage{
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
		client:      rdb,
	}
}

func (r *RedisStorage) AddSession(sessionUuid string) {
	r.client.SAdd(context.Background(), keyUser, sessionUuid)
}

func (r *RedisStorage) AddChat(sessionUuid string, ttl int, readOnly bool, chatUuid string) error {
	ctx := context.Background()
	isPresent := r.client.SIsMember(ctx, keyUser, sessionUuid).Val()
	if !isPresent {
		return repository.ErrNotFound
	}
	r.client.RPush(ctx, keyActiveChats, chatUuid)
	//Удаление чата если больше maxChats (LRU)
	excessChats := r.client.LRange(ctx, keyActiveChats, 0, -int64(r.MaxChats)-1).Val()
	for range excessChats {
		chatDeleted := r.client.LPop(ctx, keyActiveChats).Val()
		r.client.Del(ctx, fmt.Sprintf("%s%s", keyPrefixChat, chatDeleted))
		r.client.Del(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chatDeleted, keyPostfixMessages))
	}

	chatJSON, _ := json.Marshal(RedisChat{
		SessionUuid: sessionUuid,
		ChatUuid:    chatUuid,
		Ttl:         ttl,
		ReadOnly:    readOnly,
	})
	r.client.Set(ctx, fmt.Sprintf("%s%s", keyPrefixChat, chatUuid), chatJSON, 0)
	return nil
}
func (r *RedisStorage) DeleteChat(sessionUuid string, chatUuid string) error {
	ctx := context.Background()
	chat, err := getChatFromKey(r, ctx, chatUuid)
	if errors.Is(err, redis.Nil) {
		return repository.ErrNotFound
	}
	if chat.SessionUuid != sessionUuid {
		return repository.ErrProhibited
	}

	r.client.LRem(ctx, keyActiveChats, 0, chat.ChatUuid)
	r.client.Del(ctx, fmt.Sprintf("%s%s", keyPrefixChat, chat.ChatUuid))
	r.client.Del(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chat.ChatUuid, keyPostfixMessages))

	return nil
}

func (r *RedisStorage) AddMessage(sessionUuid string, chatUuid string, messageUuid string, message string) error {
	ctx := context.Background()
	chat, err := getChatFromKey(r, ctx, chatUuid)
	if errors.Is(err, redis.Nil) {
		return repository.ErrNotFound
	}
	if !r.client.SIsMember(ctx, keyUser, sessionUuid).Val() {
		return repository.ErrUserDoesntExist
	}
	if chat.ReadOnly && chat.SessionUuid != sessionUuid {
		return repository.ErrProhibited
	}

	messageJSON, _ := json.Marshal(RedisMessage{
		MessageUuid: messageUuid,
		SessionUuid: sessionUuid,
		Text:        message,
	})

	r.client.RPush(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chat.ChatUuid, keyPostfixMessages), messageJSON)
	//Удаление Сообщения если больше maxChatSize (LRU)
	r.client.LTrim(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chat.ChatUuid, keyPostfixMessages), -int64(r.MaxChatSize), -1)

	return nil
}

func getChatFromKey(r *RedisStorage, ctx context.Context, chatUuid string) (RedisChat, error) {
	chat, err := r.client.Get(ctx, fmt.Sprintf("%s%s", keyPrefixChat, chatUuid)).Result()
	if err != nil {
		return RedisChat{}, err
	}
	var chatUnmarshalled RedisChat
	err = json.Unmarshal([]byte(chat), &chatUnmarshalled)
	if err != nil{
		return RedisChat{}, err
	}
	return chatUnmarshalled, nil
}

func (r *RedisStorage) GetHistory(chatUuid string) (history []repository.Message, err error) {
	ctx := context.Background()
	messages, err := r.client.LRange(ctx, fmt.Sprintf("%s%s%s", keyPrefixChat, chatUuid, keyPostfixMessages), 0, -1).Result()
	if errors.Is(err, redis.Nil) {
		return nil, repository.ErrNotFound
	}
	for _, v := range messages {
		var message RedisMessage
		err := json.Unmarshal([]byte(v), &message)
		if err != nil{
			return nil, err
		}
		history = append(history, repository.Message{
			SessionUuid: message.SessionUuid,
			MessageUuid: message.MessageUuid,
			Text:        message.Text,
		})
	}
	return
}
func (r *RedisStorage) GetActiveChats() (chats []repository.Chat) {
	ctx := context.Background()
	chats_uuid := r.client.LRange(ctx, keyActiveChats, 0, -1).Val()
	for _, v := range chats_uuid {
		chat, _ := getChatFromKey(r, ctx, v)
		chats = append(chats, repository.Chat{
			SessionUuid: chat.SessionUuid,
			ChatUuid:    chat.ChatUuid,
			ReadOnly:    chat.ReadOnly,
			Ttl:         chat.Ttl,
		})
	}
	return
}

func GracefulStop(){
	if rdb == nil{
		return
	}
	rdb.Close()
	fmt.Println("redis closed successfully")
}