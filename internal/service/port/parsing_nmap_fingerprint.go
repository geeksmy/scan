package port

import (
	"errors"
	"io/ioutil"
	"strings"

	"scan/internal/model"
	"scan/internal/serialization"

	"go.uber.org/zap"
)

type Parse struct {
	logger *zap.Logger

	Probes  *[]model.Probe
	Exclude string
}

func NewParse(logger *zap.Logger) *Parse {
	return &Parse{
		logger: logger,
	}
}

// 从nmap-service-probes文件中解析并加载规则
func (p *Parse) ParsingNmapFingerprint(nmapFingerprintFile string) (*[]model.Probe, error) {
	p.logger.Debug("[+] 读取指纹文件 -> ", zap.String("文件名", nmapFingerprintFile))
	// 读取 nmap-service-probes 或自定义的规则文件
	fingerprintData, err := ioutil.ReadFile(nmapFingerprintFile)
	if err != nil {
		p.logger.Error("[-] ParsingNmapFingerprint -> 读取规则文件失败")
		return nil, err
	}

	// 解析规则文本获取Probe列表
	p.logger.Debug("[+] 解析指纹文件 -> ", zap.String("文件名", nmapFingerprintFile))
	if err := p.parseFingerprint2ProbesList(fingerprintData); err != nil {
		return nil, err
	}
	p.logger.Debug("[+] 解析指纹文件 -> ok")
	return p.Probes, nil

}

func (p *Parse) parseFingerprint2ProbesList(fingerprintData []byte) error {
	var (
		lines   []string
		s, flag int
	)

	linesTemp := strings.Split(string(fingerprintData), "\n")

	// 过滤注释和空行
	for i := 0; i < len(linesTemp); i++ {
		lineTemp := strings.TrimSpace(linesTemp[i])
		if lineTemp == "" || strings.HasPrefix(lineTemp, "#") {
			continue
		}
		lines = append(lines, lineTemp)
	}

	if len(lines) < 1 {
		p.logger.Error("[-] ParsingNmapFingerprint -> 请在规则文件添加规则")
		return errors.New("[-] ParsingNmapFingerprint -> 请在规则文件添加规则")
	}

	// 一份规则文件里最多只能有一个 Exclude 设置
	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "Exclude ") {
			s += 1
			flag = i
			// 取出规则文件中的 Exclude 行
			p.Exclude = lines[i][len("Exclude")+1:]
			lines = lines[1:]
		}
		if s > 1 {
			p.logger.Error("[-] ParsingNmapFingerprint -> 规则文件只能有一个Exclude 设置")
			return errors.New("[-] ParsingNmapFingerprint -> 规则文件只能有一个Exclude 设置")
		}
	}

	// 规则文件第一行必需是 Exclude
	if flag != 0 {
		p.logger.Error("[-] ParsingNmapFingerprint -> 规则文件第一行必需是 Exclude")
		return errors.New("[-] ParsingNmapFingerprint -> 规则文件第一行必需是 Exclude")
	}

	content := strings.Join(lines, "\n")
	content = "\n" + content

	// 按 "\nProbe" 拆分探针组内容
	probeBlock := strings.Split(content, "\nProbe")
	probeBlock = probeBlock[1:]

	// 进行序列化
	probes, err := serialization.Array2Probes(probeBlock)
	if err != nil {
		return err
	}

	p.Probes = probes

	return nil
}
