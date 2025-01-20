package inmemory

// Implementation of storage interface for concurrent-safe inmemory storage

import (
	"sync"

	"github.com/Rolan335/grpcMessenger/server/internal/repository"

	lru "github.com/hashicorp/golang-lru"
)

// К каждому сообщению привязан id юзера, id сообщения и само сообщение
type Message struct {
	SessionUUID string
	MessageUUID string
	Text        string
}

// Чат хранит в себе id юзера создавшего чат, могут ли другие юзеры писать в чат, время (в секундах) через сколько чат удалится,
type Chat struct {
	SessionUUID string
	ReadOnly    bool
	TTL         int
	ChatUUID    string
	Messages    *lru.Cache
}

type User struct {
	SessionUUID string
}

// mutex для конкурентой работы с мапой юзеров. MaxChatSize и MaxChats хранят максимальный размер чата и максимальное кол-во чатов соответственно.
// Все чаты хранятся в lru. Все юзеры в мапе для оптимизации поиска.
type Storage struct {
	MaxChatSize int
	MaxChats    int
	mu          *sync.RWMutex
	ChatsData   *lru.Cache
	Users       map[User]struct{}
}

func NewStorage(maxChatSize int, maxChats int) *Storage {
	//Creating lru for chats storage
	lru, _ := lru.New(maxChats)
	//Initializing new inmemory storage
	return &Storage{
		MaxChatSize: maxChatSize,
		MaxChats:    maxChats,
		mu:          &sync.RWMutex{},
		ChatsData:   lru,
		Users:       make(map[User]struct{}),
	}
}

func (s *Storage) AddSession(sessionUUID string) {
	//Add session to storage with mutex
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Users[User{SessionUUID: sessionUUID}] = struct{}{}
}

func (s *Storage) AddChat(sessionUUID string, ttl int, readOnly bool, chatUUID string) error {
	//Lock for reading to retrieve a data
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.Users[User{SessionUUID: sessionUUID}]
	//if not found send error
	if !ok {
		return repository.ErrNotFound
	}

	//Creating new lru for chat to store messages.
	lru, _ := lru.New(s.MaxChatSize)

	//Creating new chat as a pointer to add messages directly
	newChat := &Chat{
		SessionUUID: sessionUUID,
		ReadOnly:    readOnly,
		ChatUUID:    chatUUID,
		TTL:         ttl,
		Messages:    lru,
	}

	//Add new chat to lru in storage struct
	s.ChatsData.Add(chatUUID, newChat)
	return nil
}

func (s *Storage) AddMessage(sessionUUID string, chatUUID string, messageUUID string, text string) error {

	//Making message
	newMessage := Message{
		SessionUUID: sessionUUID,
		MessageUUID: messageUUID,
		Text:        text,
	}

	//Trying to get chat from lru
	chat, ok := s.ChatsData.Get(chatUUID)
	//send error if not found
	if !ok {
		return repository.ErrNotFound
	}
	s.mu.RLock()
	_, ok = s.Users[User{SessionUUID: sessionUUID}]
	s.mu.RUnlock()
	//send error if not found
	if !ok {
		return repository.ErrUserDoesntExist
	}

	//type assert retrieved chat
	chatAsserted := chat.(*Chat)

	//check if we can send message to chat
	if chatAsserted.ReadOnly && chatAsserted.SessionUUID != sessionUUID {
		return repository.ErrProhibited
	}

	//add new message to chat
	chatAsserted.Messages.Add(messageUUID, newMessage)
	return nil
}

func (s *Storage) GetHistory(chatUUID string) (history []repository.Message, err error) {
	//get chat with provided chatUUID
	chat, ok := s.ChatsData.Get(chatUUID)

	//if not found, return error
	if !ok {
		return nil, repository.ErrNotFound
	}

	//type assert chat
	chatAsserted := chat.(*Chat)

	//getting keys from chat to iterate
	keys := chatAsserted.Messages.Keys()

	//slice for storing ChatMessages
	msgArr := make([]repository.Message, 0, len(keys))

	//iterating through retrieved keys and getting all Message struct with it
	for _, v := range keys {
		//get chat struct with key
		msg, _ := chatAsserted.Messages.Get(v)
		//type assert gotten message
		msgAsserted := msg.(Message)
		//creating struct for proto response and append it to slice
		msgArr = append(msgArr, repository.Message{
			SessionUUID: msgAsserted.SessionUUID,
			MessageUUID: msgAsserted.MessageUUID,
			Text:        msgAsserted.Text,
		})
	}
	//returning completed history of messages
	return msgArr, nil
}

func (s *Storage) DeleteChat(sessionUUID string, chatUUID string) error {
	//Check if chat is present
	chat, ok := s.ChatsData.Get(chatUUID)
	if !ok {
		return repository.ErrNotFound
	}

	//type Assertion
	chatAsserted := chat.(*Chat)
	//Deletion is can be made only by chat Creator
	if chatAsserted.SessionUUID != sessionUUID {
		return repository.ErrProhibited
	}
	//Delete chat
	s.ChatsData.Remove(chatUUID)
	return nil
}

func (s *Storage) GetActiveChats() (chats []repository.Chat) {
	//Get keys for all chats in lru
	chatKeys := s.ChatsData.Keys()

	//Range over keys
	for _, key := range chatKeys {

		//Get data from lru with key
		chat, _ := s.ChatsData.Get(key)

		//type assert
		chatAsserted := *chat.(*Chat)

		//Create Chat instance and append it to returned slice
		chats = append(chats, repository.Chat{
			SessionUUID: chatAsserted.SessionUUID,
			ChatUUID:    chatAsserted.ChatUUID,
			ReadOnly:    chatAsserted.ReadOnly,
			TTL:         chatAsserted.TTL,
		})
	}
	return chats
}
