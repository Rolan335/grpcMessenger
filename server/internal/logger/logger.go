package logger

import (
	"context"
	"io"
	"log/slog"
)

type Logger struct {
	*slog.Logger
}

func Init(env string, out io.Writer) Logger {
	switch env {
	case "local":
		return Logger{slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))}
	case "dev":
		return Logger{slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelDebug}))}
	case "prod":
		return Logger{slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))}
	default:
		return Logger{slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))}
	}
}

// Function logs request
func (l *Logger) LogRequest(ctx context.Context, route string, req string, res string, err error) {
	if err != nil {
		l.LogAttrs(ctx, slog.LevelInfo, route,
			slog.String("Request", req),
			slog.String("Err", err.Error()),
		)
		return
	}
	l.LogAttrs(ctx, slog.LevelInfo, route,
		slog.String("Request", req),
		slog.String("Response", res),
	)
}

// Function logs deletion of chat
func (l *Logger) LogChatDelete(sessionUuid string, ChatUuid string, err error) {
	// if error from chatDelete provided - logging error
	if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "DeleteChat",
			slog.String("sessionUuid", sessionUuid),
			slog.String("chatUuid", ChatUuid),
			slog.String("error", err.Error()),
		)
		return
	}
	//if no error, logging without error (to avoid nil pointer dereference panic)
	l.LogAttrs(context.Background(), slog.LevelInfo, "DeleteChat",
		slog.String("sessionUuid", sessionUuid),
		slog.String("chatUuid", ChatUuid),
	)
}
