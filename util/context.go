package util

import (
	"context"

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
