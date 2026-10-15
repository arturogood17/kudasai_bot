package state

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

type LogFile struct {
	slog.Handler
}

func (l LogFile) Enabled(c context.Context, lvl slog.Level) bool {
	return lvl == slog.LevelInfo
}

func (l LogFile) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandler := l.Handler.WithAttrs(attrs)
	return LogFile{newHandler}
}

func (l LogFile) WithGroup(name string) slog.Handler {
	newHandler := l.Handler.WithGroup(name)
	return LogFile{newHandler}
}

func NewMultilogger(ErrorFile, LogInfoFile *os.File) *slog.Logger {
	LF := LogFile{
		slog.NewJSONHandler(LogInfoFile, nil),
	}
	ErrorHandler := NewHandler(ErrorFile, slog.LevelWarn)
	return slog.New(slog.NewMultiHandler(LF, ErrorHandler))
}

func NewHandler(FileToLog io.Writer, level slog.Level) *slog.JSONHandler {
	return slog.NewJSONHandler(FileToLog, &slog.HandlerOptions{Level: level})
}

func NewLogError(logger *slog.Logger, t time.Time, message string, err error) {
	logger.Error(message, "Date", t, "Error", err)
}

func NewLogPostSuccess(logger *slog.Logger, t time.Time, PubDate time.Time, message string) {
	logger.Info(message, "Date", t, "Publication_Date", PubDate)
}
