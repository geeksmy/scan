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
	Cyberspace        CyberspaceConfig     `yaml:"cyberspace,omitempty"`
	PassGen           PassGenConfig        `yaml:"passgen,omitempty"`
	IntranetAlive     IntranetAliveConfig  `yaml:"intranet-alive,omitempty"`
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
	Protocol        string   `yaml:"protocol,omitempty" default:"tcp"`
	FingerprintFile string   `yaml:"fingerprint-file,omitempty" default:"nmap-service-probes"`
	TargetIPs       []string `yaml:"target-ips,omitempty"`
	TargetFile      string   `yaml:"target-file,omitempty"`
	TargetPorts     []string `yaml:"target-ports,omitempty"`
	Timeout         int      `yaml:"timeout,omitempty" default:"1"`
	Thread          int      `yaml:"thread,omitempty" default:"20"`
	Retry           int      `yaml:"retry,omitempty" default:"1"`
	OutFile         string   `yaml:"out-file,omitempty" default:"port.txt"`
}

type BlastingConfig struct {
	TargetHost string   `yaml:"target-host,omitempty"`
	UserFile   string   `yaml:"user-file,omitempty"`
	PassFile   string   `yaml:"pass-file,omitempty"`
	Delay      int      `yaml:"delay,omitempty" default:"1"`
	Thread     int      `yaml:"thread,omitempty" default:"20"`
	Timeout    int      `yaml:"timeout,omitempty" default:"1"`
	Retry      int      `yaml:"retry,omitempty" default:"1"`
	ScanPort   bool     `yaml:"scan-port,omitempty"`
	Services   []string `yaml:"services,omitempty"`
	Path       string   `yaml:"path,omitempty" default:"/login"`
	TomcatPath string   `yaml:"tomcat-path,omitempty" default:"/manager"`
	OutFile    string   `yaml:"out-file,omitempty" default:"brute.txt"`
}

type WebFingerprintConfig struct {
	TargetUrls      string   `yaml:"target-urls,omitempty"`
	TargetPorts     []string `yaml:"target-ports,omitempty"`
	FingerprintName string   `yaml:"fingerprint-name,omitempty"`
	Thread          int      `yaml:"thread,omitempty" default:"20"`
	Timeout         int      `yaml:"timeout,omitempty" default:"1"`
	Retry           int      `yaml:"retry,omitempty" default:"1"`
	OutFile         string   `yaml:"out-file,omitempty" default:"web.txt"`
}

type Fofa struct {
	Email         string `yaml:"email,omitempty"`
	Key           string `yaml:"key,omitempty"`
	Authorization string `yaml:"authorization,omitempty"`
}

type Shodan struct {
	Key string `yaml:"key,omitempty"`
}

type CyberspaceConfig struct {
	Engine  string `yaml:"engine,omitempty" default:"fofa"`
	Timeout int    `yaml:"timeout,omitempty" default:"1"`
	Thread  int    `yaml:"thread,omitempty" default:"20"`
	Search  string `yaml:"search,omitempty"`
	Fofa    Fofa   `yaml:"fofa,omitempty"`
	Shodan  Shodan `yaml:"shodan,omitempty"`
}

type PassGenConfig struct {
	Year       string `yaml:"year,omitempty"`
	DomainName string `yaml:"domain-name,omitempty"`
	Domain     string `yaml:"domain,omitempty"`
	Device     string `yaml:"device,omitempty"`
	Length     int    `yaml:"length,omitempty" default:"1"`
	OutFile    string `yaml:"out-file,omitempty" default:"pass.txt"`
}

type IntranetAliveConfig struct {
	Target  string  `yaml:"target,omitempty"`
	Thread  int     `yaml:"thread,omitempty" default:"20"`
	Timeout int     `yaml:"timeout,omitempty" default:"1"`
	Retry   int     `yaml:"retry,omitempty" default:"1"`
	Delay   float32 `yaml:"delay,omitempty"`
	OutFile string  `yaml:"out-file,omitempty" default:"survive.txt"`
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
