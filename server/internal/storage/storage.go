package storage

import "github.com/Rolan335/grpcMessenger/proto"


//Storage interface with methods that we need to implement so our storage will be able to work in service
//Storage is responsible for checking violations, returning errors if violated (ex. trying to send message to readonly chat)
type Storage interface {
	AddSession(sessionUuid string)
	AddChat(sessionUuid string, ttl int, readOnly bool, chatUuid string) error
	DeleteChat(sessionUuid string, chatUuid string) error
	AddMessage(sessionUuid string, chatUuid string, messageUuid string, message string) error
	GetHistory(chatUuid string) (history []*proto.ChatMessage, err error)
	GetActiveChats() (chats []*proto.Chat)
}
