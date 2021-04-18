package service

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"scan/config"
	"scan/pkg/tools"
	"scan/pkg/tools/blasting"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

type BlastingSVC interface {
	/**
	 * InitArgs 初始化参数
	 * @param cmd 命令行传入参数结构体
	 */
	InitArgs(cmd *cobra.Command) (*BlastingCmdArgs, error)
	/**
	 * GenerateArgsChMainLogic 生成参数主逻辑
	 * @param requestCh 需要发送的参数 chan
	 */
	GenerateArgsChMainLogic(requestCh chan BlastingRequest, mainWG *sync.WaitGroup)
	/**
	 * BlastingMainLogic 密码爆破主逻辑
	 * @param requestCh 需要发送的参数 chan
	 * @param resultCh 返回参数 chan
	 */
	BlastingMainLogic(requestCh <-chan BlastingRequest, resultCh chan BlastingResult, mainWG *sync.WaitGroup)
	/**
	 * OutputPrinting 输出打印
	 * @param resultCh 返回参数 chan
	 */
	OutputPrinting(resultCh <-chan BlastingResult)
}

type BlastingResult struct {
	IP       string
	Port     string
	Username string
	Password string
	Server   string
	Retry    int
}

type BlastingRequest struct {
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
	OutFileName string
}

type Blasting struct {
	logger *zap.Logger

	BlastingCmdArgs BlastingCmdArgs
}

func NewBlastingSVC(logger *zap.Logger) BlastingSVC {
	return &Blasting{
		logger: logger,
	}
}

func (svc *Blasting) InitArgs(cmd *cobra.Command) (*BlastingCmdArgs, error) {
	conf := config.C
	// 设置默认配置参数

	// 如果不传命令行参数使用配置文件的配置
	targetFile, _ := cmd.Flags().GetString("target-host")
	switch targetFile {
	case "":
		targetStr, err := tools.GetFile2Strings(conf.Blasting.TargetHost)
		if err != nil {
			fmt.Println("[-] 解析目标文件失败")
			return nil, err
		}
		target, err := tools.UnfoldIPs(targetStr)
		if err != nil {
			return nil, err
		}
		tools.Shuffle(*target)
		svc.BlastingCmdArgs.Targets = target
	default:
		targetStr, err := tools.GetFile2Strings(targetFile)
		if err != nil {
			fmt.Println("[-] 解析目标文件失败")
			return nil, err
		}
		target, err := tools.UnfoldIPs(targetStr)
		if err != nil {
			return nil, err
		}
		tools.Shuffle(*target)
		svc.BlastingCmdArgs.Targets = target
	}

	userFile, _ := cmd.Flags().GetString("user-file")
	switch userFile {
	case "":
		userStr, err := tools.GetFile2Strings(conf.Blasting.UserFile)
		if err != nil {
			fmt.Println("[-] 解析用户文件失败")
			return nil, err
		}
		svc.BlastingCmdArgs.Users = &userStr
	default:
		userStr, err := tools.GetFile2Strings(userFile)
		if err != nil {
			fmt.Println("[-] 解析用户文件失败")
			return nil, err
		}
		svc.BlastingCmdArgs.Users = &userStr
	}

	passFile, _ := cmd.Flags().GetString("pass-file")
	switch passFile {
	case "":
		passStr, err := tools.GetFile2Strings(conf.Blasting.PassFile)
		if err != nil {
			fmt.Println("[-] 解析密码文件失败")
			return nil, err
		}
		svc.BlastingCmdArgs.Passwords = &passStr
	default:
		passStr, err := tools.GetFile2Strings(passFile)
		if err != nil {
			fmt.Println("[-] 解析密码文件失败")
			return nil, err
		}
		svc.BlastingCmdArgs.Passwords = &passStr
	}

	delay, _ := cmd.Flags().GetInt("delay")
	switch delay {
	case 0:
		svc.BlastingCmdArgs.Delay = conf.Blasting.Delay
	default:
		svc.BlastingCmdArgs.Delay = delay
	}

	thread, _ := cmd.Flags().GetInt("thread")
	switch thread {
	case 0:
		svc.BlastingCmdArgs.Thread = conf.Blasting.Thread
	default:
		svc.BlastingCmdArgs.Thread = thread
	}

	timeout, _ := cmd.Flags().GetInt("timeout")
	switch timeout {
	case 0:
		svc.BlastingCmdArgs.Timeout = conf.Blasting.Timeout
	default:
		svc.BlastingCmdArgs.Timeout = timeout
	}

	retry, _ := cmd.Flags().GetInt("retry")
	switch retry {
	case 0:
		svc.BlastingCmdArgs.Retry = conf.Blasting.Retry
	default:
		svc.BlastingCmdArgs.Retry = retry
	}

	scanPort, _ := cmd.Flags().GetBool("scan-port")
	switch scanPort {
	case false:
		svc.BlastingCmdArgs.ScanPort = false
	default:
		svc.BlastingCmdArgs.ScanPort = true
	}

	services, _ := cmd.Flags().GetStringArray("services")
	switch len(services) {
	case 0:
		switch len(conf.Blasting.Services) {
		case 1:
			port, _ := cmd.Flags().GetString("port")
			svc.BlastingCmdArgs.Port = port
			svc.BlastingCmdArgs.Services = &conf.Blasting.Services
		default:
			svc.BlastingCmdArgs.Services = &conf.Blasting.Services
		}
	case 1:
		s := tools.String2strings(services[0])
		svc.BlastingCmdArgs.Services = &s
		port, _ := cmd.Flags().GetString("port")
		svc.BlastingCmdArgs.Port = port
	default:
		svc.BlastingCmdArgs.Services = &services
	}

	path, _ := cmd.Flags().GetString("path")
	switch path {
	case "":
		svc.BlastingCmdArgs.Path = conf.Blasting.Path
	default:
		svc.BlastingCmdArgs.Path = path
	}

	tomcatPath, _ := cmd.Flags().GetString("tomcat-path")
	switch tomcatPath {
	case "":
		svc.BlastingCmdArgs.TomcatPath = conf.Blasting.TomcatPath
	default:
		svc.BlastingCmdArgs.TomcatPath = tomcatPath
	}

	outFile, _ := cmd.Flags().GetString("out-file")
	switch outFile {
	case "":
		svc.BlastingCmdArgs.OutFileName = conf.Blasting.OutFile
	default:
		svc.BlastingCmdArgs.OutFileName = outFile
	}

	return &svc.BlastingCmdArgs, nil
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
func generateArgsChMainLogicWorker(ipCh <-chan IPPortService, requestCh chan<- BlastingRequest, wg *sync.WaitGroup, args BlastingCmdArgs) {
	request := new(BlastingRequest)
	users := *args.Users
	passwords := *args.Passwords
	for c := range ipCh {
		request.Service = c.Service
		request.IP = c.IP
		request.Port = c.Port

		switch c.Service {
		case "redis":
			for j := 0; j < len(passwords); j++ {
				request.Password = passwords[j]
				requestCh <- *request
			}
		default:
			for i := 0; i < len(users); i++ {
				for j := 0; j < len(passwords); j++ {
					request.Username = users[i]
					request.Password = passwords[j]
					requestCh <- *request
				}
			}
		}

		wg.Done()
	}
}

func (svc *Blasting) GenerateArgsChMainLogic(requestCh chan BlastingRequest, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(requestCh)

	ipCh := make(chan IPPortService, len(*svc.BlastingCmdArgs.Services)*len(*svc.BlastingCmdArgs.Targets))
	generateIPPortServiceMainLogic(*svc.BlastingCmdArgs.Services, *svc.BlastingCmdArgs.Targets, svc.BlastingCmdArgs.Port, ipCh)

	ipPortServiceCh := make(chan IPPortService, svc.BlastingCmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(ipPortServiceCh); i++ {
		go generateArgsChMainLogicWorker(ipPortServiceCh, requestCh, &wg, svc.BlastingCmdArgs)
	}

	for i := range ipCh {
		wg.Add(1)
		ipPortServiceCh <- i
	}

	wg.Wait()
	close(ipPortServiceCh)
}

// 密码爆破主逻辑线程池
func blastingMainLogicWorker(ch <-chan BlastingRequest, resultCh chan<- BlastingResult, wg *sync.WaitGroup, args BlastingCmdArgs) {
	for response := range ch {
		result := BlastingResult{
			Server: response.Service,
			Retry:  0,
		}
		switch response.Service {
		case "mssql":
			for i := 0; i < args.Retry; i++ {
				if result.Username == "" && result.Retry < args.Retry {
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
			for i := 0; i < args.Retry; i++ {
				if result.Username == "" && result.Retry < args.Retry {
					result.Retry += 1
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnSSH(addr, response.Username, response.Password, args.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "ftp":
			for i := 0; i < args.Retry; i++ {
				if result.Username == "" && result.Retry < args.Retry {
					result.Retry += 1
					// 匿名登录
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnFTP(addr, "anonymous", "", args.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = "anonymous"
						result.Password = ""
						resultCh <- result
					}
					if blasting.NewConnFTP(addr, response.Username, response.Password, args.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "mysql":
			for i := 0; i < args.Retry; i++ {
				if result.Username == "" && result.Retry < args.Retry {
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
			for i := 0; i < args.Retry; i++ {
				if result.Password == "" && result.Retry < args.Retry {
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
			for i := 0; i < args.Retry; i++ {
				if result.Username == "" && result.Retry < args.Retry {
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
			for i := 0; i < args.Retry; i++ {
				if result.Username == "" && result.Retry < args.Retry {
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
			for i := 0; i < args.Retry; i++ {
				if result.Username == "" && result.Retry < args.Retry {
					result.Retry += 1
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnHttpBasic(addr, response.Username, response.Password, args.Path, args.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "tomcat":
			for i := 0; i < args.Retry; i++ {
				if result.Username == "" && result.Retry < args.Retry {
					result.Retry += 1
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnTomcat(addr, response.Username, response.Password, args.TomcatPath, args.Timeout) {
						result.IP = response.IP
						result.Port = response.Port
						result.Username = response.Username
						result.Password = response.Password
						resultCh <- result
					}
				}
			}
		case "telnet":
			for i := 0; i < args.Retry; i++ {
				if result.Username == "" && result.Retry < args.Retry {
					result.Retry += 1
					addr := fmt.Sprintf("%s:%s", response.IP, response.Port)
					if blasting.NewConnTelnet(addr, response.Username, response.Password, args.Timeout) {
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

func (svc *Blasting) BlastingMainLogic(requestCh <-chan BlastingRequest, resultCh chan BlastingResult, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(resultCh)

	ch := make(chan BlastingRequest, svc.BlastingCmdArgs.Thread)
	var wg sync.WaitGroup

	for i := 0; i < cap(ch); i++ {
		go blastingMainLogicWorker(ch, resultCh, &wg, svc.BlastingCmdArgs)
	}

	for request := range requestCh {
		wg.Add(1)
		ch <- request
	}

	wg.Wait()
	close(ch)
}

// 输出文件
func outFile(res BlastingResult, file *os.File, id int) {
	_, _ = file.WriteString(fmt.Sprintf("%-5s%-20s%-10s%-15s%-10s%-15s\n", strconv.Itoa(id), res.IP, res.Port, res.Server, res.Username, res.Password))
}

// 输出屏幕
func outCmd(res BlastingResult, id int) {
	fmt.Printf("%-5s%-20s%-10s%-15s%-10s%-15s\n", strconv.Itoa(id), res.IP, res.Port, res.Server, res.Username, res.Password)
}

func (svc *Blasting) OutputPrinting(resultCh <-chan BlastingResult) {
	var id int

	_, err := os.Stat(svc.BlastingCmdArgs.OutFileName)
	if err == nil {
		// 如果文件存在
		_ = os.Remove(svc.BlastingCmdArgs.OutFileName)
	}

	file, _ := os.Create(svc.BlastingCmdArgs.OutFileName)
	_, _ = file.WriteString(fmt.Sprintf("%-5s%-20s%-10s%-15s%-10s%-15s\n", "id", "ip", "port", "server", "user", "pass"))
	fmt.Printf("%-5s%-20s%-10s%-15s%-10s%-15s\n", "id", "ip", "port", "server", "user", "pass")

	for res := range resultCh {
		id += 1
		outFile(res, file, id)
		outCmd(res, id)
	}

	_ = file.Close()
}
