package util

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type requestIDKye struct{}

func WithRequestID(parent context.Context) context.Context {
	reqID := uuid.New().String()

	return context.WithValue(parent, requestIDKye{}, reqID)
}

func GetRequestID(ctx context.Context) string {
	reqID, ok := ctx.Value(requestIDKye{}).(string)
	if !ok {
		return ""
	}

	return reqID
}

type DetachedCtx struct {
	Parent context.Context
}

func (d DetachedCtx) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (d DetachedCtx) Done() <-chan struct{} {
	return nil
}

func (d DetachedCtx) Err() error {
	return nil
}

func (d DetachedCtx) Value(key any) any {
	return d.Parent.Value(key)
}
