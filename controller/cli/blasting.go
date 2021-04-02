package cli

import (
	"fmt"
	"os"
	"sync"
	"time"

	"scan/config"
	"scan/pkg/tools"
	"scan/pkg/tools/blasting"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type BlastingResult struct {
	IP       string
	Port     string
	Username string
	Password string
	Server   string
	Retry    int
}

type BlastingResponse struct {
	IP       string
	Username string
	Password string
	Port     string
	Service  string
}

type IPPortService struct {
	IP      string
	Port    string
	Service string
}

// 命令行参数
type BlastingCmdArgs struct {
	Targets     *[]string
	Users       *[]string
	Passwords   *[]string
	Delay       int
	Thread      int
	Timeout     int
	Retry       int
	ScanPort    bool
	Services    *[]string
	Port        string
	Path        string
	TomcatPath  string
	OutPut      string
	OutFileName string
}

type Blasting struct {
	cmd    *cobra.Command
	logger *zap.Logger

	BlastingCmdArgs BlastingCmdArgs
}

func NewBlasting(cmd *cobra.Command, logger *zap.Logger) *Blasting {
	return &Blasting{
		cmd:    cmd,
		logger: logger,
	}
}

func (b *Blasting) BlastingMain() error {
	start := time.Now()
	err := b.initArgs()
	if err != nil {
		return err
	}

	responseCh := make(chan BlastingResponse, len(*b.BlastingCmdArgs.Services)*len(*b.BlastingCmdArgs.Targets)*len(*b.BlastingCmdArgs.Users)*len(*b.BlastingCmdArgs.Passwords))
	resultCh := make(chan BlastingResult, b.BlastingCmdArgs.Thread)
	var mainWG sync.WaitGroup

	mainWG.Add(1)
	go b.generateArgsChMainLogic(responseCh, &mainWG)

	mainWG.Add(1)
	go b.blastingMainLogic(responseCh, resultCh, &mainWG)

	b.outputPrinting(resultCh)

	mainWG.Wait()
	elapsed := time.Since(start)
	fmt.Println("耗时 ", elapsed)
	return nil
}

// 初始化命令参数
func (b *Blasting) initArgs() error {
	conf := config.C
	// 设置默认配置参数

	// 如果不传命令行参数使用配置文件的配置
	targetFile, _ := b.cmd.Flags().GetString("target-host")
	switch targetFile {
	case "":
		targetStr, err := tools.GetFile2Strings(conf.Blasting.TargetHost)
		if err != nil {
			b.logger.Error("[-] initArgs -> 解析目标文件失败")
			return err
		}
		target, err := tools.UnfoldIPs(targetStr)
		if err != nil {
			return err
		}
		tools.Shuffle(*target)
		b.BlastingCmdArgs.Targets = target
	default:
		targetStr, err := tools.GetFile2Strings(targetFile)
		if err != nil {
			b.logger.Error("[-] initArgs -> 解析目标文件失败")
			return err
		}
		target, err := tools.UnfoldIPs(targetStr)
		if err != nil {
			return err
		}
		tools.Shuffle(*target)
		b.BlastingCmdArgs.Targets = target
	}

	userFile, _ := b.cmd.Flags().GetString("user-file")
	switch userFile {
	case "":
		userStr, err := tools.GetFile2Strings(conf.Blasting.UserFile)
		if err != nil {
			b.logger.Error("[-] initArgs -> 解析目标文件失败")
			return err
		}
		b.BlastingCmdArgs.Users = &userStr
	default:
		userStr, err := tools.GetFile2Strings(userFile)
		if err != nil {
			b.logger.Error("[-] initArgs -> 解析目标文件失败")
			return err
		}
		b.BlastingCmdArgs.Users = &userStr
	}

	passFile, _ := b.cmd.Flags().GetString("pass-file")
	switch passFile {
	case "":
		passStr, err := tools.GetFile2Strings(conf.Blasting.PassFile)
		if err != nil {
			b.logger.Error("[-] initArgs -> 解析目标文件失败")
			return err
		}
		b.BlastingCmdArgs.Passwords = &passStr
	default:
		passStr, err := tools.GetFile2Strings(passFile)
		if err != nil {
			b.logger.Error("[-] initArgs -> 解析目标文件失败")
			return err
		}
		b.BlastingCmdArgs.Passwords = &passStr
	}

	delay, _ := b.cmd.Flags().GetInt("delay")
	switch delay {
	case 0:
		b.BlastingCmdArgs.Delay = conf.Blasting.Delay
	default:
		b.BlastingCmdArgs.Delay = delay
	}

	if b.BlastingCmdArgs.Delay <= 1 {
		b.BlastingCmdArgs.Delay = 1
	}

	thread, _ := b.cmd.Flags().GetInt("thread")
	switch thread {
	case 0:
		b.BlastingCmdArgs.Thread = conf.Blasting.Thread
	default:
		b.BlastingCmdArgs.Thread = thread
	}

	if b.BlastingCmdArgs.Thread <= 1 {
		b.BlastingCmdArgs.Thread = 1
	}

	timeout, _ := b.cmd.Flags().GetInt("timeout")
	switch timeout {
	case 0:
		b.BlastingCmdArgs.Timeout = conf.Blasting.Timeout
	default:
		b.BlastingCmdArgs.Timeout = timeout
	}

	if b.BlastingCmdArgs.Timeout <= 1 {
		b.BlastingCmdArgs.Timeout = 1
	}

	retry, _ := b.cmd.Flags().GetInt("retry")
	switch retry {
	case 0:
		b.BlastingCmdArgs.Retry = conf.Blasting.Retry
	default:
		b.BlastingCmdArgs.Retry = retry
	}

	if b.BlastingCmdArgs.Retry <= 1 {
		b.BlastingCmdArgs.Retry = 1
	}

	scanPort, _ := b.cmd.Flags().GetBool("scan-port")
	switch scanPort {
	case false:
		b.BlastingCmdArgs.ScanPort = false
	default:
		b.BlastingCmdArgs.ScanPort = true
	}

	services, _ := b.cmd.Flags().GetStringArray("services")
	switch len(services) {
	case 0:
		switch len(conf.Blasting.Services) {
		case 1:
			port, _ := b.cmd.Flags().GetString("port")
			b.BlastingCmdArgs.Port = port
			b.BlastingCmdArgs.Services = &conf.Blasting.Services
		default:
			b.BlastingCmdArgs.Services = &conf.Blasting.Services
		}
	case 1:
		s := tools.String2strings(services[0])
		b.BlastingCmdArgs.Services = &s
		port, _ := b.cmd.Flags().GetString("port")
		b.BlastingCmdArgs.Port = port
	default:
		b.BlastingCmdArgs.Services = &services
	}

	path, _ := b.cmd.Flags().GetString("path")
	switch path {
	case "":
		b.BlastingCmdArgs.Path = conf.Blasting.Path
	default:
		b.BlastingCmdArgs.Path = path
	}

	tomcatPath, _ := b.cmd.Flags().GetString("tomcat-path")
	switch tomcatPath {
	case "":
		b.BlastingCmdArgs.TomcatPath = conf.Blasting.TomcatPath
	default:
		b.BlastingCmdArgs.TomcatPath = tomcatPath
	}

	outFile, _ := b.cmd.Flags().GetString("out-file")
	switch outFile {
	case "":
		b.BlastingCmdArgs.OutPut = "print"
	default:
		b.BlastingCmdArgs.OutPut = "file"
		b.BlastingCmdArgs.OutFileName = outFile
	}

	return nil
}

// 生成 IP Port Service Ch
func generateIPPortService(ch <-chan string, ipCh chan<- IPPortService, targets []string, port string, wg *sync.WaitGroup) {
	ipPortService := new(IPPortService)
	for c := range ch {
		ipPortService.Service = c
		switch c {
		case "mssql":
			ipPortService.Port = "1443"
		case "ssh":
			ipPortService.Port = "22"
		case "ftp":
			ipPortService.Port = "21"
		case "mysql":
			ipPortService.Port = "3306"
		case "redis":
			ipPortService.Port = "6379"
		case "postgresql":
			ipPortService.Service = "postgres"
			ipPortService.Port = "5432"
		case "oracle":
			ipPortService.Port = "1521"
		case "http_basic":
			ipPortService.Port = "80"
		case "tomcat":
			ipPortService.Port = "8080"
		case "telnet":
			ipPortService.Port = "23"
		}
		if port != "" {
			ipPortService.Port = port
		}
		for i := 0; i < len(targets); i++ {
			ipPortService.IP = targets[i]
			ipCh <- *ipPortService
		}
		wg.Done()
	}
}

func generateIPPortServiceMainLogic(services, targets []string, port string, ipCh chan IPPortService) {
	ch := make(chan string, len(services))
	var wg sync.WaitGroup

	for i := 0; i < cap(ch); i++ {
		go generateIPPortService(ch, ipCh, targets, port, &wg)
	}

	for i := 0; i < len(services); i++ {
		wg.Add(1)
		ch <- services[i]
	}

	wg.Wait()
	close(ch)
	close(ipCh)
}

// 生成参数 chan 主逻辑线程池
func (b *Blasting) generateArgsChMainLogicWorker(ipCh <-chan IPPortService, responseCh chan<- BlastingResponse, wg *sync.WaitGroup) {
	response := new(BlastingResponse)
	users := *b.BlastingCmdArgs.Users
	passwords := *b.BlastingCmdArgs.Passwords
	for c := range ipCh {
		response.Service = c.Service
		response.IP = c.IP
		response.Port = c.Port

		switch c.Service {
		case "redis":
			for j := 0; j < len(passwords); j++ {
				response.Password = passwords[j]
				responseCh <- *response
			}
		default:
			for i := 0; i < len(users); i++ {
				for j := 0; j < len(passwords); j++ {
					response.Username = users[i]
					response.Password = passwords[j]
					responseCh <- *response
				}
			}
		}

		wg.Done()
	}
}

// 生成参数 chan 主逻辑
func (b *Blasting) generateArgsChMainLogic(responseCh chan BlastingResponse, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(responseCh)

	ipCh := make(chan IPPortService, len(*b.BlastingCmdArgs.Services)*len(*b.BlastingCmdArgs.Targets))
	generateIPPortServiceMainLogic(*b.BlastingCmdArgs.Services, *b.BlastingCmdArgs.Targets, b.BlastingCmdArgs.Port, ipCh)

	ipPortServiceCh := make(chan IPPortService, b.BlastingCmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(ipPortServiceCh); i++ {
		go b.generateArgsChMainLogicWorker(ipPortServiceCh, responseCh, &wg)
	}

	for i := range ipCh {
		wg.Add(1)
		ipPortServiceCh <- i
	}

	wg.Wait()
	close(ipPortServiceCh)
}

// 密码爆破主逻辑线程池
func (b *Blasting) blastingMainLogicWorker(ch <-chan BlastingResponse, resultCh chan<- BlastingResult, wg *sync.WaitGroup) {
	for response := range ch {
		result := BlastingResult{
			Server: response.Service,
			Retry:  0,
		}
		switch response.Service {
		case "mssql":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Username == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					dataSourceName := fmt.Sprintf("server=%s;user id=%s;password=%s", response.IP, response.Username, response.Password)
					if blasting.NewGormConnMssql(response.Service, dataSourceName) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "ssh":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Username == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnSSH(addr, response.Username, response.Password, b.BlastingCmdArgs.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "ftp":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Username == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					// 匿名登录
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnFTP(addr, "anonymous", "", b.BlastingCmdArgs.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = "anonymous"
						result.Password = ""
						resultCh <- result
					}
					if blasting.NewConnFTP(addr, response.Username, response.Password, b.BlastingCmdArgs.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "mysql":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Username == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql", response.Username, response.Password, response.IP, response.Port)
					if blasting.NewXormConnMysql(response.Service, dataSourceName) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "redis":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Password == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnRedis(addr, "") {
						result.IP = response.IP
						result.Port = response.Port
						result.Password = " "
						resultCh <- result
					}
					if blasting.NewConnRedis(addr, response.Password) {
						result.IP = response.IP
						result.Port = response.Port
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "postgres":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Username == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					result.Server = "postgresql"
					dataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable", response.IP, response.Port, response.Username, response.Password)
					if blasting.NewConnPgSql(response.Service, dataSourceName) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "oracle":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Username == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					dataSourceName := fmt.Sprintf("%s/%s@%s:%s/ORCL", response.Username, response.Password, response.IP, response.Port)
					if blasting.NewConnOracle("oci8", dataSourceName) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "http_basic":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Username == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnHttpBasic(addr, response.Username, response.Password, b.BlastingCmdArgs.Path, b.BlastingCmdArgs.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "tomcat":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Username == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnTomcat(addr, response.Username, response.Password, b.BlastingCmdArgs.TomcatPath, b.BlastingCmdArgs.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "telnet":
			for i := 0; i < b.BlastingCmdArgs.Retry; i++ {
				if result.Username == "" && result.Retry < b.BlastingCmdArgs.Retry {
					result.Retry += 1
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnTelnet(addr, response.Username, response.Password, b.BlastingCmdArgs.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		}

		wg.Done()
	}
}

// 密码爆破主逻辑
func (b *Blasting) blastingMainLogic(responseCh <-chan BlastingResponse, resultCh chan BlastingResult, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(resultCh)

	ch := make(chan BlastingResponse, b.BlastingCmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(ch); i++ {
		go b.blastingMainLogicWorker(ch, resultCh, &wg)
	}

	for response := range responseCh {
		wg.Add(1)
		ch <- response
	}

	wg.Wait()
	close(ch)
}

// 输出文件
func (b *Blasting) outFile(res BlastingResult, file *os.File) {
	_, _ = file.WriteString(fmt.Sprintf("%-20s%-10s%-20s%-20s%-20s\n", res.IP, res.Port, res.Server, res.Username, res.Password))
}

// 输出屏幕
func (b *Blasting) outCmd(res BlastingResult) {
	fmt.Printf("%-20s%-10s%-20s%-20s%-20s\n", res.IP, res.Port, res.Server, res.Username, res.Password)
}

// 输出打印
func (b *Blasting) outputPrinting(resultCh <-chan BlastingResult) {

	_, err := os.Stat(b.BlastingCmdArgs.OutFileName)
	if err == nil {
		// 如果文件存在
		_ = os.Remove(b.BlastingCmdArgs.OutFileName)
	}

	file, _ := os.Create(b.BlastingCmdArgs.OutFileName)

	switch b.BlastingCmdArgs.OutPut {
	case "file":
		_, _ = file.WriteString(fmt.Sprintf("%-20s%-10s%-20s%-20s%-20s\n", "ip", "port", "server", "user", "pass"))
		for res := range resultCh {
			b.outFile(res, file)
		}

	default:
		fmt.Printf("%-20s%-10s%-20s%-20s%-20s\n", "ip", "port", "server", "user", "pass")
		for res := range resultCh {
			b.outCmd(res)
		}
	}

	_ = file.Close()
}
