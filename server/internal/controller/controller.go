// Package controller provides implementation of Service routes with proceeding requests, sending data to storage, logging and sending response.
package controller

import (
	"context"
	"errors"

	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/service/messenger"
	"github.com/Rolan335/grpcMessenger/server/pkg/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	m      *messenger.Messenger
	proto.UnimplementedMessengerServiceServer
}

// Server creation. Initializing storage and returning service
func NewServer(config config.ServiceCfg, logger logger.Logger) Server {
	return Server{
		m:      messenger.NewMessenger(config, logger),
	}
}

// Implementation of InitSession rpc
func (s Server) InitSession(ctx context.Context, r *proto.InitSessionRequest) (*proto.InitSessionResponse, error) {
	//Call initsession
	sessionUuid := s.m.InitSession()

	//Creating, sending response
	response := &proto.InitSessionResponse{SessionUuid: sessionUuid}
	return response, nil
}

// Implementation of CreateChat rpc
func (s Server) CreateChat(ctx context.Context, r *proto.CreateChatRequest) (*proto.CreateChatResponse, error) {
	ChatUuid, err := s.m.CreateChat(r.GetSessionUuid(), int(r.GetTtl()), r.GetReadOnly())
	if err != nil {
		if errors.Is(err, messenger.ErrInvalidSessionUuid) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, messenger.ErrUserDoesNotExist) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Creating, sending response
	response := &proto.CreateChatResponse{ChatUuid: ChatUuid}
	return response, nil
}

// Implementation of SendMessage rpc
func (s Server) SendMessage(ctx context.Context, r *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	err := s.m.SendMessage(r.GetSessionUuid(), r.GetChatUuid(), r.GetMessage())
	if err != nil {
		if errors.Is(err, messenger.ErrInvalidSessionUuid) || errors.Is(err, messenger.ErrInvalidChatUuid) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, messenger.ErrChatNotFound) || errors.Is(err, messenger.ErrUserDoesNotExist) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, messenger.ErrProhibited) {
			return nil, status.Error(codes.PermissionDenied, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Creating, sending response
	response := &proto.SendMessageResponse{}
	return response, nil
}

// Implementation of GetHistory rpc
func (s Server) GetHistory(ctx context.Context, r *proto.GetHistoryRequest) (*proto.GetHistoryResponse, error) {
	messages, err := s.m.GetHistory(r.GetChatUuid())
	if err != nil {
		if errors.Is(err, messenger.ErrChatNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Create response
	var history []*proto.ChatMessage
	for _, v := range messages {
		history = append(history, &proto.ChatMessage{
			SessionUuid: v.SessionUuid,
			MessageUuid: v.MessageUuid,
			Text:        v.Text,
		})
	}
	response := &proto.GetHistoryResponse{Messages: history}
	return response, nil
}

// implementation of GetActiveChats rpc
func (s Server) GetActiveChats(ctx context.Context, r *proto.GetActiveChatsRequest) (*proto.GetActiveChatsResponse, error) {
	chats := s.m.GetActiveChats()
	var res []*proto.Chat
	for _, v := range chats {
		res = append(res, &proto.Chat{
			SessionUuid: v.SessionUuid,
			ChatUuid:    v.ChatUuid,
			Ttl:         int32(v.Ttl),
			ReadOnly:    v.ReadOnly,
		})
	}
	response := &proto.GetActiveChatsResponse{Chats: res}
	return response, nil
}
