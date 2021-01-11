package zap

import (
	"context"

	"go.uber.org/zap"
)

type key uint8

var logKey key

func NewContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, logKey, logger)
}

func FromContext(ctx context.Context) (*zap.Logger, bool) {
	l, ok := ctx.Value(logKey).(*zap.Logger)
	return l, ok
}
