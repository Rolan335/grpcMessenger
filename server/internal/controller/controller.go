// Package controller provides implementation of Service routes with proceeding requests, sending data to storage, logging and sending response.
package controller

import (
	"context"
	"errors"

	"github.com/Rolan335/grpcMessenger/proto"
	chatttl "github.com/Rolan335/grpcMessenger/server/internal/chatTtl"
	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/internal/serviceErrors"
	"github.com/Rolan335/grpcMessenger/server/internal/storage"
	"github.com/Rolan335/grpcMessenger/server/internal/util/checkuuid"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		errCreated := status.Error(codes.InvalidArgument, serviceErrors.ErrInvalidUuid.Error())
		logger.LogRequestWithError(ctx, "CreateChat", r.String(), errCreated)
		return nil, errCreated
	}

	//Creating uuid for chat
	id, _ := uuid.NewRandom()

	//Add new chat to server storage
	err := Storage.AddChat(r.SessionUuid, int(r.Ttl), r.ReadOnly, id.String())

	//If nonExistent session-uuid provided - logging and returning error
	if errors.Is(err, serviceErrors.ErrUserDoesNotExist) {
		errCreated := status.Error(codes.NotFound, err.Error())
		logger.LogRequestWithError(ctx, "CreateChat", r.String(), errCreated)
		return nil, errCreated
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
		errCreated := status.Error(codes.InvalidArgument, serviceErrors.ErrInvalidUuid.Error())
		logger.LogRequestWithError(ctx, "SendMessage", r.String(), errCreated)
		return nil, errCreated
	}

	//Creating uuid for message
	id, _ := uuid.NewRandom()
	//Adding new message to storage and if failed - returns error
	err := Storage.AddMessage(r.SessionUuid, r.ChatUuid, id.String(), r.Message)
	//Check if chat not found - send error and log it
	if errors.Is(err, serviceErrors.ErrChatNotFound) {
		errCreated := status.Error(codes.NotFound, err.Error())
		logger.LogRequestWithError(ctx, "SendMessage", r.String(), errCreated)
		return nil, errCreated
	}

	//Check if chat is readonly
	if errors.Is(err, serviceErrors.ErrProhibited) {
		errCreated := status.Error(codes.PermissionDenied, err.Error())
		logger.LogRequestWithError(ctx, "SendMessage", r.String(), errCreated)
		return nil, errCreated
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
		errCreated := status.Error(codes.InvalidArgument, serviceErrors.ErrInvalidUuid.Error())
		logger.LogRequestWithError(ctx, "GetHistory", r.String(), errCreated)
		return nil, errCreated
	}

	//get history from storage with chatUuid provided
	history, err := Storage.GetHistory(r.ChatUuid)

	//If chat not found - logging and returning error
	if errors.Is(err, serviceErrors.ErrChatNotFound) {
		errCreated := status.Error(codes.NotFound, err.Error())
		logger.LogRequestWithError(ctx, "GetHistory", r.String(), errCreated)
		return nil, errCreated
	}

	//Create response, log and send
	response := &proto.GetHistoryResponse{Messages: history}
	logger.LogRequest(ctx, "GetHistory", r.String(), "History for "+r.String()+" successfully sent")
	return response, nil
}
