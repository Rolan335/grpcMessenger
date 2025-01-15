// Implementation of storage interface for concurrent-safe inmemory storage
package inmemory

import (
	"sync"

	"github.com/Rolan335/grpcMessenger/server/internal/repository"

	lru "github.com/hashicorp/golang-lru"
)

// К каждому сообщению привязан id юзера, id сообщения и само сообщение
type InmemoryMessage struct {
	SessionUuid string
	MessageUuid string
	Text        string
}

// Чат хранит в себе id юзера создавшего чат, могут ли другие юзеры писать в чат, время (в секундах) через сколько чат удалится,
type InmemoryChat struct {
	SessionUuid string
	ReadOnly    bool
	Ttl         int
	ChatUuid    string
	Messages    *lru.Cache
}

type InmemoryUser struct {
	SessionUuid string
}

// mutex для конкурентой работы с мапой юзеров. MaxChatSize и MaxChats хранят максимальный размер чата и максимальное кол-во чатов соответственно.
// Все чаты хранятся в lru. Все юзеры в мапе для оптимизации поиска.
type InMemoryStorage struct {
	MaxChatSize int
	MaxChats    int
	mu          *sync.RWMutex
	ChatsData   *lru.Cache
	Users       map[InmemoryUser]struct{}
}

func NewInMemoryStorage(maxChatSize int, maxChats int) *InMemoryStorage {
	//Creating lru for chats storage
	lru, _ := lru.New(maxChats)
	//Initializing new inmemory storage
	return &InMemoryStorage{
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
		mu:          &sync.RWMutex{},
		ChatsData:   lru,
		Users:       make(map[InmemoryUser]struct{}),
	}
}

func (s *InMemoryStorage) AddSession(sessionUuid string) {
	//Add session to storage with mutex
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Users[InmemoryUser{SessionUuid: sessionUuid}] = struct{}{}
}

func (s *InMemoryStorage) AddChat(sessionUuid string, ttl int, readOnly bool, chatUuid string) error {
	//Lock for reading to retrieve a data
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.Users[InmemoryUser{SessionUuid: sessionUuid}]
	//if not found send error
	if !ok {
		return repository.ErrNotFound
	}

	//Creating new lru for chat to store messages.
	lru, _ := lru.New(s.MaxChatSize)

	//Creating new chat as a pointer to add messages directly
	newChat := &InmemoryChat{
		SessionUuid: sessionUuid,
		ReadOnly:    readOnly,
		ChatUuid:    chatUuid,
		Messages:    lru,
	}

	//Add new chat to lru in storage struct
	s.ChatsData.Add(chatUuid, newChat)
	return nil
}

func (s *InMemoryStorage) AddMessage(sessionUuid string, chatUuid string, messageUuid string, text string) error {

	//Making message
	newMessage := InmemoryMessage{
		SessionUuid: sessionUuid,
		MessageUuid: messageUuid,
		Text:        text,
	}

	//Trying to get chat from lru
	chat, ok := s.ChatsData.Get(chatUuid)
	//send error if not found
	if !ok {
		return repository.ErrNotFound
	}
	s.mu.RLock()
	_, ok = s.Users[InmemoryUser{SessionUuid: sessionUuid}]
	s.mu.RUnlock()
	//send error if not found
	if !ok {
		return repository.ErrUserDoesntExist
	}

	//type assert retrieved chat
	chatAsserted := chat.(*InmemoryChat)

	//check if we can send message to chat
	if chatAsserted.ReadOnly && chatAsserted.SessionUuid != sessionUuid {
		return repository.ErrProhibited
	}

	//add new message to chat
	chatAsserted.Messages.Add(messageUuid, newMessage)
	return nil
}

func (s *InMemoryStorage) GetHistory(chatUuid string) (history []repository.Message, err error) {
	//get chat with provided chatUuid
	chat, ok := s.ChatsData.Get(chatUuid)

	//if not found, return error
	if !ok {
		return nil, repository.ErrNotFound
	}

	//type assert chat
	chatAsserted := chat.(*InmemoryChat)

	//getting keys from chat to iterate
	keys := chatAsserted.Messages.Keys()

	//slice for storing ChatMessages
	var msgArr []repository.Message

	//iterating through retrieved keys and getting all Message struct with it
	for _, v := range keys {
		//get chat struct with key
		msg, _ := chatAsserted.Messages.Get(v)
		//type assert gotten message
		msgAsserted := msg.(InmemoryMessage)
		//creating struct for proto response and append it to slice
		msgArr = append(msgArr, repository.Message{
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
		return repository.ErrNotFound
	}

	//type Assertion
	chatAsserted := chat.(*InmemoryChat)
	//Deletion is can be made only by chat Creator
	if chatAsserted.SessionUuid != sessionUuid {
		return repository.ErrProhibited
	}
	//Delete chat
	s.ChatsData.Remove(chatUuid)
	return nil
}

func (s *InMemoryStorage) GetActiveChats() (chats []repository.Chat) {
	//Get keys for all chats in lru
	chatKeys := s.ChatsData.Keys()

	//Range over keys
	for _, key := range chatKeys {

		//Get data from lru with key
		chat, _ := s.ChatsData.Get(key)

		//type assert
		chatAsserted := *chat.(*InmemoryChat)

		//Create Chat instance and append it to returned slice
		chats = append(chats, repository.Chat{
			SessionUuid: chatAsserted.SessionUuid,
			ChatUuid:    chatAsserted.ChatUuid,
			ReadOnly:    chatAsserted.ReadOnly,
			Ttl:         chatAsserted.Ttl,
		})
	}
	return chats
}
