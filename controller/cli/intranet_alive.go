package cli

import (
	"fmt"
	"sync"
	"time"

	"scan/internal/service"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

type IntranetAlive struct {
	cmd    *cobra.Command
	logger *zap.Logger
}

func NewIntranetAlive(cmd *cobra.Command, logger *zap.Logger) *IntranetAlive {
	return &IntranetAlive{
		cmd:    cmd,
		logger: logger,
	}
}

func (i *IntranetAlive) IntranetAliveMain() error {
	start := time.Now()

	svc := service.NewIntranetAlive(i.logger)
	args, err := svc.InitCmdArgs(i.cmd)
	if err != nil {
		return err
	}

	var mainWG sync.WaitGroup
	ipsCh := make(chan string, len(args.Targets)/2)

	mainWG.Add(1)
	go svc.ICMPSendPackage(ipsCh, &mainWG)

	svc.OutputPrinting(ipsCh)

	mainWG.Wait()
	elapsed := time.Since(start)
	fmt.Println("耗时 ", elapsed)
	return nil
}
