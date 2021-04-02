package config

import (
	"os"

	"github.com/geeksmy/go-libs/redis"
	"github.com/jinzhu/configor"
	"go.uber.org/zap"
)

type Config struct {
	Debug             bool                 `yaml:"debug,omitempty" default:"false" `
	Logger            LoggerConfig         `yaml:"logger,omitempty"`
	DisableStacktrace bool                 `yaml:"disable-stacktrace,omitempty" default:"true"`
	Redis             redis.Conf           `yaml:"redis,omitempty"`
	Database          DatabaseConfig       `yaml:"database,omitempty"`
	Port              PortConfig           `yaml:"port,omitempty"`
	Blasting          BlastingConfig       `yaml:"brute,omitempty"`
	WebFingerprint    WebFingerprintConfig `yaml:"web,omitempty"`
	Cyberspace        Cyberspace           `yaml:"cyberspace,omitempty"`
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
	Level string `yaml:"level,omitempty" default:"info"`
	// json or text
	Format string `yaml:"format,omitempty" default:"json"`
	// file
	Output string `yaml:"output,omitempty" default:""`
}

type PortConfig struct {
	Protocol        string   `yaml:"protocol,omitempty"`
	FingerprintFile string   `yaml:"fingerprint-file,omitempty"`
	TargetIPs       []string `yaml:"target-ips,omitempty"`
	TargetFile      string   `yaml:"target-file,omitempty" default:""`
	TargetPorts     []string `yaml:"target-ports,omitempty"`
	Timeout         int      `yaml:"timeout,omitempty"`
	Thread          int      `yaml:"thread,omitempty"`
	Retry           int      `yaml:"retry,omitempty"`
}

type BlastingConfig struct {
	TargetHost string   `yaml:"target-host,omitempty"`
	UserFile   string   `yaml:"user-file,omitempty"`
	PassFile   string   `yaml:"pass-file,omitempty"`
	Delay      int      `yaml:"delay,omitempty"`
	Thread     int      `yaml:"thread,omitempty"`
	Timeout    int      `yaml:"timeout,omitempty"`
	Retry      int      `yaml:"retry,omitempty"`
	ScanPort   bool     `yaml:"scan-port,omitempty"`
	Services   []string `yaml:"services,omitempty"`
	Path       string   `yaml:"path,omitempty"`
	TomcatPath string   `yaml:"tomcat-path,omitempty"`
}

type WebFingerprintConfig struct {
	TargetUrls      string   `yaml:"target-urls,omitempty"`
	TargetPorts     []string `yaml:"target-ports,omitempty"`
	FingerprintName string   `yaml:"fingerprint-name,omitempty"`
	Thread          int      `yaml:"thread,omitempty"`
	Timeout         int      `yaml:"timeout,omitempty"`
	Retry           int      `yaml:"retry,omitempty"`
}

type Cyberspace struct {
	Engine  string `yaml:"engine,omitempty"`
	Keyword string `yaml:"keyword,omitempty"`
}

func InitLogger(debug, disableStacktrace bool, level, output string) {
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
	conf.DisableStacktrace = disableStacktrace
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

	InitLogger(C.Debug, C.DisableStacktrace, C.Logger.Level, C.Logger.Output)
	// InitRedis(C.Redis)

	zap.L().Debug("[+]: 加载配置文件 ->", zap.String("文件名", cfgFile))
}

func init() {
	C = &Config{}
}
