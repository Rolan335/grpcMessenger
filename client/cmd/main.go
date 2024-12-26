package main

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"github.com/Rolan335/grpcMessenger/client/requests"
	"github.com/Rolan335/grpcMessenger/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var UUID string

func main() {
	address := flag.String("addr", "localhost:50051", "enter address of the service [127.0.0.1:50051]")
	flag.Parse()
	conn, err := grpc.NewClient(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	ctx := context.TODO()
	c := proto.NewMessengerServiceClient(conn)
	resp, err := requests.InitSession(ctx, c, &proto.InitSessionRequest{})
	if err != nil {
		panic("server connection error: " + err.Error())
	}
	UUID = resp.GetSessionUuid()
	fmt.Println("Your session: ", UUID)
	for {
		var action string
		fmt.Scan(&action)
		switch strings.ToLower(action) {
		case "createchat":
			var (
				ttl      int
				readOnly bool
			)
			fmt.Println("Enter [ttl] [readonly]")
			fmt.Scan(&ttl, &readOnly)
			if ttl > 0 {

			}
			request := &proto.CreateChatRequest{Chat: &proto.Chat{SessionUuid: UUID, Ttl: int32(ttl), ReadOnly: readOnly}}
			_, err := requests.CreateChat(ctx, c, request)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "sendmessage":
			var (
				chatUuid string
				message  string
			)
			fmt.Println("Enter chatUuid")
			fmt.Scan(&chatUuid)
			fmt.Println("Enter your message")
			fmt.Scan(&message)
			request := &proto.SendMessageRequest{SessionUuid: UUID, Message: message, ChatUuid: chatUuid}
			_, err := requests.SendMessage(ctx, c, request)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("message sent")
		case "gethistory":
			var (
				chatUuid string
			)
			fmt.Println("Enter chatUuid")
			fmt.Scan(&chatUuid)
			request := &proto.GetHistoryRequest{ChatUuid: chatUuid}
			resp, err := requests.GetHistory(ctx, c, request)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(resp.String())
		case "getactivechats":
			request := &proto.GetActiveChatsRequest{}
			resp, err := requests.GetActiveChats(ctx, c, request)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(resp.String())
		default:
			fmt.Println("Enter command")
		}
	}
}
