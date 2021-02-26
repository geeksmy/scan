package cli

import (
	"net"
	"strings"
	"time"

	"scan/config"
	"scan/internal/model"
	"scan/internal/service/cli"
	"scan/pkg/tools"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Port struct {
	cmd    *cobra.Command
	logger *zap.Logger
	Probes *[]model.Probe

	Protocol        string
	FingerprintFile string
	TargetIPs       *[]string
	TargetPorts     *[]string
	Timeout         int
	Thread          int
	Retry           int
}

type ScanResult struct {
	IP       string
	Port     string
	State    string
	Response []byte
	Retry    int
}

func NewPort(cmd *cobra.Command, logger *zap.Logger) *Port {
	return &Port{
		cmd:    cmd,
		logger: logger,
	}
}

func (p *Port) PortMain() error {
	// 初始化参数
	err := p.initArgs()
	if err != nil {
		return err
	}

	// 初始化规则文件
	parse := cli.NewParse(p.logger)
	probes, err := parse.ParsingNmapFingerprint(p.FingerprintFile)
	if err != nil {
		return err
	}
	p.Probes = probes

	// 初始化扫描
	p.initScanPort()

	return nil
}

func (p *Port) initArgs() error {
	conf := config.C
	// 设置默认配置参数

	protocol, _ := p.cmd.Flags().GetString("protocol")
	switch protocol {
	case "":
		// 如果不传命令行参数使用配置文件的配置
		p.Protocol = conf.Port.Protocol
	default:
		p.Protocol = strings.ToLower(protocol)
	}
	targetIPs, _ := p.cmd.Flags().GetStringArray("target-ips")
	switch len(targetIPs) {
	case 1:
		ips, err := tools.UnfoldIPs(p.string2strings(targetIPs[0]))
		if err != nil {
			return err
		}
		p.TargetIPs = ips
	case 0:

		ips, err := tools.UnfoldIPs(conf.Port.TargetIPs)
		if err != nil {
			return err
		}
		p.TargetIPs = ips
	default:
		ips, err := tools.UnfoldIPs(targetIPs)
		if err != nil {
			return err
		}
		p.TargetIPs = ips
	}
	targetPorts, _ := p.cmd.Flags().GetStringArray("target-ports")
	switch len(targetPorts) {
	case 1:
		ports, err := tools.UnfoldPort(p.string2strings(targetPorts[0]))
		if err != nil {
			return err
		}
		p.TargetPorts = ports
	case 0:
		// 如果不传命令行参数使用配置文件的配置
		ports, err := tools.UnfoldPort(conf.Port.TargetPorts)
		if err != nil {
			return err
		}
		p.TargetPorts = ports
	default:
		ports, err := tools.UnfoldPort(targetPorts)
		if err != nil {
			return err
		}
		p.TargetPorts = ports
	}

	timeout, _ := p.cmd.Flags().GetInt("timeout")
	switch timeout {
	case 0:
		p.Timeout = conf.Port.Timeout
	default:
		p.Timeout = timeout
	}
	thread, _ := p.cmd.Flags().GetInt("thread")
	switch thread {
	case 0:
		p.Thread = conf.Port.Thread
	default:
		p.Thread = thread
	}
	retry, _ := p.cmd.Flags().GetInt("retry")
	switch retry {
	case 0:
		p.Retry = conf.Port.Retry
	default:
		p.Retry = retry
	}
	fingerprintFile, _ := p.cmd.Flags().GetString("fingerprint-file")
	switch fingerprintFile {
	case "":
		p.FingerprintFile = conf.Port.FingerprintFile
	default:
		p.FingerprintFile = fingerprintFile
	}

	return nil
}

func (p *Port) getArgs(argKey string) string {
	argValue := p.cmd.Flags().Lookup(argKey).Value.String()
	return argValue
}

func (p *Port) string2strings(s string) []string {
	s = s[1 : len(s)-1]
	return strings.Split(s, ",")
}

func (p *Port) initScanPort() []*ScanResult {
	p.logger.Info("[+] 开始扫描")
	var results []*ScanResult

	targetIps := *p.TargetIPs
	targetPorts := *p.TargetPorts

	// 暂时使用单线程
	for i := 0; i < len(targetIps); i++ {
		p.logger.Info("[+] 扫描 -> ", zap.String("IP", targetIps[i]))
		p.logger.Info("[+] 扫描 -> ", zap.Strings("Port", targetPorts))
		for j := 0; j < len(targetPorts); j++ {
			result := p.scanPort(p.Protocol, targetIps[i], targetPorts[j])
			results = append(results, result)
		}
	}

	return results
}

func (p *Port) scanPort(protocol, targetIP, targetPort string) *ScanResult {
	probe := "Hello Hacking!"

	result := ScanResult{IP: targetIP, Port: targetPort}
	address := targetIP + ":" + targetPort

	dialer := net.Dialer{Timeout: time.Duration(p.Timeout) * time.Second}
	conn, err := dialer.Dial(protocol, address)
	if err != nil {
		result.State = "Close"
		return &result
	}
	defer conn.Close()

	for i := 0; i < p.Retry; i++ {
		// 发送从nmap-service-probes拿到的 probe data
		if len(probe) > 0 {
			_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(p.Timeout) * time.Second))
			_, err = conn.Write([]byte(probe))
			if err != nil {
				result.State = "Close"
				return &result
			}
		}

		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(p.Timeout) * time.Second))
		for true {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				if len(result.Response) > 0 {
					break
				} else {
					result.State = "Close"
					return &result
				}
			}
			if n > 0 {
				result.State = "Open"
				result.Response = append(result.Response, buff[:n]...)
			}
		}
	}

	return &result
}

/**
 * 1. 发送数据包
 *  1.1 整理需要发送数据包的参数
 *  1.2 发送数据包
 * 2. 指纹识别
 * 3. 输出打印
 */

// 数据包发送
func (p *Port) packageSend() {
	/**
	 * 1. 整理需要发送的数据包参数
	 *
	 */
}

// 指纹识别
func (p *Port) fingerprintRecognition() {

}

// 输出打印
func (p *Port) outputPrinting() {

}
