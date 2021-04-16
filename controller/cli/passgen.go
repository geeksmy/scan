package cli

import (
	"fmt"
	"sync"
	"time"

	"scan/internal/service"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

type PassGen struct {
	cmd    *cobra.Command
	logger *zap.Logger
}

func NewPassGen(cmd *cobra.Command, logger *zap.Logger) *PassGen {
	return &PassGen{
		cmd:    cmd,
		logger: logger,
	}
}

func (p *PassGen) PassGenMain() error {
	start := time.Now()

	svc := service.NewPassGen(p.logger)

	_, err := svc.InitCmdArgs(p.cmd)
	if err != nil {
		return err
	}

	var mainWG sync.WaitGroup
	passwordCh := make(chan string, 10000)

	mainWG.Add(1)
	go svc.GeneratePass(passwordCh, &mainWG)

	svc.OutFile(passwordCh)

	mainWG.Wait()
	elapsed := time.Since(start)
	fmt.Println("耗时 ", elapsed)
	return nil
}