package service

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"scan/config"
	"scan/internal/model"
	"scan/pkg/tools"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

type PortSVC interface {
	/**
	 * InitArgs 初始化参数
	 * @param cmd 命令行传入参数结构体
	 */
	InitArgs(cmd *cobra.Command) (*PortCmdArgs, error)
	/**
	 * InitPackageArgs 初始化扫描参数
	 * @param packageArgsCh 数据包回包 chan
	 */
	InitPackageArgs(packageArgsCh chan PackageArgs, mainWG *sync.WaitGroup, probes *[]model.Probe)
	/**
	 * SendPackage 发送数据包
	 * @param packageArgsCh 数据包参数 chan
	 * @param scanResultCh 数据包回包 chan
	 */
	SendPackage(packageArgsCh <-chan PackageArgs, scanResultCh chan ScanResult, mainWG *sync.WaitGroup)
	/**
	 * fingerprintRecognition 指纹识别
	 * @param scanResultCh 数据包回包 chan
	 * @param resultCh 指纹识别后的 chan
	 */
	FingerprintRecognition(scanResultCh <-chan ScanResult, resultCh chan Result, mainWG *sync.WaitGroup, probes *[]model.Probe)
	/**
	 * OutputPrinting 输出打印
	 * @param resultCh 指纹识别后的 chan
	 */
	OutputPrinting(resultCh <-chan Result)
}

type PackageArgs struct {
	Protocol   string
	TargetIP   string
	TargetPort string
	Probe      model.Probe
}

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

type PortCmdArgs struct {
	Protocol        string
	FingerprintFile string
	TargetIPs       *[]string
	TargetPorts     *[]string
	Timeout         int
	Thread          int
	Retry           int
	OutFileName     string
}

type Port struct {
	cmd    *cobra.Command
	logger *zap.Logger

	CmdArgs PortCmdArgs

	Results []Result
}

func NewPortSVC(logger *zap.Logger) PortSVC {
	return &Port{
		logger: logger,
	}
}

func (svc *Port) InitArgs(cmd *cobra.Command) (*PortCmdArgs, error) {
	conf := config.C
	// 设置默认配置参数

	protocol, _ := cmd.Flags().GetString("protocol")
	switch protocol {
	case "":
		// 如果不传命令行参数使用配置文件的配置
		svc.CmdArgs.Protocol = conf.Port.Protocol
	default:
		svc.CmdArgs.Protocol = strings.ToLower(protocol)
	}

	targetIPs, _ := cmd.Flags().GetStringArray("target-ips")
	switch len(targetIPs) {
	case 1:
		ips, err := tools.UnfoldIPs(tools.String2strings(targetIPs[0]))
		if err != nil {
			fmt.Println("[-] 参数错误")
			return nil, err
		}
		tools.Shuffle(*ips)
		svc.CmdArgs.TargetIPs = ips
	case 0:
		ips, err := tools.UnfoldIPs(conf.Port.TargetIPs)
		if err != nil {
			fmt.Println("[-] 参数错误")
			return nil, err
		}
		tools.Shuffle(*ips)
		svc.CmdArgs.TargetIPs = ips
	default:
		ips, err := tools.UnfoldIPs(targetIPs)
		if err != nil {
			fmt.Println("[-] 参数错误")
			return nil, err
		}
		tools.Shuffle(*ips)
		svc.CmdArgs.TargetIPs = ips
	}

	targetPorts, _ := cmd.Flags().GetStringArray("target-ports")
	switch len(targetPorts) {
	case 1:
		ports, err := tools.UnfoldPort(tools.String2strings(targetPorts[0]))
		if err != nil {
			fmt.Println("[-] 参数错误")
			return nil, err
		}
		svc.CmdArgs.TargetPorts = ports
	case 0:
		// 如果不传命令行参数使用配置文件的配置
		ports, err := tools.UnfoldPort(conf.Port.TargetPorts)
		if err != nil {
			fmt.Println("[-] 参数错误")
			return nil, err
		}
		svc.CmdArgs.TargetPorts = ports
	default:
		ports, err := tools.UnfoldPort(targetPorts)
		if err != nil {
			fmt.Println("[-] 参数错误")
			return nil, err
		}
		svc.CmdArgs.TargetPorts = ports
	}

	timeout, _ := cmd.Flags().GetInt("timeout")
	switch timeout {
	case 0:
		svc.CmdArgs.Timeout = conf.Port.Timeout
	default:
		svc.CmdArgs.Timeout = timeout
	}

	thread, _ := cmd.Flags().GetInt("thread")
	switch thread {
	case 0:
		svc.CmdArgs.Thread = conf.Port.Thread
	default:
		svc.CmdArgs.Thread = thread
	}

	retry, _ := cmd.Flags().GetInt("retry")
	switch retry {
	case 0:
		svc.CmdArgs.Retry = conf.Port.Retry
	default:
		svc.CmdArgs.Retry = retry
	}

	fingerprintFile, _ := cmd.Flags().GetString("fingerprint-file")
	switch fingerprintFile {
	case "":
		svc.CmdArgs.FingerprintFile = conf.Port.FingerprintFile
	default:
		svc.CmdArgs.FingerprintFile = fingerprintFile
	}

	outFile, _ := cmd.Flags().GetString("out-file")
	switch outFile {
	case "":
		svc.CmdArgs.OutFileName = conf.Port.OutFile
	default:
		svc.CmdArgs.OutFileName = outFile
	}

	hostsFile, _ := cmd.Flags().GetString("target-file")
	switch hostsFile {
	case "":
		if conf.Port.TargetFile == "" {
			return &svc.CmdArgs, nil
		}
		ipsData, err := tools.GetFile2Strings(conf.Port.TargetFile)
		if err != nil {
			fmt.Println("[-] 解析目标文件失败")
			return nil, err
		}
		ips, err := tools.UnfoldIPs(ipsData)
		if err != nil {
			fmt.Println("[-] 参数错误")
			return nil, err
		}
		tools.Shuffle(*ips)
		svc.CmdArgs.TargetIPs = ips
	default:
		ipsData, err := tools.GetFile2Strings(hostsFile)
		if err != nil {
			fmt.Println("[-] 解析目标文件失败")
			return nil, err
		}
		ips, err := tools.UnfoldIPs(ipsData)
		if err != nil {
			fmt.Println("[-] 参数错误")
			return nil, err
		}
		tools.Shuffle(*ips)
		svc.CmdArgs.TargetIPs = ips
	}

	return &svc.CmdArgs, nil
}

func initPackageArgsWorker(portCh <-chan string, packageArgsCh chan<- PackageArgs, wg *sync.WaitGroup, args PortCmdArgs, probes *[]model.Probe) {
	targetIPs := *args.TargetIPs
	for po := range portCh {
		for i := 0; i < len(targetIPs); i++ {
			packageArgs := new(PackageArgs)
			packageArgs.Protocol = args.Protocol
			packageArgs.TargetIP = targetIPs[i]
			packageArgs.TargetPort = po
			for _, probe := range *probes {
				if tools.IncludePort(po, probe.Ports) && (strings.ToLower(args.Protocol) == strings.ToLower(probe.Protocol)) {
					packageArgs.Probe = probe
					packageArgsCh <- *packageArgs
					continue
				}
				if tools.IncludePort(po, probe.SSLPorts) && (strings.ToLower(args.Protocol) == strings.ToLower(probe.Protocol)) {
					packageArgs.Probe = probe
					packageArgsCh <- *packageArgs
					continue
				}
				if probe.Name == "NULL" {
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

func (svc *Port) InitPackageArgs(packageArgsCh chan PackageArgs, mainWG *sync.WaitGroup, probes *[]model.Probe) {
	defer mainWG.Done()
	defer close(packageArgsCh)

	svc.logger.Debug("[+] 初始化扫描参数")
	targetPorts := *svc.CmdArgs.TargetPorts
	portCh := make(chan string, svc.CmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(portCh); i++ {
		go initPackageArgsWorker(portCh, packageArgsCh, &wg, svc.CmdArgs, probes)
	}

	for i := 0; i < len(targetPorts); i++ {
		wg.Add(1)
		portCh <- targetPorts[i]
	}

	wg.Wait()
	close(portCh)
}

func sendPackageWorker(argsCh <-chan PackageArgs, scanResultCh chan<- ScanResult, wg *sync.WaitGroup, cmdArgs PortCmdArgs) {
	for args := range argsCh {
		scanResult := ScanResult{
			IP:       args.TargetIP,
			Port:     args.TargetPort,
			Protocol: args.Protocol,
			Retry:    0,
		}
		address := fmt.Sprintf("%s:%s", args.TargetIP, args.TargetPort)

		// dialer := net.Dialer{Timeout: time.Duration(p.CmdArgs.Timeout) * time.Second}
		conn, err := net.DialTimeout(args.Protocol, address, time.Duration(cmdArgs.Timeout)*time.Second)
		if err != nil {
			scanResult.Error = err
			scanResult.Retry += 1
			scanResultCh <- scanResult
			wg.Done()
			continue
		}

		scanResult.State = true

		// 需要发送的数据
		if len(args.Probe.Data) > 0 {
			_ = conn.SetWriteDeadline(time.Now().Add(time.Duration(cmdArgs.Timeout) * time.Second))
			_, err = conn.Write(args.Probe.Data)
			if err != nil {
				scanResult.Error = err
				scanResult.Retry += 1
				scanResultCh <- scanResult
				continue
			}
		}

		_ = conn.SetReadDeadline(time.Now().Add(time.Duration(cmdArgs.Timeout) * time.Second))
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

func (svc *Port) SendPackage(packageArgsCh <-chan PackageArgs, scanResultCh chan ScanResult, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(scanResultCh)

	svc.logger.Debug("[+] 发送数据包")
	argsCh := make(chan PackageArgs, svc.CmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(argsCh); i++ {
		go sendPackageWorker(argsCh, scanResultCh, &wg, svc.CmdArgs)
	}

	for args := range packageArgsCh {
		wg.Add(1)
		argsCh <- args
	}

	wg.Wait()
	close(argsCh)
}

// 指纹识别逻辑
func fingerprintRecognitionLogic(matchs []*model.Match, resultCh chan<- Result, scanResult ScanResult, i int, probes *[]model.Probe) *Result {
	result := Result{
		IP:       scanResult.IP,
		Port:     scanResult.Port,
		Protocol: scanResult.Protocol,
		Retry:    i,
		Banner:   string(scanResult.Response),
	}

	// if strings.Contains(string(scanResult.Response), "\\u") {
	// 	// 转换
	// 	result.Banner = tools.Unicode2UTF8(result.Banner)
	// }

	switch scanResult.State {
	case false:
		result.State = "Close"
	default:
		result.State = "Open"
	}

	for _, match := range matchs {
		if match.PatternCompiled.MatcherString(result.Banner, 0).Matches() {
			// 方便格式化输出
			result.ServerType = match.Service
			result.IsSoft = match.IsSoft
			switch strings.Contains(match.VersionInfo, "$") {
			case false:
				var (
					infos []string
					infoS string
				)
				if strings.Contains(match.VersionInfo, "p/") {
					infos = strings.Split(match.VersionInfo, "p/")
				}
				switch len(infos) {
				case 0:
				default:
					infoS = strings.Split(infos[1], "/")[0]
					result.Version = infoS
					resultCh <- result
				}
			default:
				m, err := regexp.Compile(match.Pattern)
				if err != nil {
					continue
				}
				version := m.ReplaceAllString(string(scanResult.Response), match.VersionInfo)
				// version = match.PatternCompiled.ReplaceAllString()
				if strings.Contains(match.VersionInfo, "$I") {
					version = match.VersionInfo
				}
				var (
					infos        []string
					infoS, infoV string
				)
				if strings.Contains(version, "p/") {
					infos = strings.Split(version, "p/")
				}

				switch len(infos) {
				case 0:
				case 1:
					infoS = strings.Split(infos[1], "/")[0]
				default:
					infoS = strings.Split(infos[1], "/")[0]
					if strings.Contains(infos[1], "/ v/") {
						infos = strings.Split(infos[1], "/ v/")
						infoS = infos[0]
						infoV = strings.Split(infos[1], "/")[0]
					}
				}
				result.Version = fmt.Sprintf("%s %s", infoS, infoV)
				resultCh <- result
			}
		}
	}

	// 如果上面没有找到,则全局搜索
	if result.ServerType == "" && result.Version == "" {
		for _, probe := range *probes {
			if strings.ToLower(result.Protocol) == strings.ToLower(probe.Protocol) {
				for _, match := range probe.Matchs {
					if match.PatternCompiled.Matcher(scanResult.Response, 0).Matches() {
						// 方便格式化输出
						result.ServerType = match.Service
						result.IsSoft = match.IsSoft
						switch strings.Contains(match.VersionInfo, "$") {
						case false:
							m, err := regexp.Compile(match.Pattern)
							if err != nil {
								continue
							}
							version := m.ReplaceAllString(string(scanResult.Response), match.VersionInfo)
							var (
								infos []string
								infoS string
							)
							if strings.Contains(version, "p/") {
								infos = strings.Split(version, "p/")
							}
							switch len(infos) {
							case 0:
							default:
								infoS = strings.Split(infos[1], "/")[0]
								result.Version = infoS
								resultCh <- result
							}
						default:
							m, err := regexp.Compile(match.Pattern)
							if err != nil {
								continue
							}
							version := m.ReplaceAllString(string(scanResult.Response), match.VersionInfo)
							var (
								infos        []string
								infoS, infoV string
							)
							if strings.Contains(version, "p/") {
								infos = strings.Split(version, "p/")
							}

							switch len(infos) {
							case 0:
							case 1:
								infoS = strings.Split(infos[1], "/")[0]
							default:
								infoS = strings.Split(infos[1], "/")[0]
								if strings.Contains(infos[1], "/ v/") {
									infos = strings.Split(infos[1], "/ v/")
									infoS = infos[0]
									infoV = strings.Split(infos[1], "/")[0]
								}
							}
							result.Version = fmt.Sprintf("%s %s", infoS, infoV)
							resultCh <- result
						}
					}
				}
			}
		}
	}

	return &result
}

// 指纹识别线程池
func fingerprintRecognitionWorker(scanCh <-chan ScanResult, resultCh chan<- Result, wg *sync.WaitGroup, args PortCmdArgs, probes *[]model.Probe) {
	for scan := range scanCh {
		result := new(Result)

		for i := 0; i < args.Retry; i++ {
			if (result.Version == "" || result.Version == " ") && result.Retry < args.Retry {
				result = fingerprintRecognitionLogic(scan.Probe.Matchs, resultCh, scan, i, probes)
			}
		}

		resultCh <- *result
		wg.Done()
	}
}

func (svc *Port) FingerprintRecognition(scanResultCh <-chan ScanResult, resultCh chan Result, mainWG *sync.WaitGroup, probes *[]model.Probe) {
	defer mainWG.Done()
	defer close(resultCh)

	svc.logger.Debug("[+] 指纹识别")
	scanCh := make(chan ScanResult, svc.CmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(scanCh); i++ {
		go fingerprintRecognitionWorker(scanCh, resultCh, &wg, svc.CmdArgs, probes)
	}

	for res := range scanResultCh {
		// 过滤掉端口关闭的包
		if res.State != false {
			wg.Add(1)
			scanCh <- res
		}
	}

	wg.Wait()
	close(scanCh)
}

// 去重
func chanRemoveSliceMap(res Result, hashMap *sync.Map) {
	// p.logger.Info("[+] 去重", zap.Any("result", res))
	key := fmt.Sprintf("%s:%s", res.IP, res.Port)

	val, ok := hashMap.Load(key)

	if ok {
		if res.ServerType != "" && val.(Result).ServerType != "" {
			if len(res.ServerType) > len(val.(Result).ServerType) {
				hashMap.Store(key, res)
			}
		} else if res.ServerType != "" {
			hashMap.Store(key, res)
		}
		if res.Version != "" && val.(Result).Version != "" {
			if len(res.Version) > len(val.(Result).Version) {
				hashMap.Store(key, res)
			}
		} else if res.Version != "" {
			hashMap.Store(key, res)
		}
		if (res.Version == " " && val.(Result).Version == " ") || (res.Version == "" && val.(Result).Version == "") {
			if res.Retry > val.(Result).Retry {
				hashMap.Store(key, res)
			}
		}
	} else {
		hashMap.Store(key, res)
	}
}

func (svc *Port) OutputPrinting(resultCh <-chan Result) {
	var (
		hashMap   sync.Map
		id, oldID int
	)

	fmt.Printf("%-5s%-20s%-10s%-20s%-50.45s%-5s\n", "ID", "IP", "Port", "Server", "Server Info", "retry")

	// 生成日志文件
	_, err := os.Stat(fmt.Sprintf("%s-%s", "old", svc.CmdArgs.OutFileName))
	if err == nil {
		// 如果文件存在
		_ = os.Remove(fmt.Sprintf("%s-%s", "old", svc.CmdArgs.OutFileName))
	}

	oldFile, _ := os.Create(fmt.Sprintf("%s-%s", "old", svc.CmdArgs.OutFileName))
	_, _ = oldFile.WriteString(fmt.Sprintf("%-5s%-20s%-10s%-20s%-50.45s%-100.95s\n", "ID", "IP", "Port", "Server", "Server Info", "Banner"))
	// 去重
	for res := range resultCh {
		oldID += 1
		_, _ = oldFile.WriteString(fmt.Sprintf("%-5s%-20s%-10s%-20s%-50.45s%-100.95s\n", strconv.Itoa(oldID), res.IP, res.Port, res.ServerType, res.Version, res.Banner))
		chanRemoveSliceMap(res, &hashMap)
	}

	// 删除日志文件
	_ = os.Remove(fmt.Sprintf("%s-%s", "old", svc.CmdArgs.OutFileName))

	_, err = os.Stat(svc.CmdArgs.OutFileName)
	if err == nil {
		// 如果文件存在
		_ = os.Remove(svc.CmdArgs.OutFileName)
	}

	file, _ := os.Create(svc.CmdArgs.OutFileName)

	_, _ = file.WriteString(fmt.Sprintf("%-5s%-20s%-10s%-20s%-50.45s%-100.95s\n", "ID", "IP", "Port", "Server", "Server Info", "Banner"))
	hashMap.Range(func(key, value interface{}) bool {
		id += 1
		fmt.Printf("%-5s%-20s%-10s%-20s%-50.45s%-1d\n", strconv.Itoa(id), value.(Result).IP, value.(Result).Port, value.(Result).ServerType, value.(Result).Version, value.(Result).Retry)
		_, _ = file.WriteString(fmt.Sprintf("%-5s%-20s%-10s%-20s%-50.45s\n", strconv.Itoa(id), value.(Result).IP, value.(Result).Port, value.(Result).ServerType, value.(Result).Version))
		return true
	})

	_ = file.Close()
}
