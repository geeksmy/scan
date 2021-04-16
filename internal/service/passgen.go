package service

import (
	"fmt"
	"os"
	"sync"

	"scan/config"
	"scan/pkg/tools"
	"scan/pkg/util"

	"github.com/geeksmy/cobra"
	"go.uber.org/zap"
)

type PassGenSVC interface {
	/**
	 * InitArgs 初始化参数
	 * @param cmd 命令行传入参数结构体
	 */
	InitCmdArgs(cmd *cobra.Command) (*PassGenCmdArgs, error)
	/**
	 * GeneratePass 生成密码
	 * @param passwordCh 密码管道
	 */
	GeneratePass(passwordCh chan string, mainWG *sync.WaitGroup)
	/**
	 * OutFile 生成密码文件
	 * @param passwordCh 密码管道
	 */
	OutFile(passwordCh <-chan string)
}

type PassGenCmdArgs struct {
	Year        []string
	DomainName  []string
	Domain      []string
	Device      []string
	SpecialLen  int
	OutFileName string
}

type PassGen struct {
	logger *zap.Logger

	Args PassGenCmdArgs
}

func NewPassGen(logger *zap.Logger) PassGenSVC {
	return &PassGen{
		logger: logger,
	}
}

func (svc *PassGen) InitCmdArgs(cmd *cobra.Command) (*PassGenCmdArgs, error) {
	conf := config.C

	year, _ := cmd.Flags().GetString("year")
	switch year {
	case "":
		svc.Args.Year = tools.String2strings(conf.PassGen.Year)
	default:
		svc.Args.Year = tools.String2strings(year)
	}

	domainName, _ := cmd.Flags().GetString("domain-name")
	switch domainName {
	case "":
		svc.Args.DomainName = tools.String2strings(conf.PassGen.DomainName)
	default:
		svc.Args.DomainName = tools.String2strings(domainName)
	}

	domain, _ := cmd.Flags().GetString("domain")
	switch domain {
	case "":
		svc.Args.Domain = tools.String2strings(conf.PassGen.Domain)
	default:
		svc.Args.Domain = tools.String2strings(domain)
	}

	device, _ := cmd.Flags().GetString("device")
	switch device {
	case "":
		svc.Args.Device = tools.String2strings(conf.PassGen.Device)
	default:
		svc.Args.Device = tools.String2strings(device)
	}

	fileName, _ := cmd.Flags().GetString("out-file")
	switch fileName {
	case "":
		svc.Args.OutFileName = conf.PassGen.OutFile
	default:
		svc.Args.OutFileName = fileName
	}

	length, _ := cmd.Flags().GetInt("length")
	switch length {
	case 0:
		svc.Args.SpecialLen = conf.PassGen.Length
	default:
		svc.Args.SpecialLen = length
	}

	return &svc.Args, nil
}

func (svc *PassGen) GeneratePass(passwordCh chan string, mainWG *sync.WaitGroup) {
	defer close(passwordCh)
	defer mainWG.Done()

	var (
		genWG                                    sync.WaitGroup
		isYear, isDomain, isDomainName, isDevice bool
	)

	switch len(svc.Args.Year) {
	case 0:
		isYear = false
	case 1:
		if svc.Args.Year[0] != "" {
			isYear = true
		}
	default:
		isYear = true
	}

	switch len(svc.Args.DomainName) {
	case 0:
		isDomainName = false
	case 1:
		if svc.Args.DomainName[0] != "" {
			isDomainName = true
		}
	default:
		isDomainName = true
	}

	switch len(svc.Args.Domain) {
	case 0:
		isDomain = false
	case 1:
		if svc.Args.Domain[0] != "" {
			isDomain = true
		}
	default:
		isDomain = true
	}

	switch len(svc.Args.Device) {
	case 0:
		isDevice = false
	case 1:
		if svc.Args.Device[0] != "" {
			isDevice = true
		}
	default:
		isDevice = true
	}

	if isYear {
		// 0x01 基于 "月份 + 年份 英文简写" 设置的密码
		genWG.Add(1)
		go tools.Generate(util.Month, svc.Args.Year, util.Special, svc.Args.SpecialLen, passwordCh, &genWG)

		// 0x01 基于 "月份 + 年份 英文全称" 设置的密码
		genWG.Add(1)
		go tools.Generate(util.MonthsAbbreviation, svc.Args.Year, util.Special, svc.Args.SpecialLen, passwordCh, &genWG)

		// 0x02 季节的英文全称
		genWG.Add(1)
		go tools.Generate(util.Season, svc.Args.Year, util.Special, svc.Args.SpecialLen, passwordCh, &genWG)

		// 0x03 季度的英文简写
		genWG.Add(1)
		go tools.Generate(util.Quarterly, svc.Args.Year, util.Special, svc.Args.SpecialLen, passwordCh, &genWG)
	}

	if isDomainName {
		// 0x04 集团 / 公司名全称 或 谐音拼音
		var digital []string
		digital = util.Digital
		if isYear {
			for i := 0; i < len(svc.Args.Year); i++ {
				digital = append(digital, svc.Args.Year[i])
			}
		}
		genWG.Add(1)
		go tools.Generate(svc.Args.DomainName, digital, util.Special, svc.Args.SpecialLen, passwordCh, &genWG)
	}

	if isDevice {
		// 0x05 特定设备名
		genWG.Add(1)
		go tools.Generate(svc.Args.Device, util.Digital, util.Special, svc.Args.SpecialLen, passwordCh, &genWG)
	}

	if isDomain {
		for _, domain := range svc.Args.Domain {
			passwordCh <- domain
		}
	}

	genWG.Wait()
}

func (svc *PassGen) OutFile(passwordCh <-chan string) {
	_, err := os.Stat(svc.Args.OutFileName)
	if err == nil {
		// 如果文件存在
		_ = os.Remove(svc.Args.OutFileName)
	}

	file, _ := os.Create(svc.Args.OutFileName)

	for pass := range passwordCh {
		_, _ = file.WriteString(fmt.Sprintf("%s\n", pass))
	}
}
