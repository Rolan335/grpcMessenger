package serviceErrors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrChatNotFound = status.Error(codes.NotFound, "chat not found")
var ErrUserDoesNotExist = status.Error(codes.NotFound, "user doesn't exist")
var ErrProhibited = status.Error(codes.PermissionDenied, "prohibited. Only creator can send")

var ErrInvalidUuid = status.Error(codes.InvalidArgument, "invalid uuid provided")
var ErrUnknownMsgType = status.Error(codes.Unknown, "Unknown message type")