// grpc interceptor for request log
package interceptors

import (
	"context"

	"github.com/Rolan335/grpcMessenger/server/internal/logger"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// interface from fmt package to type assert and call method string
type Stringer interface {
	String() string
}

//nolint:wrapcheck
func Log(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//making request
	resp, err := handler(ctx, req)
	//log request
	if ctx.Err() == context.DeadlineExceeded {
		err = status.Error(codes.DeadlineExceeded, "context deadline exceeded")
	}
	logger.LogRequest(ctx, info.FullMethod, req.(Stringer).String(), resp.(Stringer).String(), err)
	return resp, err
}
