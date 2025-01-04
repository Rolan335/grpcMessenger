// Implementation of storage interface for concurrent-safe inmemory storage
package storage

import (
	"sync"

	"github.com/Rolan335/grpcMessenger/server/proto"
	"github.com/Rolan335/grpcMessenger/server/internal/serviceErrors"
	lru "github.com/hashicorp/golang-lru"
)

// К каждому сообщению привязан id юзера, id сообщения и само сообщение
type Message struct {
	SessionUuid string `json:"session_uuid"`
	MessageUuid string `json:"message_uuid"`
	Text        string `json:"text"`
}

// Чат хранит в себе id юзера создавшего чат, могут ли другие юзеры писать в чат, время (в секундах) через сколько чат удалится,
// id чата и сообщения в формате lru
type Chat struct {
	SessionUuid string `json:"session_uuid"`
	Readonly    bool   `json:"readonly"`
	ttl         int
	ChatUuid    string     `json:"chat_uuid"`
	Messages    *lru.Cache `json:"messages"`
}

// mutex для конкурентой работы с мапой юзеров. MaxChatSize и MaxChats хранят максимальный размер чата и максимальное кол-во чатов соответственно.
// Все чаты хранятся в lru. Все юзеры в мапе для оптимизации поиска.
type InMemoryStorage struct {
	MaxChatSize int
	MaxChats    int
	Users       map[string]struct{}
	mu          *sync.RWMutex
	ChatsData   *lru.Cache
}

func NewInMemoryStorage(maxChatSize int, maxChats int) *InMemoryStorage {
	//Creating lru for chats storage
	lru, _ := lru.New(maxChats)
	//Initializing new inmemory storage
	return &InMemoryStorage{
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
		Users:       make(map[string]struct{}),
		mu:          &sync.RWMutex{},
		ChatsData:   lru,
	}
}

func (s *InMemoryStorage) AddSession(sessionUuid string) {
	//Add session to storage with mutex
	s.mu.Lock()
	s.Users[sessionUuid] = struct{}{}
	s.mu.Unlock()
}

func (s *InMemoryStorage) AddChat(sessionUuid string, ttl int, readOnly bool, chatUuid string) error {
	//Lock for reading to retrieve a data
	s.mu.RLock()
	_, ok := s.Users[sessionUuid]
	//if not found send error
	if !ok {
		return serviceErrors.ErrUserDoesNotExist
	}
	s.mu.RUnlock()

	//Creating new lru for chat to store messages.
	lru, _ := lru.New(s.MaxChatSize)

	//Creating new chat as a pointer to add messages directly
	newChat := &Chat{
		SessionUuid: sessionUuid,
		Readonly:    readOnly,
		ChatUuid:    chatUuid,
		Messages:    lru,
	}

	//Add new chat to lru in storage struct
	s.ChatsData.Add(chatUuid, newChat)
	return nil
}

func (s *InMemoryStorage) AddMessage(sessionUuid string, chatUuid string, messageUuid string, text string) error {
	//Making message
	newMessage := Message{
		SessionUuid: sessionUuid,
		MessageUuid: messageUuid,
		Text:        text,
	}

	//Trying to get chat from lru
	chat, ok := s.ChatsData.Get(chatUuid)

	//send error if not found
	if !ok {
		return serviceErrors.ErrChatNotFound
	}

	//type assert retrieved chat
	chatAsserted := chat.(*Chat)

	//check if we can send message to chat
	if chatAsserted.Readonly && chatAsserted.SessionUuid != sessionUuid {
		return serviceErrors.ErrProhibited
	}

	//add new message to chat
	chatAsserted.Messages.Add(messageUuid, newMessage)
	return nil
}

func (s *InMemoryStorage) GetHistory(chatUuid string) (history []*proto.ChatMessage, err error) {
	//get chat with provided chatUuid
	chat, ok := s.ChatsData.Get(chatUuid)

	//if not found, return error
	if !ok {
		return nil, serviceErrors.ErrChatNotFound
	}

	//type assert chat
	chatAsserted := chat.(*Chat)

	//getting keys from chat to iterate
	keys := chatAsserted.Messages.Keys()

	//slice for storing ChatMessages
	var msgArr []*proto.ChatMessage

	//iterating through retrieved keys and getting all Message struct with it
	for _, v := range keys {
		//get chat struct with key
		msg, _ := chatAsserted.Messages.Get(v)
		//type assert gotten message
		msgAsserted := msg.(Message)
		//creating struct for proto response and append it to slice
		msgArr = append(msgArr, &proto.ChatMessage{
			SessionUuid: msgAsserted.SessionUuid,
			MessageUuid: msgAsserted.MessageUuid,
			Text:        msgAsserted.Text,
		})
	}
	//returning completed history of messages
	return msgArr, nil
}

func (s *InMemoryStorage) DeleteChat(sessionUuid string, chatUuid string) error {
	//Check if chat is present
	chat, ok := s.ChatsData.Get(chatUuid)
	if !ok {
		return serviceErrors.ErrChatNotFound
	}

	//type Assertion
	chatAsserted := chat.(*Chat)
	//Deletion is can be made only by chat Creator
	if chatAsserted.SessionUuid != sessionUuid {
		return serviceErrors.ErrProhibited
	}
	//Delete chat
	s.ChatsData.Remove(chatUuid)
	return nil
}

func (s *InMemoryStorage) GetActiveChats() (chats []*proto.Chat) {
	//Get keys for all chats in lru
	chatKeys := s.ChatsData.Keys()

	//Range over keys
	for _, key := range chatKeys {

		//Get data from lru with key
		chat, _ := s.ChatsData.Get(key)

		//type assert
		chatAsserted := chat.(*Chat)

		//Create proto.Chat instance and append it to returned slice
		chats = append(chats, &proto.Chat{
			ChatUuid:    chatAsserted.ChatUuid,
			SessionUuid: chatAsserted.SessionUuid,
			Ttl:         int32(chatAsserted.ttl),
			ReadOnly:    chatAsserted.Readonly,
		})
	}
	return chats
}
