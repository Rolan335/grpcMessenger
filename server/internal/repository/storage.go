package repository

//Storage interface with methods that we need to implement so our storage will be able to work in service
//Storage is responsible for checking violations, returning errors if violated (ex. trying to send message to readonly chat)
type Storage interface {
	AddSession(sessionUUID string)
	AddChat(sessionUUID string, ttl int, readOnly bool, chatUUID string) error
	DeleteChat(sessionUUID string, chatUUID string) error
	AddMessage(sessionUUID string, chatUUID string, messageUUID string, message string) error
	GetHistory(chatUUID string) (history []Message, err error)
	GetActiveChats() (chats []Chat)
}

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
}

type User struct {
	SessionUUID string
}
