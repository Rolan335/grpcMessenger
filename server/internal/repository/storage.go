package repository

//Storage interface with methods that we need to implement so our storage will be able to work in service
//Storage is responsible for checking violations, returning errors if violated (ex. trying to send message to readonly chat)
type Storage interface {
	AddSession(sessionUuid string)
	AddChat(sessionUuid string, ttl int, readOnly bool, chatUuid string) error
	DeleteChat(sessionUuid string, chatUuid string) error
	AddMessage(sessionUuid string, chatUuid string, messageUuid string, message string) error
	GetHistory(chatUuid string) (history []Message, err error)
	GetActiveChats() (chats []Chat)
}

// К каждому сообщению привязан id юзера, id сообщения и само сообщение
type Message struct {
	SessionUuid string
	MessageUuid string
	Text        string
}

// Чат хранит в себе id юзера создавшего чат, могут ли другие юзеры писать в чат, время (в секундах) через сколько чат удалится,
type Chat struct {
	SessionUuid string
	ReadOnly    bool
	Ttl         int
	ChatUuid    string
}

type User struct {
	SessionUuid string
}
