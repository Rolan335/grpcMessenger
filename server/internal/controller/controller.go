// Package controller provides implementation of Service routes with proceeding requests, sending data to storage, logging and sending response.
//
//nolint:wrapcheck
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
	m *messenger.Messenger
	proto.UnimplementedMessengerServiceServer
}

// Server creation. Initializing storage and returning service
func NewServer(config config.ServiceCfg, logger logger.Logger) Server {
	return Server{
		m: messenger.NewMessenger(config, logger),
	}
}

// Implementation of InitSession rpc
func (s Server) InitSession(_ context.Context, _ *proto.InitSessionRequest) (*proto.InitSessionResponse, error) {
	//Call initsession
	sessionUUID := s.m.InitSession()

	//Creating, sending response
	response := &proto.InitSessionResponse{SessionUuid: sessionUUID}
	return response, nil
}

// Implementation of CreateChat rpc
func (s Server) CreateChat(_ context.Context, r *proto.CreateChatRequest) (*proto.CreateChatResponse, error) {
	ChatUUID, err := s.m.CreateChat(r.GetSessionUuid(), int(r.GetTtl()), r.GetReadOnly())
	if err != nil {
		if errors.Is(err, messenger.ErrInvalidSessionUUID) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, messenger.ErrUserDoesNotExist) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Creating, sending response
	response := &proto.CreateChatResponse{ChatUuid: ChatUUID}
	return response, nil
}

// Implementation of SendMessage rpc
func (s Server) SendMessage(_ context.Context, r *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	err := s.m.SendMessage(r.GetSessionUuid(), r.GetChatUuid(), r.GetMessage())
	if err != nil {
		if errors.Is(err, messenger.ErrInvalidSessionUUID) || errors.Is(err, messenger.ErrInvalidChatUUID) {
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
func (s Server) GetHistory(_ context.Context, r *proto.GetHistoryRequest) (*proto.GetHistoryResponse, error) {
	messages, err := s.m.GetHistory(r.GetChatUuid())
	if err != nil {
		if errors.Is(err, messenger.ErrChatNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	//Create response
	history := make([]*proto.ChatMessage, 0, len(messages))
	for _, v := range messages {
		history = append(history, &proto.ChatMessage{
			SessionUuid: v.SessionUUID,
			MessageUuid: v.MessageUUID,
			Text:        v.Text,
		})
	}
	response := &proto.GetHistoryResponse{Messages: history}
	return response, nil
}

// implementation of GetActiveChats rpc
func (s Server) GetActiveChats(_ context.Context, _ *proto.GetActiveChatsRequest) (*proto.GetActiveChatsResponse, error) {
	chats := s.m.GetActiveChats()
	res := make([]*proto.Chat, 0, len(chats))
	for _, v := range chats {
		res = append(res, &proto.Chat{
			SessionUuid: v.SessionUUID,
			ChatUuid:    v.ChatUUID,
			Ttl:         int32(v.TTL),
			ReadOnly:    v.ReadOnly,
		})
	}
	response := &proto.GetActiveChatsResponse{Chats: res}
	return response, nil
}

func (s Server) HealthCheck(_ context.Context, _ *proto.HealthCheckRequest) (*proto.HealthCheckResponse, error) {
	return &proto.HealthCheckResponse{Status: "SERVING"}, nil
}
