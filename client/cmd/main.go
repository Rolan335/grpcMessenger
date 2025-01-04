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
	if err != nil {
		panic("cannot start client " + err.Error())
	}
	defer conn.Close()
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
		_, err := fmt.Scan(&action)
		if err != nil {
			fmt.Println(err)
            continue
		}
		switch strings.ToLower(action) {
		case "createchat":
			var (
				ttl      int
				readOnly bool
			)
			fmt.Println("Enter [ttl] [readonly]")
			_, err := fmt.Scan(&ttl, &readOnly)
			if err != nil {
				fmt.Println(err)
                continue
			}
			request := &proto.CreateChatRequest{SessionUuid: UUID, Ttl: int32(ttl), ReadOnly: readOnly}
			resp, err := requests.CreateChat(ctx, c, request)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(resp.GetChatUuid())
		case "sendmessage":
			var (
				chatUuid string
				message  string
			)
			fmt.Println("Enter chatUuid")
			_, err := fmt.Scan(&chatUuid)
			if err != nil {
				fmt.Println(err)
                continue
			}
			fmt.Println("Enter your message")
			_, err = fmt.Scan(&message)
			if err != nil {
				fmt.Println(err)
                continue
			}
			request := &proto.SendMessageRequest{SessionUuid: UUID, Message: message, ChatUuid: chatUuid}
			_, err = requests.SendMessage(ctx, c, request)
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
			_, err := fmt.Scan(&chatUuid)
			if err != nil {
				fmt.Println(err)
                continue
			}
			request := &proto.GetHistoryRequest{ChatUuid: chatUuid}
			resp, err := requests.GetHistory(ctx, c, request)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(resp.GetMessages())
		case "getactivechats":
			request := &proto.GetActiveChatsRequest{}
			resp, err := requests.GetActiveChats(ctx, c, request)
			if err != nil {
				fmt.Println(err)
				continue
			}
			chats := resp.GetChats()
			for _, v := range chats {
				fmt.Printf("chatUuid: %v, readonly: %v, creatorUuid: %v\n", v.ChatUuid, v.ReadOnly, v.SessionUuid)
			}
		default:
			fmt.Println("Enter command")
		}
	}
}
