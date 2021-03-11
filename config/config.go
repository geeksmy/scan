package config

import (
	"os"

	"github.com/geeksmy/go-lib/redis"
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
)

type Config struct {
	Debug    bool           `yaml:"debug,omitempty" default:"false" `
	Logger   LoggerConfig   `yaml:"logger,omitempty"`
	Redis    redis.Conf     `yaml:"redis,omitempty"`
	Database DatabaseConfig `yaml:"database,omitempty"`
	Port     PortConfig     `yaml:"port,omitempty"`
}

type DatabaseConfig struct {
	// 仅支持 mysql
	DSN          string `yaml:"dsn"`
	LogMode      bool   `yaml:"log_mode" default:"false"`
	MaxIdleConns int    `yaml:"max_idle_conns" default:"10"`
	MaxOpenConns int    `yaml:"max_open_conns" default:"100"`
	// format: https://golang.org/pkg/time/#ParseDuration
	ConnMaxLifetime string `yaml:"conn_max_lifetime" default:"1h"`
}

type LoggerConfig struct {
	Level string `yaml:"level,omitempty" default:"debug"`
	// json or text
	Format string `yaml:"format,omitempty" default:"json"`
	// file
	Output string `yaml:"output,omitempty" default:""`
}

type PortConfig struct {
	Protocol        string   `yaml:"protocol,omitempty"`
	FingerprintFile string   `yaml:"fingerprint_file,omitempty"`
	TargetIPs       []string `yaml:"target_ips,omitempty"`
	TargetPorts     []string `yaml:"target_ports,omitempty"`
	Timeout         int      `yaml:"timeout,omitempty"`
	Thread          int      `yaml:"thread,omitempty"`
	Retry           int      `yaml:"retry,omitempty"`
}

func InitLogger(debug bool, level, output string) {
	var conf zap.Config
	if debug {
		conf = zap.NewDevelopmentConfig()
	} else {
		conf = zap.NewProductionConfig()
	}

	var zapLevel = zap.NewAtomicLevel()
	if err := zapLevel.UnmarshalText([]byte(level)); err != nil {
		zap.L().Panic("设置日志记录级别失败",
			zap.Strings("only", []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}),
			zap.Error(err),
		)
	}

	conf.Level = zapLevel
	conf.Encoding = "console"

	if output != "" {
		conf.OutputPaths = []string{output}
		conf.ErrorOutputPaths = []string{output}
	}

	logger, _ := conf.Build()

	zap.RedirectStdLog(logger)
	zap.ReplaceGlobals(logger)
}

func InitRedis(conf redis.Conf) {
	redis.C = conf
}

var C = &Config{}

func Init(cfgFile string) {
	_ = os.Setenv("SCAN", "-")

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	if cfgFile != "" {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(C, cfgFile); err != nil {
			zap.L().Panic("init config fail", zap.Error(err))
		}
	} else {
		if err := configor.New(&configor.Config{AutoReload: true}).Load(C); err != nil {
			zap.L().Panic("init config fail", zap.Error(err))
		}
	}

	InitLogger(C.Debug, C.Logger.Level, C.Logger.Output)
	InitRedis(C.Redis)

	zap.L().Debug("[+]: 加载配置文件 ->", zap.String("文件名", cfgFile))
}

func init() {
	C = &Config{}
}
