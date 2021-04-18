package service

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"scan/config"
	"scan/pkg/tools"
	"scan/pkg/util"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

type IntranetAliveSVC interface {
	/**
	 * InitArgs 初始化参数
	 * @param cmd 命令行传入参数结构体
	 */
	InitCmdArgs(cmd *cobra.Command) (*IntranetAliveCmdArgs, error)
	/**
	 * ICMPSendPackage 发送 icmp 包进行存活检测
	 */
	ICMPSendPackage(ipsCh chan string, mainWG *sync.WaitGroup)
	/**
	 * OutputPrinting 输出打印
	 * @param resultCh 返回参数 chan
	 */
	OutputPrinting(ipsCh <-chan string)
}

type IntranetAliveCmdArgs struct {
	Targets     []string
	Timeout     int
	Thread      int
	Retry       int
	Delay       int
	OutFileName string
}

type IntranetAlive struct {
	logger *zap.Logger

	Args IntranetAliveCmdArgs
}

func NewIntranetAlive(logger *zap.Logger) IntranetAliveSVC {
	return &IntranetAlive{
		logger: logger,
	}
}

func (svc *IntranetAlive) InitCmdArgs(cmd *cobra.Command) (*IntranetAliveCmdArgs, error) {
	conf := config.C

	target, _ := cmd.Flags().GetString("target")
	switch target {
	case "":
		svc.Args.Targets = tools.GenerateIntranet(util.TargetIPs, tools.String2strings(conf.IntranetAlive.Target))
	default:
		svc.Args.Targets = tools.GenerateIntranet(util.TargetIPs, tools.String2strings(target))
	}

	timeout, _ := cmd.Flags().GetInt("timeout")
	switch timeout {
	case 0:
		svc.Args.Timeout = conf.IntranetAlive.Timeout
	default:
		svc.Args.Timeout = timeout
	}

	thread, _ := cmd.Flags().GetInt("thread")
	switch thread {
	case 0:
		svc.Args.Thread = conf.IntranetAlive.Thread
	default:
		svc.Args.Thread = thread
	}

	retry, _ := cmd.Flags().GetInt("retry")
	switch retry {
	case 0:
		svc.Args.Retry = conf.IntranetAlive.Retry
	default:
		svc.Args.Retry = retry
	}

	delay, _ := cmd.Flags().GetFloat32("delay")
	switch delay {
	case 0:
		svc.Args.Delay = int(conf.IntranetAlive.Delay * 1000)
	default:
		svc.Args.Delay = int(delay * 1000)
	}

	outFile, _ := cmd.Flags().GetString("out-file")
	switch outFile {
	case "":
		svc.Args.OutFileName = conf.IntranetAlive.OutFile
	default:
		svc.Args.OutFileName = outFile
	}

	if target == "" && conf.IntranetAlive.Target == "" {
		return nil, errors.New("[-] target参数必填")
	}

	return &svc.Args, nil
}

func ICMPSendPackageWork(argsCh <-chan string, ipsCh chan<- string, wg *sync.WaitGroup, timeout, retry, delay int) {
	for args := range argsCh {
		time.Sleep(time.Duration(delay) * time.Millisecond)
		conn, err := net.DialTimeout("ip4:icmp", args, time.Duration(timeout)*time.Second)
		if err != nil {
			wg.Done()
			continue
		}

		s := []byte{0x08, 0x00, 0xdd, 0x39, 0x00, 0x00, 0x00, 0x00}
		data := []byte("hello icmp")
		s = append(s, data...)

		_, err = conn.Write(s)
		if err != nil {
			wg.Done()
			continue
		}

		buf := make([]byte, 512)
		err = conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		if err != nil {
			wg.Done()
			continue
		}
		n, err := conn.Read(buf)
		if err != nil {
			wg.Done()
			continue
		}

		if string(buf[28:n]) == string(data) {
			ipsCh <- args
		}

		wg.Done()
		conn.Close()
	}

}

func (svc *IntranetAlive) ICMPSendPackage(ipsCh chan string, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(ipsCh)

	var wg sync.WaitGroup
	argsCh := make(chan string, svc.Args.Thread)

	for i := 0; i < cap(argsCh); i++ {
		go ICMPSendPackageWork(argsCh, ipsCh, &wg, svc.Args.Timeout, svc.Args.Retry, svc.Args.Delay)
	}

	for _, target := range svc.Args.Targets {
		wg.Add(1)
		argsCh <- target
	}

	wg.Wait()
	close(argsCh)
}

func ipsChRemoveSliceMap(ip string, ipsMap *sync.Map) {
	s := strings.Split(ip, ".")
	key := fmt.Sprintf("%s.%s.%s.0/24", s[0], s[1], s[2])

	_, ok := ipsMap.Load(key)

	if !ok {
		ipsMap.Store(key, ip)
	}
}

func (svc *IntranetAlive) OutputPrinting(ipsCh <-chan string) {
	var (
		ipsMap    sync.Map
		id, oldId int
	)

	fmt.Printf("%-5s%-20s\n", "ID", "IP")

	// 生成日志文件
	_, err := os.Stat(fmt.Sprintf("%s-%s", "old", svc.Args.OutFileName))
	if err == nil {
		// 如果文件存在
		_ = os.Remove(fmt.Sprintf("%s-%s", "old", svc.Args.OutFileName))
	}

	oldFile, _ := os.Create(fmt.Sprintf("%s-%s", "old", svc.Args.OutFileName))
	_, _ = oldFile.WriteString(fmt.Sprintf("%-5s,%-20s\n", "ID", "IP"))

	// 去重
	for ip := range ipsCh {
		oldId += 1
		_, _ = oldFile.WriteString(fmt.Sprintf("%-5s%-20s", strconv.Itoa(oldId), ip))
		go ipsChRemoveSliceMap(ip, &ipsMap)
	}

	_ = os.Remove(fmt.Sprintf("%s-%s", "old", svc.Args.OutFileName))

	_, err = os.Stat(svc.Args.OutFileName)
	if err == nil {
		// 如果文件存在
		_ = os.Remove(svc.Args.OutFileName)
	}

	file, _ := os.Create(svc.Args.OutFileName)
	_, _ = file.WriteString(fmt.Sprintf("%-5s,%-20s\n", "ID", "IP"))
	ipsMap.Range(func(key, value interface{}) bool {
		id += 1
		fmt.Printf("%-5s%-20s\n", strconv.Itoa(id), key)
		_, _ = file.WriteString(fmt.Sprintf("%-5s%-20s", strconv.Itoa(id), key))
		return true
	})

}
