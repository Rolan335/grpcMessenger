package logger

import (
	"context"
	"io"
	"log/slog"
)

var Logger *slog.Logger

// initializing logger depending on env that currently is on
func LoggerInit(env string, out io.Writer) {
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

// Function logs successful request
func LogRequest(ctx context.Context, route string, req string, res string) {
	Logger.LogAttrs(ctx, slog.LevelInfo, route,
		slog.String("Request", req),
		slog.String("Response", res),
	)
}

// Function logs request with error
func LogRequestWithError(ctx context.Context, route string, req string, err error) {
	Logger.LogAttrs(ctx, slog.LevelInfo, route,
		slog.String("Request", req),
		slog.String("Error", err.Error()),
	)
}

// Function logs deletion of chat
func LogChatDelete(sessionUuid string, ChatUuid string, err error) {
	// if error from chatDelete provided - logging error
	if err != nil {
		Logger.LogAttrs(context.Background(), slog.LevelError, "DeleteChat",
			slog.String("sessionUuid", sessionUuid),
			slog.String("chatUuid", ChatUuid),
			slog.String("error", err.Error()),
		)
		return
	}
	//if no error, logging without error (to avoid nil pointer dereference panic)
	Logger.LogAttrs(context.Background(), slog.LevelInfo, "DeleteChat",
		slog.String("sessionUuid", sessionUuid),
		slog.String("chatUuid", ChatUuid),
	)
}
