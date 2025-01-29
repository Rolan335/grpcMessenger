package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}

// initializing logger depending on env that currently is on
func Init(env string, out io.Writer) {
	switch env {
	case "local":
		Logger = slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "dev":
		Logger = slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		Logger = slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		Logger = slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
}

// Function logs request
func LogRequest(ctx context.Context, route string, req string, res string, err error) {
	if err != nil {
		Logger.LogAttrs(ctx, slog.LevelInfo, route,
			slog.String("Request", req),
			slog.String("Err", err.Error()),
		)
		return
	}
	Logger.LogAttrs(ctx, slog.LevelInfo, route,
		slog.String("Request", req),
		slog.String("Response", res),
	)
}

func LogKafkaSuccess(partition int, offset int) {
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "Successfully sent message to kafka",
		slog.Int("Partition", partition),
		slog.Int("Offset", offset),
	)
}

func LogKafkaError(err error) {
	Logger.LogAttrs(context.Background(), slog.LevelError, "Failed to send message to kafka",
		slog.String("Err", err.Error()),
	)
}

// Function logs deletion of chat
func LogChatDelete(sessionUUID string, ChatUUID string, err error) {
	// if error from chatDelete provided - logging error
	if err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "DeleteChat",
			slog.String("sessionUuid", sessionUUID),
			slog.String("chatUuid", ChatUUID),
			slog.String("error", err.Error()),
		)
		return
	}
	//if no error, logging without error (to avoid nil pointer dereference panic)
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "DeleteChat",
		slog.String("sessionUuid", sessionUUID),
		slog.String("chatUuid", ChatUUID),
	)
}
