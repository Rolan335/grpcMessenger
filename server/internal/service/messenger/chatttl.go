package messenger

import (
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/repository"
)

// ttl is in seconds. Starting goroutine that will invoke DeleteChat method when time elapsed.
func DeleteAfter(ttl int, sessionUuid string, chatUuid string, storage repository.Storage, logger logger.Logger) {
	go func() {
		<-time.After(time.Duration(ttl) * time.Second)
		err := storage.DeleteChat(sessionUuid, chatUuid)
		//when chat deleted - log it.
		logger.LogChatDelete(sessionUuid, chatUuid, err)
	}()
}
