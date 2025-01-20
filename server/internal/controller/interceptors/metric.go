// grpc interceptor for metric
package interceptors

import (
	"context"
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/metric"
	"github.com/Rolan335/grpcMessenger/server/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Metric(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//start timer
	start := time.Now()

	//making request
	resp, err := handler(ctx, req)

	//OK code in default
	statusCode := codes.OK

	//if request returned error, get wrapped code
	if err != nil {
		statusCode = status.Code(err)
	}

	if _, ok := req.(*proto.InitSessionRequest); ok && statusCode == codes.OK {
		metric.UsersRegisteredTotal.Inc()
	}
	if reqAsserted, ok := req.(*proto.CreateChatRequest); ok && statusCode == codes.OK {
		if reqAsserted.Ttl < 0 {
			reqAsserted.Ttl = 0
		}
		metric.ChatsCreatedTTL.Observe(float64(reqAsserted.Ttl))
	}

	if reqAsserted, ok := req.(*proto.SendMessageRequest); ok && statusCode == codes.OK {
		metric.MessagesPerChat.WithLabelValues(reqAsserted.ChatUuid).Inc()
	}

	//incrementing request counter according to label values
	metric.RequestsCounter.WithLabelValues(info.FullMethod, statusCode.String()).Inc()

	if statusCode == codes.OK {
		metric.ResponseDuration.WithLabelValues(statusCode.String()).Observe(time.Since(start).Seconds())
	}
	return resp, err
}
