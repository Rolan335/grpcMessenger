package test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/app"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/service/messenger"
	"github.com/Rolan335/grpcMessenger/server/pkg/proto"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type noOpWriter struct{}

func (n *noOpWriter) Write(p []byte) (int, error) {
	return 0, nil
}

func parseConfig() config.ServiceCfg {
	err := godotenv.Load(".env.test")
	if err != nil {
		panic("cannot load .env file: " + err.Error())
	}

	maxChatSize, err := strconv.Atoi(os.Getenv("APP_MAXCHATSIZE"))
	if err != nil {
		panic("failed to parse maxChatSize .env:" + err.Error())
	}
	maxChats, err := strconv.Atoi(os.Getenv("APP_MAXCHATS"))
	if err != nil {
		panic("failed to parse maxChats .env:" + err.Error())
	}
	return config.MustConfigInit(
		os.Getenv("APP_ADDRESS"),
		os.Getenv("APP_PORTGRPC"),
		os.Getenv("APP_PORTHTTP"),
		os.Getenv("APP_ENV"),
		maxChatSize,
		maxChats,
		os.Getenv("APP_DB"),
	)
}

// test for all storages
func TestRepositories(t *testing.T) {
	testCases := []string{"inmemory", "redis", "postgres"}
	logger.Init("dev", &noOpWriter{})
	for _, v := range testCases {
		t.Run(v, func(t *testing.T) {
			serverConfig := parseConfig()
			serverConfig.StorageType = v
			server := app.NewServiceServer(serverConfig)
			defer server.GracefulStop()
			go server.MustStartGRPC()
			conn, err := grpc.NewClient(serverConfig.Address+serverConfig.PortGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				panic(err)
			}
			defer conn.Close()
			c := proto.NewMessengerServiceClient(conn)

			a := assert.New(t)
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
				a.EqualValues(serverConfig.MaxChatSize, len(resp.GetMessages()), "there should be no more than maxChatSize messages in chat")
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
				a.NoError(err, "c.CreateChat shouldn't return an error")
				time.Sleep(time.Second * 3)
				resp, _ := c.GetActiveChats(ctx, &proto.GetActiveChatsRequest{})
				a.Equal(len(resp.Chats), 1, "chat should be deleted after 2 seconds") // 1 len because we created chat previously
			})

			t.Run("Test readonly", func(t *testing.T) {
				newUser, _ := c.InitSession(ctx, &proto.InitSessionRequest{})
				_, err := c.SendMessage(ctx, &proto.SendMessageRequest{
					ChatUuid:    chatsCreated[0],
					SessionUuid: newUser.GetSessionUuid(),
					Message:     "hello",
				})
				a.ErrorIs(err, status.Error(codes.PermissionDenied, "prohibited. Only creator can send"), "should return proper mistake")
			})

			t.Run("Invalid uuid createChat", func(t *testing.T) {
				uuid := "siwroieqrw-214124-wwrwrr-2222"
				_, err := c.CreateChat(ctx, &proto.CreateChatRequest{
					SessionUuid: uuid,
					Ttl:         -1,
					ReadOnly:    false,
				})
				a.ErrorIs(err, status.Error(codes.InvalidArgument, messenger.ErrInvalidSessionUUID.Error()), "should return proper mistake")
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
				a.ErrorIs(err, status.Error(codes.InvalidArgument, messenger.ErrInvalidSessionUUID.Error()), "should return proper mistake")
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
					a.ErrorIs(err, status.Error(codes.NotFound, "chat not found"), "returned proper error")
				}
			})
		})
	}
}

// standart test with inmemory
func TestServer(t *testing.T) {
	logger.Init("dev", &noOpWriter{})
	serverConfig := parseConfig()
	fmt.Println(serverConfig)
	server := app.NewServiceServer(serverConfig)
	defer server.GracefulStop()
	go server.MustStartGRPC()
	conn, err := grpc.NewClient(serverConfig.Address+serverConfig.PortGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := proto.NewMessengerServiceClient(conn)

	a := assert.New(t)
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
		a.EqualValues(serverConfig.MaxChatSize, len(resp.GetMessages()), "there should be no more than maxChatSize messages in chat")
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
		a.NoError(err, "c.CreateChat shouldn't return an error")
		time.Sleep(time.Second * 3)
		resp, _ := c.GetActiveChats(ctx, &proto.GetActiveChatsRequest{})
		a.Equal(len(resp.Chats), 1, "chat should be deleted after 2 seconds") // 1 len because we created chat previously
	})

	t.Run("Test readonly", func(t *testing.T) {
		newUser, _ := c.InitSession(ctx, &proto.InitSessionRequest{})
		_, err := c.SendMessage(ctx, &proto.SendMessageRequest{
			ChatUuid:    chatsCreated[0],
			SessionUuid: newUser.GetSessionUuid(),
			Message:     "hello",
		})
		a.ErrorIs(err, status.Error(codes.PermissionDenied, "prohibited. Only creator can send"), "should return proper mistake")
	})

	t.Run("Invalid uuid createChat", func(t *testing.T) {
		uuid := "siwroieqrw-214124-wwrwrr-2222"
		_, err := c.CreateChat(ctx, &proto.CreateChatRequest{
			SessionUuid: uuid,
			Ttl:         -1,
			ReadOnly:    false,
		})
		a.ErrorIs(err, status.Error(codes.InvalidArgument, messenger.ErrInvalidSessionUUID.Error()), "should return proper mistake")
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
		a.ErrorIs(err, status.Error(codes.InvalidArgument, messenger.ErrInvalidSessionUUID.Error()), "should return proper mistake")
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
			a.ErrorIs(err, status.Error(codes.NotFound, "chat not found"), "returned proper error")
		}
	})
}
