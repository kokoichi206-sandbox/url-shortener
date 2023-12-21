package logger

import (
	"context"
)

type Logger interface {
	LogFunc
	FormatLogFunc
	Print(ctx context.Context, lv Level, msg string)
}

type LogFunc interface {
	Critical(ctx context.Context, msg string)
	Error(ctx context.Context, mst string)
	Warn(ctx context.Context, mst string)
	Info(ctx context.Context, mst string)
	Debug(ctx context.Context, mst string)
}

type FormatLogFunc interface {
	Criticalf(ctx context.Context, msg string, a ...interface{})
	Errorf(ctx context.Context, mst string, a ...interface{})
	Warnf(ctx context.Context, mst string, a ...interface{})
	Infof(ctx context.Context, mst string, a ...interface{})
	Debugf(ctx context.Context, mst string, a ...interface{})
}

// Logging level.
type Level int

const (
	Degub    = -1
	Info     = 0
	Warn     = 1
	Error    = 2
	Critical = 3
)

func (lv Level) String() string {
	switch lv {
	case Degub:
		return "DEGUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Critical:
		return "CRITICAL"
	default:
		return "UNKNOWN"
	}
}

func shouldPrint(setting, given Level) bool {
	return setting <= given
}
