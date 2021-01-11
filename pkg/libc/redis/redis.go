// package redis 提供便利方式创建 redis 连接, 主要是提供 pflag 统一命令行参数风格
// ps: 暂未提供创建 ClusterClient 的能力
package redis

import (
	"time"

	"github.com/go-redis/redis/v7"
)

const (
	MaxRetries = 0
	MaxConnAge = time.Hour
)

var (
	// 默认 Client, 使用 Connect() 创建连接后会自动初始化
	Client *redis.Client
	// // 默认的 ClusterClient
	// ClusterClient *redis.ClusterClient
)

// 使用默认配置创建一个全局 Client
func Connect() error {
	if C.Sentinel {
		connector := NewFailoverConnector(C)
		cli, err := connector.Connect()
		if err != nil {
			return err
		}

		Client = cli
		return nil
	} else {
		connector, err := NewConnector(C)
		if err != nil {
			return err
		}

		cli, err := connector.Connect()
		if err != nil {
			return err
		}

		Client = cli
		return nil
	}
}

// 使用 Connector 可以提供 Options 的修改, 方便开发者设置具体的连接参数
func NewConnector(conf Conf) (*Connector, error) {
	opts, err := redis.ParseURL(conf.URI)
	if err != nil {
		return nil, err
	}
	opts.PoolSize = conf.PoolSize

	connector := &Connector{
		options: opts,
	}
	return connector, nil
}

type Connector struct {
	options *redis.Options
}

// Options 返回 *redis.Options 可以设置连接选项
func (c *Connector) Options() *redis.Options {
	return c.options
}

// Connect 创建 redis 连接
func (c *Connector) Connect() (*redis.Client, error) {
	cli := redis.NewClient(c.options)
	if err := cli.Ping().Err(); err != nil {
		return nil, err
	}
	return cli, nil
}

func NewFailoverConnector(conf Conf) *FailoverConnector {
	opt := redis.FailoverOptions{
		MasterName:       conf.MasterName,
		Password:         conf.Password,
		SentinelAddrs:    conf.SentinelAddrs,
		SentinelPassword: conf.SentinelPassword,
		DB:               conf.DB,
		PoolSize:         conf.PoolSize,
		MaxRetries:       MaxRetries,
		MaxConnAge:       MaxConnAge,
	}

	connector := FailoverConnector{options: &opt}
	return &connector
}

type FailoverConnector struct {
	options *redis.FailoverOptions
}

func (c *FailoverConnector) Options() *redis.FailoverOptions {
	return c.options
}
func (c *FailoverConnector) Connect() (*redis.Client, error) {
	cli := redis.NewFailoverClient(c.options)
	if err := cli.Ping().Err(); err != nil {
		return nil, err
	}

	return cli, nil
}
