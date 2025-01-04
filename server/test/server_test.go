package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Rolan335/grpcMessenger/proto"
	"github.com/Rolan335/grpcMessenger/server/internal/app/serverinit"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/serviceErrors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func clientInit() (proto.MessengerServiceClient, *grpc.ClientConn) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
	}
	c := proto.NewMessengerServiceClient(conn)
	return c, conn
}

func startClientAndServer(maxChatSize int, maxChats int) (proto.MessengerServiceClient, *grpc.ClientConn) {
	go serverinit.Start(config.ServiceCfg{PortGrpc: ":50051", Env: "dev", MaxChatSize: maxChatSize, MaxChats: maxChats})
	return clientInit()
}

func TestServer(t *testing.T) {
	maxChats := 3
	maxChatSize := 5
	a := assert.New(t)
	c, conn := startClientAndServer(maxChatSize, maxChats)
	defer conn.Close()
	ctx := context.Background()
	var clientUuid string
	chatsCreated := make([]string, 0)

	t.Run("InitSession", func(t *testing.T) {
		resp, err := c.InitSession(ctx, &proto.InitSessionRequest{})
		if a.NoError(err, "c.InitSession shouldn't return an error") {
			_, err := uuid.Parse(resp.GetSessionUuid())
			a.NoError(err, "returned uuid parsed successfully")
		}

		clientUuid = resp.GetSessionUuid()
	})
	t.Run("CreateChat", func(t *testing.T) {
		resp, err := c.CreateChat(ctx, &proto.CreateChatRequest{
			SessionUuid: clientUuid,
			ReadOnly:    true,
			Ttl:         -1,
		})
		if a.NoError(err, "c.CreateChat shouldn't return an error") {
			_, err := uuid.Parse(resp.GetChatUuid())
			a.NoError(err, "returned uuid parsed successfully")
			chatsCreated = append(chatsCreated, resp.GetChatUuid())
		}
	})

	messages := []*proto.SendMessageRequest{
		{
			ChatUuid:    chatsCreated[0],
			SessionUuid: clientUuid,
			Message:     "1",
		},
		{
			ChatUuid:    chatsCreated[0],
			SessionUuid: clientUuid,
			Message:     "2",
		},
		{
			ChatUuid:    chatsCreated[0],
			SessionUuid: clientUuid,
			Message:     "3",
		},
		{
			ChatUuid:    chatsCreated[0],
			SessionUuid: clientUuid,
			Message:     "4",
		},
		{
			ChatUuid:    chatsCreated[0],
			SessionUuid: clientUuid,
			Message:     "5",
		},
		{
			ChatUuid:    chatsCreated[0],
			SessionUuid: clientUuid,
			Message:     "6",
		},
	}
	t.Run("SendMessage", func(t *testing.T) {
		for _, v := range messages {
			_, err := c.SendMessage(ctx, v)
			a.NoError(err, "c.SendMessage shouldn't return an error")
		}
	})

	t.Run("GetHistory", func(t *testing.T) {
		resp, err := c.GetHistory(ctx, &proto.GetHistoryRequest{
			ChatUuid: chatsCreated[0],
		})
		a.NoError(err, "c.GetHistory shouldn't return an error")
		a.EqualValues(maxChatSize, len(resp.GetMessages()), "there should be no more than maxChatSize messages in chat")
		for i, v := range resp.GetMessages() {
			_, err := uuid.Parse(v.GetMessageUuid())
			a.NoError(err, "returned uuid parsed successfully")
			//6 messages sent, 5 should stay, so checking from sent messages as i+1
			a.Equal(messages[i+1].GetMessage(), v.GetText(), "text in messages should be equal")
			a.Equal(messages[i+1].GetSessionUuid(), v.GetSessionUuid(), "Session uuid should be equal")
		}
	})

	t.Run("GetActiveChats", func(t *testing.T) {
		resp, err := c.GetActiveChats(ctx, &proto.GetActiveChatsRequest{})
		a.NoError(err, "c.GetActiveChats shouldn't return an error")
		chats := resp.GetChats()
		a.Equal(len(chatsCreated), len(chats))
		a.EqualValues(chats[0].GetChatUuid(), chatsCreated[0], "should return same uuid for chat created")
	})

	t.Run("Test ttl", func(t *testing.T) {
		_, err := c.CreateChat(ctx, &proto.CreateChatRequest{
			SessionUuid: clientUuid,
			ReadOnly:    false,
			Ttl:         2,
		})
		a.NoError(err,"c.CreateChat shouldn't return an error")
		time.Sleep(time.Second * 3)
		resp, _ := c.GetActiveChats(ctx, &proto.GetActiveChatsRequest{})
		a.Equal(len(resp.Chats), 1, "chat should be deleted after 2 seconds") // 1 len because we created chat previously
		fmt.Println(resp.Chats)
	})

	t.Run("Test readonly", func(t *testing.T) {
		newUser, _ := c.InitSession(ctx, &proto.InitSessionRequest{})
		_, err := c.SendMessage(ctx, &proto.SendMessageRequest{
			ChatUuid:    chatsCreated[0],
			SessionUuid: newUser.GetSessionUuid(),
			Message:     "hello",
		})
		a.ErrorIs(err, serviceErrors.ErrProhibited, "should return proper mistake")
	})

	t.Run("Invalid uuid createChat", func(t *testing.T) {
		uuid := "siwroieqrw-214124-wwrwrr-2222"
		_, err := c.CreateChat(ctx, &proto.CreateChatRequest{
			SessionUuid: uuid,
			Ttl:         -1,
			ReadOnly:    false,
		})
		a.ErrorIs(err, serviceErrors.ErrInvalidUuid, "should return proper mistake")
		resp, _ := c.GetActiveChats(ctx, &proto.GetActiveChatsRequest{})
		a.Equal(1, len(resp.GetChats()), "chat should not be created")
	})

	t.Run("Invalid uuid sendMessage", func(t *testing.T) {
		uuid := "isapodiasfpo-qwfpowqf-ir124fg-332ii"
		_, err := c.SendMessage(ctx, &proto.SendMessageRequest{
			ChatUuid:    uuid,
			SessionUuid: uuid,
			Message:     "hello",
		})
		a.ErrorIs(err, serviceErrors.ErrInvalidUuid, "should return proper mistake")
	})

	t.Run("ChatNotFound", func(t *testing.T) {
		resp, err := c.CreateChat(ctx, &proto.CreateChatRequest{
			SessionUuid: clientUuid,
			Ttl:         2,
			ReadOnly:    false,
		})
		a.NoError(err, "no error returned")
		time.Sleep(3 * time.Second)
		{
			_, err := c.SendMessage(ctx, &proto.SendMessageRequest{
				SessionUuid: clientUuid,
				Message:     "hello",
				ChatUuid:    resp.GetChatUuid(),
			})
			a.ErrorIs(err, serviceErrors.ErrChatNotFound, "returned proper error")
		}
	})
}
