package metric

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// counter of all successful and failed requests
var RequestsCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "requests_total",
		Help: "total number of requests to service",
	},
	[]string{"method", "status"},
)

// histogram of time of the query completion
var ResponseDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "grpc_duration_seconds",
		Help:    "Histogram of response durations for gRPC calls.",
		Buckets: prometheus.DefBuckets,
	},
	[]string{"method"},
)

var ChatsCreatedTTL = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "chats_ttl",
		Help:    "ttl time of chats created in seconds (zero is no ttl)",
		Buckets: prometheus.LinearBuckets(0, 600, 20),
	},
)

var UsersRegisteredTotal = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "user_registered_total",
		Help: "total number of users registered",
	},
)

var MessagesPerChat = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "messages_per_chat_total",
		Help: "Total number of sent messages in chat",
	},
	[]string{"chat_uuid"},
)

var once sync.Once

func MustInit() {
	once.Do(func() {
		prometheus.MustRegister(
			RequestsCounter,
			ResponseDuration,
			ChatsCreatedTTL,
			UsersRegisteredTotal,
			MessagesPerChat,
		)
	})
}
