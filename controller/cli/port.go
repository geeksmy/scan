package cli

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"scan/config"
	"scan/internal/model"
	"scan/internal/service/cli"
	"scan/pkg/tools"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// 回包结构图
type ScanResult struct {
	IP       string
	Port     string
	State    bool
	Protocol string
	Response []byte
	Probe    model.Probe
	Retry    int
	Error    error
}

// 命令行参数结构图
type CmdArgs struct {
	Protocol        string
	FingerprintFile string
	TargetIPs       *[]string
	TargetPorts     *[]string
	Timeout         int
	Thread          int
	Retry           int
	OutPut          string
	OutFileName     string
}

// 数据包参数结构图
type PackageArgs struct {
	Protocol   string
	TargetIP   string
	TargetPort string
	Probe      model.Probe
}

// 指纹识别后的结构图
type Result struct {
	IP         string
	Port       string
	State      string
	Protocol   string
	Retry      int
	ServerType string
	Version    string
	Banner     string
	IsSoft     bool
}

type Port struct {
	cmd    *cobra.Command
	logger *zap.Logger

	Probes  *[]model.Probe
	CmdArgs CmdArgs
}

func NewPort(cmd *cobra.Command, logger *zap.Logger) *Port {
	return &Port{
		cmd:    cmd,
		logger: logger,
	}
}

func (p *Port) PortMain() error {
	// 初始化命令参数
	err := p.initArgs()
	if err != nil {
		return err
	}

	// 初始化规则文件
	parse := cli.NewParse(p.logger)
	probes, err := parse.ParsingNmapFingerprint(p.CmdArgs.FingerprintFile)
	if err != nil {
		return err
	}
	p.Probes = probes

	// 初始化扫描参数
	var mainWG sync.WaitGroup
	packageArgsCh := make(chan PackageArgs, len(*p.CmdArgs.TargetIPs)*len(*p.CmdArgs.TargetPorts))
	scanResultCh := make(chan ScanResult, len(*p.CmdArgs.TargetIPs)*len(*p.CmdArgs.TargetPorts))
	resultCh := make(chan Result, len(*p.CmdArgs.TargetIPs)*len(*p.CmdArgs.TargetPorts))
	mainWG.Add(1)
	go p.initPackageArgs(packageArgsCh, &mainWG)

	// 扫描
	mainWG.Add(1)
	go p.sendPackage(packageArgsCh, scanResultCh, &mainWG)

	// 指纹识别
	mainWG.Add(1)
	go p.fingerprintRecognition(scanResultCh, resultCh, &mainWG)

	// 输出打印
	mainWG.Add(1)
	go p.outputPrinting(resultCh, &mainWG)

	mainWG.Wait()
	return nil
}

func (p *Port) initArgs() error {
	conf := config.C
	// 设置默认配置参数

	protocol, _ := p.cmd.Flags().GetString("protocol")
	switch protocol {
	case "":
		// 如果不传命令行参数使用配置文件的配置
		p.CmdArgs.Protocol = conf.Port.Protocol
	default:
		p.CmdArgs.Protocol = strings.ToLower(protocol)
	}
	targetIPs, _ := p.cmd.Flags().GetStringArray("target-ips")
	switch len(targetIPs) {
	case 1:
		ips, err := tools.UnfoldIPs(string2strings(targetIPs[0]))
		if err != nil {
			return err
		}
		p.CmdArgs.TargetIPs = ips
	case 0:

		ips, err := tools.UnfoldIPs(conf.Port.TargetIPs)
		if err != nil {
			return err
		}
		tools.Shuffle(*ips)
		p.CmdArgs.TargetIPs = ips
	default:
		ips, err := tools.UnfoldIPs(targetIPs)
		if err != nil {
			return err
		}
		tools.Shuffle(*ips)
		p.CmdArgs.TargetIPs = ips
	}
	targetPorts, _ := p.cmd.Flags().GetStringArray("target-ports")
	switch len(targetPorts) {
	case 1:
		ports, err := tools.UnfoldPort(string2strings(targetPorts[0]))
		if err != nil {
			return err
		}
		p.CmdArgs.TargetPorts = ports
	case 0:
		// 如果不传命令行参数使用配置文件的配置
		ports, err := tools.UnfoldPort(conf.Port.TargetPorts)
		if err != nil {
			return err
		}
		p.CmdArgs.TargetPorts = ports
	default:
		ports, err := tools.UnfoldPort(targetPorts)
		if err != nil {
			return err
		}
		p.CmdArgs.TargetPorts = ports
	}

	timeout, _ := p.cmd.Flags().GetInt("timeout")
	switch timeout {
	case 0:
		p.CmdArgs.Timeout = conf.Port.Timeout
	default:
		p.CmdArgs.Timeout = timeout
	}
	thread, _ := p.cmd.Flags().GetInt("thread")
	switch thread {
	case 0:
		p.CmdArgs.Thread = conf.Port.Thread
	default:
		p.CmdArgs.Thread = thread
	}
	retry, _ := p.cmd.Flags().GetInt("retry")
	switch retry {
	case 0:
		p.CmdArgs.Retry = conf.Port.Retry
	default:
		p.CmdArgs.Retry = retry
	}
	fingerprintFile, _ := p.cmd.Flags().GetString("fingerprint-file")
	switch fingerprintFile {
	case "":
		p.CmdArgs.FingerprintFile = conf.Port.FingerprintFile
	default:
		p.CmdArgs.FingerprintFile = fingerprintFile
	}
	outFile, _ := p.cmd.Flags().GetString("out-file")
	switch outFile {
	case "":
		p.CmdArgs.OutPut = "print"
	default:
		p.CmdArgs.OutPut = "file"
		p.CmdArgs.OutFileName = outFile
	}
	hostsFile, _ := p.cmd.Flags().GetString("target-file")
	switch hostsFile {
	case "":
		return nil
	default:
		ipsData, err := getTargetFile(hostsFile)
		if err != nil {
			p.logger.Error("[-] initArgs -> 解析目标文件失败")
			return err
		}
		ips, err := tools.UnfoldIPs(ipsData)
		if err != nil {
			return err
		}
		tools.Shuffle(*ips)
		p.CmdArgs.TargetIPs = ips
	}

	return nil
}

// 扫描参数线程池
func (p *Port) initPackageArgsWorker(portCh <-chan string, packageArgsCh chan<- PackageArgs, wg *sync.WaitGroup) {
	targetIPs := *p.CmdArgs.TargetIPs
	for port := range portCh {
		for i := 0; i < len(targetIPs); i++ {
			packageArgs := new(PackageArgs)
			packageArgs.Protocol = p.CmdArgs.Protocol
			packageArgs.TargetIP = targetIPs[i]
			packageArgs.TargetPort = port
			for _, probe := range *p.Probes {
				if tools.IncludePort(port, probe.Ports) && (strings.ToLower(p.CmdArgs.Protocol) == strings.ToLower(probe.Protocol)) {
					packageArgs.Probe = probe
					packageArgsCh <- *packageArgs
					continue
				}
				if probe.Name == "NULL" {
					packageArgs.Probe = probe
					packageArgsCh <- *packageArgs
					continue
				}
				if tools.IncludePort(port, probe.SSLPorts) && (strings.ToLower(p.CmdArgs.Protocol) == strings.ToLower(probe.Protocol)) {
					packageArgs.Probe = probe
					packageArgsCh <- *packageArgs
					continue
				}
				continue
			}
			// fmt.Printf("%s:%s--%d\n", packageArgs.TargetIP, packageArgs.TargetPort, len(packageArgs.Probe))
		}
		wg.Done()
	}
}

// 初始化扫描参数
func (p *Port) initPackageArgs(packageArgsCh chan PackageArgs, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(packageArgsCh)

	targetPorts := *p.CmdArgs.TargetPorts
	portCh := make(chan string, p.CmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(portCh); i++ {
		go p.initPackageArgsWorker(portCh, packageArgsCh, &wg)
	}

	for i := 0; i < len(targetPorts); i++ {
		wg.Add(1)
		portCh <- targetPorts[i]
	}

	wg.Wait()
	close(portCh)

}

// 发送数据包线程池
func (p *Port) sendPackageWorker(argsCh <-chan PackageArgs, scanResultCh chan<- ScanResult, wg *sync.WaitGroup) {
	for args := range argsCh {
		scanResult := ScanResult{
			IP:       args.TargetIP,
			Port:     args.TargetPort,
			Protocol: args.Protocol,
			Retry:    0,
		}
		address := fmt.Sprintf("%s:%s", args.TargetIP, args.TargetPort)

		// dialer := net.Dialer{Timeout: time.Duration(p.CmdArgs.Timeout) * time.Second}
		conn, err := net.DialTimeout(args.Protocol, address, time.Duration(p.CmdArgs.Timeout)*time.Second)
		if err != nil {
			scanResult.Error = err
			scanResult.Retry += 1
			scanResultCh <- scanResult
			wg.Done()
			continue
		}

		// 需要发送的数据
		if len(args.Probe.Data) > 0 && scanResult.Retry < 5 {
			_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(p.CmdArgs.Timeout) * time.Second))
			data, err := tools.DecodeData(args.Probe.Data)
			if err != nil {
				scanResult.Error = err
				scanResult.Retry += 1
				scanResultCh <- scanResult
				continue
			}
			_, err = conn.Write(data)
			if err != nil {
				scanResult.Error = err
				scanResult.Retry += 1
				scanResultCh <- scanResult
				continue
			}
		}

		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(p.CmdArgs.Timeout) * time.Second))
		for true {
			buff := make([]byte, 1024)
			n, err := conn.Read(buff)
			if err != nil {
				if len(scanResult.Response) > 0 {
					break
				} else {
					scanResult.Error = err
					scanResult.Retry += 1
					scanResultCh <- scanResult
					break
				}
			}
			if n > 0 {
				scanResult.State = true
				scanResult.Error = nil
				scanResult.Retry += 1
				scanResult.Probe = args.Probe
				scanResult.Response = append(scanResult.Response, buff[:n]...)
				scanResultCh <- scanResult
			}
		}
		wg.Done()
		conn.Close()
	}
}

// 发送数据包
func (p *Port) sendPackage(packageArgsCh <-chan PackageArgs, scanResultCh chan ScanResult, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(scanResultCh)

	argsCh := make(chan PackageArgs, p.CmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(argsCh); i++ {
		go p.sendPackageWorker(argsCh, scanResultCh, &wg)
	}

	for args := range packageArgsCh {
		wg.Add(1)
		argsCh <- args
	}

	wg.Wait()
	close(argsCh)
}

// 指纹识别线程池
func (p *Port) fingerprintRecognitionWorker(scanCh <-chan ScanResult, resultCh chan<- Result, wg *sync.WaitGroup) {
	for scan := range scanCh {
		result := Result{
			IP:       scan.IP,
			Port:     scan.Port,
			Protocol: scan.Protocol,
			Retry:    scan.Retry,
		}

		switch scan.State {
		case false:
			result.State = "Close"
		default:
			result.State = "Open"
		}

		for _, match := range scan.Probe.Matchs {
			banner := match.PatternCompiled.FindStringSubmatch(string(scan.Response))
			if len(banner) > 0 && !(result.Banner != "") {
				// 方便格式化输出
				result.ServerType = match.Service
				result.IsSoft = match.IsSoft
				version := match.PatternCompiled.ReplaceAllString(string(scan.Response), match.VersionInfo)
				if len(version) < 50 {
					result.Version = version
					result.Banner = strings.Replace(string(scan.Response), "\n", "", -1)
				}
			}
		}

		// 如果上面没有找到,则全局搜索
		if result.ServerType == "" && result.Version == "" {
			for _, probe := range *p.Probes {
				for _, match := range probe.Matchs {
					banner := match.PatternCompiled.FindStringSubmatch(string(scan.Response))
					if len(banner) > 0 && !(result.Banner != "") {
						// 方便格式化输出
						result.ServerType = match.Service
						result.IsSoft = match.IsSoft
						version := match.PatternCompiled.ReplaceAllString(string(scan.Response), match.VersionInfo)
						if len(version) < 50 {
							result.Version = version
							result.Banner = strings.Replace(string(scan.Response), "\n", "", -1)
						}
					}
				}
			}
		}

		resultCh <- result
		wg.Done()
	}
}

// 指纹识别
func (p *Port) fingerprintRecognition(scanResultCh <-chan ScanResult, resultCh chan Result, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(resultCh)

	scanCh := make(chan ScanResult, p.CmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(scanCh); i++ {
		go p.fingerprintRecognitionWorker(scanCh, resultCh, &wg)
	}

	for res := range scanResultCh {
		// 过滤掉端口关闭的包
		if res.State != false && res.Error == nil {
			wg.Add(1)
			scanCh <- res
		}
	}

	wg.Wait()
	close(scanCh)
}

// 输出文件
func (p *Port) outFile(res Result, wg *sync.WaitGroup, file *os.File) {
	_, _ = file.WriteString(fmt.Sprintf("%s:%s\t%s\t\t\t%s\t\t\t%s\n", res.IP, res.Port, res.ServerType, res.Version, res.Banner))
	wg.Done()
}

// 输出屏幕
func (p *Port) outCmd(res Result, wg *sync.WaitGroup) {
	fmt.Printf("%s:%s\t\t%s\t\t\t%s\t\t\t%s\n", res.IP, res.Port, res.ServerType, res.Version, res.Banner)
	wg.Done()
}

// 输出打印
func (p *Port) outputPrinting(resultCh <-chan Result, mainWG *sync.WaitGroup) {
	defer mainWG.Done()

	var wg sync.WaitGroup

	_, err := os.Stat(p.CmdArgs.OutFileName)
	if err == nil {
		// 如果文件存在
		_ = os.Remove("test.txt")
	}

	file, _ := os.Create(p.CmdArgs.OutFileName)

	switch p.CmdArgs.OutPut {
	case "file":
		_, _ = file.WriteString(fmt.Sprintf("%s:%s\t\t%s\t\t\t%s\t\t\t%s\n", "目标地址", "目标端口", "服务类型", "版本信息", "Banner"))
		for res := range resultCh {
			wg.Add(1)
			go p.outFile(res, &wg, file)
		}

	default:
		fmt.Printf("%s:%s\t\t%s\t\t\t%s\t\t\t%s\n", "目标地址", "目标端口", "服务类型", "版本信息", "Banner")
		for res := range resultCh {
			wg.Add(1)
			go p.outCmd(res, &wg)
		}
	}

	wg.Wait()
	_ = file.Close()
}

func string2strings(s string) []string {
	s = s[1 : len(s)-1]
	return strings.Split(s, ",")
}

func getTargetFile(targetFile string) ([]string, error) {
	fileData, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return nil, err
	}
	ips := strings.Split(string(fileData), "\n")

	return ips, nil
}
