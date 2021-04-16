package cli

import (
	"fmt"
	"sync"
	"time"

	"scan/internal/service"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

type Blasting struct {
	cmd    *cobra.Command
	logger *zap.Logger
}

func NewBlasting(cmd *cobra.Command, logger *zap.Logger) *Blasting {
	return &Blasting{
		cmd:    cmd,
		logger: logger,
	}
}

func (b *Blasting) BlastingMain() error {
	start := time.Now()

	svc := service.NewBlastingSVC(b.logger)

	args, err := svc.InitArgs(b.cmd)
	if err != nil {
		return err
	}

	responseCh := make(chan service.BlastingRequest, len(*args.Services)*len(*args.Targets)*len(*args.Users)*len(*args.Passwords))
	resultCh := make(chan service.BlastingResult, args.Thread)
	var mainWG sync.WaitGroup

	mainWG.Add(1)
	go svc.GenerateArgsChMainLogic(responseCh, &mainWG)

	mainWG.Add(1)
	go svc.BlastingMainLogic(responseCh, resultCh, &mainWG)

	svc.OutputPrinting(resultCh)

	mainWG.Wait()
	elapsed := time.Since(start)
	fmt.Println("耗时 ", elapsed)
	return nil
}
