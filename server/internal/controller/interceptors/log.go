// grpc interceptor for request log
package interceptors

import (
	"context"

	"github.com/Rolan335/grpcMessenger/server/internal/logger"

	"google.golang.org/grpc"
)

// interface from fmt package to type assert and call method string
type Stringer interface {
	String() string
}

func Log(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		//making request
		resp, err := handler(ctx, req)
		//log request
		logger.LogRequest(ctx, info.FullMethod, req.(Stringer).String(), resp.(Stringer).String(), err)
		return resp, err
	}
}
