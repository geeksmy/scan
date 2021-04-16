package cli

import (
	"fmt"
	"sync"
	"time"

	"scan/internal/model"
	"scan/internal/service"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

type WebFingerprint struct {
	cmd    *cobra.Command
	logger *zap.Logger
}

func NewWebFingerprint(cmd *cobra.Command, logger *zap.Logger) *WebFingerprint {
	return &WebFingerprint{
		cmd:    cmd,
		logger: logger,
	}
}

func (w *WebFingerprint) WebFingerprintMain() error {
	start := time.Now()
	svc := service.NewWebFingerPrintSVC(w.logger)
	args, err := svc.InitCmdArgs(w.cmd)
	if err != nil {
		return err
	}

	err = svc.InitFingerPrintFile()
	if err != nil {
		return err
	}

	svc.InitRequestArgs()

	var mainWG sync.WaitGroup
	responses := make(chan service.WebFingerPrintResponse, len(*args.TargetIPs)*len(*args.TargetPorts))
	results := make(chan model.Web, len(*args.TargetIPs)*len(*args.TargetPorts))

	mainWG.Add(1)
	go svc.SendRequest(responses, &mainWG)

	mainWG.Add(1)
	go svc.IdentifyResponse(responses, results, &mainWG)

	svc.OutputPrinting(results)

	mainWG.Wait()
	elapsed := time.Since(start)
	fmt.Println("耗时 ", elapsed)
	return nil
}
