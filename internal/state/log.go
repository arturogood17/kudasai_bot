package state

import (
	"context"
	"log/slog"
	"os"
	"time"
)

func NewJSONLogger() *slog.JSONHandler {
	return slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
}

func NewLogError(ctx context.Context, logger *slog.JSONHandler, t time.Time, level slog.Level, message string, err error) {
	record := slog.NewRecord(t, level, message, 0)
	record.Add("Error", err)
	logger.Handle(ctx, record)
}

func NewLogPostSuccess(ctx context.Context, logger *slog.JSONHandler, t time.Time, PubDate time.Time, level slog.Level, message string) {
	record := slog.NewRecord(t, level, message, 0)
	record.Add("Publication_Date", PubDate)
	logger.Handle(ctx, record)
}
