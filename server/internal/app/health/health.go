package health

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/Rolan335/grpcMessenger/server/internal/repository/postgres"
	"github.com/Rolan335/grpcMessenger/server/internal/repository/redis"
	"github.com/Rolan335/grpcMessenger/server/pkg/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// nolint
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// nolint
func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	grpcAddr := os.Getenv("APP_ADDRESS") + os.Getenv("APP_PORTGRPC")
	grpcConn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		http.Error(w, "failed to dial grpc", http.StatusServiceUnavailable)
		return
	}
	defer grpcConn.Close()

	grpcClient := proto.NewMessengerServiceClient(grpcConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := grpcClient.HealthCheck(ctx, &proto.HealthCheckRequest{})
	if err != nil || resp.Status != "SERVING" {
		http.Error(w, "grpc is not ready", http.StatusServiceUnavailable)
		return
	}

	db := os.Getenv("APP_DB")
	switch db {
	case "postgres":
		if err := postgres.Ping(); err != nil {
			http.Error(w, "postgres is not ready", http.StatusServiceUnavailable)
			return
		}
	case "redis":
		if err := redis.Ping(); err != nil {
			http.Error(w, "redis is not ready", http.StatusServiceUnavailable)
			return
		}
	default:
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
	return
}
