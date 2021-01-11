package zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var level = zap.NewAtomicLevel()

func InitZap(env string, opts ...zap.Option) {
	var conf zap.Config
	if env == ENV_PRODDUCTION {
		conf = zap.NewProductionConfig()
	} else {
		conf = zap.NewDevelopmentConfig()
	}

	conf.Level = level

	logger, _ := conf.Build(opts...)

	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)
}

func SetLevelFromString(lv string) error {
	return level.UnmarshalText([]byte(lv))
}

func SetLevel(lv zapcore.Level) {
	level.SetLevel(lv)
}

// 是否可以设置 *zap.Logger
type CanZapSettable interface {
	// 设置 *zap.Logger
	SetZap(l *zap.Logger)
}

func BindZap(i interface{}, l *zap.Logger) {
	ii, ok := i.(CanZapSettable)
	if !ok {
		panic(fmt.Sprintf("%+v is not settable for zap", i))
	}

	ii.SetZap(l)
}
