package entities

// К каждому сообщению привязан id юзера, id сообщения и само сообщение
type Message struct {
	SessionUUID string `json:"session_uuid,omitempty"`
	MessageUUID string `json:"message_uuid,omitempty"`
	Text        string `json:"text,omitempty"`
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
