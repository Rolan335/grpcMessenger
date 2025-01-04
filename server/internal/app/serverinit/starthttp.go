package serverinit

import (
	"context"
	"net/http"

	"github.com/Rolan335/grpcMessenger/server/internal/config"
	"github.com/Rolan335/grpcMessenger/server/internal/logger"
	"github.com/Rolan335/grpcMessenger/server/proto"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartHttp(ctx context.Context, serverConfig config.ServiceCfg) {
	mux := runtime.NewServeMux()
	err := proto.RegisterMessengerServiceHandlerFromEndpoint(
		ctx,
		mux,
		serverConfig.Address+serverConfig.PortGrpc,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		logger.Logger.Info("cannot register http handler: " + err.Error())
		panic("cannot register http handler: " + err.Error())
	}
	logger.Logger.Info("started http server")
	if err := http.ListenAndServe(serverConfig.PortHttp, mux); err != nil {
		logger.Logger.Info("cannot start http server: " + err.Error())
		panic("cannot start http server: " + err.Error())
	}
}
