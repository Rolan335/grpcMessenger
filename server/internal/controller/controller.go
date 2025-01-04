// Package controller provides implementation of Service routes with proceeding requests, sending data to storage, logging and sending response.
package controller

import (
	"context"
	"errors"

	"github.com/Rolan335/grpcMessenger/proto"
	"github.com/Rolan335/grpcMessenger/server/internal/chatttl"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/serviceErrors"
	"github.com/Rolan335/grpcMessenger/server/internal/storage"
	"github.com/Rolan335/grpcMessenger/server/internal/util/checkuuid"

	"github.com/google/uuid"
)

// Variable for storage struct that we are using
var Storage storage.Storage

type Server struct {
	proto.UnimplementedMessengerServiceServer
}

// Server creation. Initializing storage and returning service
func NewServer(config config.ServiceCfg) Server {
	Storage = storage.NewInMemoryStorage(config.MaxChatSize, config.MaxChats)
	return Server{proto.UnimplementedMessengerServiceServer{}}
}

// Implementation of InitSession rpc
func (s Server) InitSession(ctx context.Context, r *proto.InitSessionRequest) (*proto.InitSessionResponse, error) {
	//Creating uuid for user
	id, _ := uuid.NewRandom()

	//Add Session to server storage
	Storage.AddSession(id.String())

	//Creating, logging and sending response
	response := &proto.InitSessionResponse{SessionUuid: id.String()}
	logger.LogRequest(ctx, "InitSession", r.String(), response.String())
	return response, nil
}

// Implementation of CreateChat rpc
func (s Server) CreateChat(ctx context.Context, r *proto.CreateChatRequest) (*proto.CreateChatResponse, error) {
	//If invalid sessionUuid provided - request cannot be completed, return invalidUuid error.
	if !checkuuid.IsParsed(r.GetSessionUuid()) {
		logger.LogRequestWithError(ctx, "CreateChat", r.String(), serviceErrors.ErrInvalidUuid)
		return nil, serviceErrors.ErrInvalidUuid
	}

	//Creating uuid for chat
	id, _ := uuid.NewRandom()

	//Add new chat to server storage
	err := Storage.AddChat(r.GetSessionUuid(), int(r.Ttl), r.ReadOnly, id.String())

	//If nonExistent session-uuid provided - logging and returning error
	if errors.Is(err, serviceErrors.ErrUserDoesNotExist) {
		logger.LogRequestWithError(ctx, "CreateChat", r.String(), serviceErrors.ErrUserDoesNotExist)
		return nil, serviceErrors.ErrUserDoesNotExist
	}

	//If ttl is set, chat will be deleted after time elapsed
	if r.Ttl > 0 {
		chatttl.DeleteAfter(int(r.Ttl), r.SessionUuid, id.String(), Storage)
	}

	//Creating, logging and sending response
	response := &proto.CreateChatResponse{ChatUuid: id.String()}
	logger.LogRequest(ctx, "CreateChat", r.String(), response.String())
	return response, nil
}

// Implementation of SendMessage rpc
func (s Server) SendMessage(ctx context.Context, r *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	//If invalid sessionUuid or chatUuid provided  - request cannot be completed, return invalidargs error.
	if !checkuuid.IsParsed(r.GetSessionUuid(), r.GetChatUuid()) {
		logger.LogRequestWithError(ctx, "SendMessage", r.String(), serviceErrors.ErrInvalidUuid)
		return nil, serviceErrors.ErrInvalidUuid
	}

	//Creating uuid for message
	id, _ := uuid.NewRandom()
	//Adding new message to storage and if failed - returns error
	err := Storage.AddMessage(r.SessionUuid, r.ChatUuid, id.String(), r.Message)
	//Check if chat not found - send error and log it
	if errors.Is(err, serviceErrors.ErrChatNotFound) {
		logger.LogRequestWithError(ctx, "SendMessage", r.String(), serviceErrors.ErrChatNotFound)
		return nil, serviceErrors.ErrChatNotFound
	}

	//Check if chat is readonly
	if errors.Is(err, serviceErrors.ErrProhibited) {
		logger.LogRequestWithError(ctx, "SendMessage", r.String(), serviceErrors.ErrProhibited)
		return nil, serviceErrors.ErrProhibited
	}

	//Creating, logging and sending response
	response := &proto.SendMessageResponse{}
	logger.LogRequest(ctx, "SendMessage", r.String(), response.String())
	return response, nil
}

// Implementation of GetHistory rpc
func (s Server) GetHistory(ctx context.Context, r *proto.GetHistoryRequest) (*proto.GetHistoryResponse, error) {
	// if invalid chatUuid provided - request cannot be completed, return error
	if !checkuuid.IsParsed(r.GetChatUuid()) {
		logger.LogRequestWithError(ctx, "GetHistory", r.String(), serviceErrors.ErrInvalidUuid)
		return nil, serviceErrors.ErrInvalidUuid
	}

	//get history from storage with chatUuid provided
	history, err := Storage.GetHistory(r.ChatUuid)

	//If chat not found - logging and returning error
	if errors.Is(err, serviceErrors.ErrChatNotFound) {
		logger.LogRequestWithError(ctx, "GetHistory", r.String(), serviceErrors.ErrChatNotFound)
		return nil, serviceErrors.ErrChatNotFound
	}

	//Create response, log and send
	response := &proto.GetHistoryResponse{Messages: history}
	logger.LogRequest(ctx, "GetHistory", r.String(), "History for "+r.String()+" successfully sent")
	return response, nil
}

// implementation of GetActiveChats rpc
func (s Server) GetActiveChats(ctx context.Context, r *proto.GetActiveChatsRequest) (*proto.GetActiveChatsResponse, error) {
	chats := Storage.GetActiveChats()
	response := &proto.GetActiveChatsResponse{Chats: chats}
	logger.LogRequest(ctx, "GetActiveChats", r.String(), "Active chats successfully sent")
	return response, nil
}
