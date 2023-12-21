package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/kokoichi206-sandbox/url-shortener/util"
)

type basicLogger struct {
	level   Level
	writer  io.Writer
	host    string
	service string
}

type logMessage struct {
	Host      string `json:"hostname"`
	Service   string `json:"service"`
	Message   string `json:"message"`
	Status    string `json:"status"`
	RequestID string `json:"request_id,omitempty"`
}

func NewBasicLogger(
	writer io.Writer,
	host string,
	service string,
) Logger {
	logger := &basicLogger{
		level:   Info,
		writer:  writer,
		host:    host,
		service: service,
	}

	return logger
}

func (b *basicLogger) Critical(ctx context.Context, msg string) {
	b.Print(ctx, Critical, msg)
}

func (b *basicLogger) Error(ctx context.Context, msg string) {
	b.Print(ctx, Error, msg)
}

func (b *basicLogger) Warn(ctx context.Context, msg string) {
	b.Print(ctx, Warn, msg)
}

func (b *basicLogger) Info(ctx context.Context, msg string) {
	b.Print(ctx, Info, msg)
}

func (b *basicLogger) Debug(ctx context.Context, msg string) {
	b.Print(ctx, Degub, msg)
}

func (b *basicLogger) Criticalf(ctx context.Context, msg string, a ...interface{}) {
	b.Print(ctx, Critical, fmt.Sprintf(msg, a...))
}

func (b *basicLogger) Errorf(ctx context.Context, msg string, a ...interface{}) {
	b.Print(ctx, Error, fmt.Sprintf(msg, a...))
}

func (b *basicLogger) Warnf(ctx context.Context, msg string, a ...interface{}) {
	b.Print(ctx, Warn, fmt.Sprintf(msg, a...))
}

func (b *basicLogger) Infof(ctx context.Context, msg string, a ...interface{}) {
	b.Print(ctx, Info, fmt.Sprintf(msg, a...))
}

func (b *basicLogger) Debugf(ctx context.Context, msg string, a ...interface{}) {
	b.Print(ctx, Degub, fmt.Sprintf(msg, a...))
}

func (b *basicLogger) Print(ctx context.Context, level Level, msg string) {
	if !shouldPrint(b.level, level) {
		return
	}

	logMsg := logMessage{
		Host:      b.host,
		Service:   b.service,
		Message:   msg,
		Status:    level.String(),
		RequestID: util.GetRequestID(ctx),
	}

	//nolint:errchkjson
	jsonBytes, _ := json.Marshal(logMsg)

	//nolint:errcheck
	b.writer.Write(jsonBytes)
}

// Set Level after struct is initialized.
// The default log level is set to Info.
func (b *basicLogger) SetLevel(lv Level) {
	b.level = lv
}
