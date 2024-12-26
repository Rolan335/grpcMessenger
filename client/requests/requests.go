package requests

import (
	"context"

	"github.com/Rolan335/grpcMessenger/proto"
)

func InitSession(ctx context.Context, c proto.MessengerServiceClient, req *proto.InitSessionRequest) (*proto.InitSessionResponse, error) {
	resp, err := c.InitSession(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func CreateChat(ctx context.Context, c proto.MessengerServiceClient, req *proto.CreateChatRequest) (*proto.CreateChatResponse, error) {
	resp, err := c.CreateChat(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SendMessage(ctx context.Context, c proto.MessengerServiceClient, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	resp, err := c.SendMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetHistory(ctx context.Context, c proto.MessengerServiceClient, req *proto.GetHistoryRequest) (*proto.GetHistoryResponse, error) {
	resp, err := c.GetHistory(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetActiveChats(ctx context.Context, c proto.MessengerServiceClient, req *proto.GetActiveChatsRequest) (*proto.GetActiveChatsResponse, error) {
	resp, err := c.GetActiveChats(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
