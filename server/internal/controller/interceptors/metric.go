// grpc interceptor for metric
package interceptors

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// counter of all successful and failed requests
var RequestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "requests_total",
		Help: "total number of requests to service",
	},
	[]string{"method", "status"},
)

func Metric(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	//making request
	resp, err := handler(ctx, req)

	//OK code in default
	statusCode := codes.OK

	//if request returned error, get wrapped code
	if err != nil {
		statusCode = status.Code(err)
	}

	//incrementing request counter according to label values
	RequestsCounter.WithLabelValues(info.FullMethod, statusCode.String()).Inc()
	return resp, err
}
