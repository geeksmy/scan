package zap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type zapBindable struct{}

func (z *zapBindable) SetZap(l *zap.Logger) {
	l.Info("call `SetZap` success")
}

func TestBindZap(t *testing.T) {
	bindable := &zapBindable{}
	logger := zap.L().With(zap.String("svc", "for test"))

	assert.NotPanics(t, func() { BindZap(bindable, logger) })

	cantBind := struct{}{}

	assert.Panics(t, func() { BindZap(cantBind, logger) })

	BindZap(bindable, logger)
}
