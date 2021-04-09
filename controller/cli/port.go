package cli

import (
	"fmt"
	"sync"
	"time"

	"scan/internal/service"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type Port struct {
	cmd    *cobra.Command
	logger *zap.Logger
}

func NewPort(cmd *cobra.Command, logger *zap.Logger) *Port {
	return &Port{
		cmd:    cmd,
		logger: logger,
	}
}

func (p *Port) PortMain() error {
	start := time.Now()

	svc := service.NewPortSVC(p.logger)
	// 初始化命令参数
	args, err := svc.InitArgs(p.cmd)
	if err != nil {
		return err
	}

	// 初始化规则文件
	parse := service.NewParse(p.logger)
	probes, err := parse.ParsingNmapFingerprint(args.FingerprintFile)
	if err != nil {
		return err
	}

	// 初始化扫描参数
	var mainWG sync.WaitGroup
	packageArgsCh := make(chan service.PackageArgs, len(*args.TargetIPs)*len(*args.TargetPorts))
	scanResultCh := make(chan service.ScanResult, len(*args.TargetIPs)*len(*args.TargetPorts))
	resultCh := make(chan service.Result, len(*args.TargetIPs)*len(*args.TargetPorts))
	mainWG.Add(1)
	go svc.InitPackageArgs(packageArgsCh, &mainWG, probes)

	// 扫描
	mainWG.Add(1)
	go svc.SendPackage(packageArgsCh, scanResultCh, &mainWG)

	// 指纹识别
	mainWG.Add(1)
	go svc.FingerprintRecognition(scanResultCh, resultCh, &mainWG, probes)

	// 输出打印
	// mainWG.Add(1)
	svc.OutputPrinting(resultCh)

	mainWG.Wait()
	elapsed := time.Since(start)
	fmt.Println("耗时 ", elapsed)
	return nil
}
