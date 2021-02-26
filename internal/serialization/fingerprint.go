package serialization

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"scan/internal/model"
	"scan/pkg/tools"
)

func Array2Probes(array []string) (*[]model.Probe, error) {
	res := make([]model.Probe, len(array))

	for i := 0; i < len(array); i++ {
		probe, err := String2Probe(array[i])
		if err != nil {
			return nil, err
		}
		res[i] = *probe
	}

	return &res, nil
}

func String2Probe(s string) (*model.Probe, error) {
	var res model.Probe

	s = strings.TrimSpace(s)
	lines := strings.Split(s, "\n")

	// 提取Name Data Protocol
	header := lines[0]
	protocol := header[:4]
	other := header[4:]
	if !(protocol == "TCP " || protocol == "UDP ") {
		return nil, errors.New("[-] serialization -> Probe Protocol必须是TCP或者UDP ")
	}
	if len(other) == 0 {
		return nil, errors.New("[-] serialization -> Probe Name 损坏 ")
	}

	d := tools.NewDirective()
	directive, err := d.GetDirective(other)
	if err != nil {
		return nil, err
	}

	res.Name = directive.DirectiveName
	res.Data = strings.Split(directive.DirectiveStr, directive.Delimiter)[0]
	res.Protocol = strings.ToLower(strings.TrimSpace(protocol))

	for i := 0; i < len(lines); i++ {
		if strings.HasPrefix(lines[i], "match ") {
			match, err := String2Match(lines[i], "match")
			if err != nil {
				if err.Error() == "解析正则表达式失败" {
					continue
				}
				return nil, err
			}
			res.Matchs = append(res.Matchs, match)
		}
		if strings.HasPrefix(lines[i], "softmatch ") {
			softMatch, err := String2Match(lines[i], "softmatch")
			if err != nil {
				if err.Error() == "解析正则表达式失败" {
					continue
				}
				return nil, err
			}
			res.Matchs = append(res.Matchs, softMatch)
		}
		if strings.HasPrefix(lines[i], "ports ") {
			res.Ports = lines[i][len("ports")+1:]
		}
		if strings.HasPrefix(lines[i], "sslports ") {
			res.SSLPorts = lines[i][len("sslports")+1:]
		}
		if strings.HasPrefix(lines[i], "totalwaitms ") {
			res.TotalWaitMS, err = strconv.Atoi(lines[i][len("totalwaitms")+1:])
			if err != nil {
				return nil, errors.New("[-] serialization -> totalwaitms 转换错误 ")
			}
		}
		if strings.HasPrefix(lines[i], "tcpwrappedms ") {
			res.TCPWrappedMS, err = strconv.Atoi(lines[i][len("tcpwrappedms")+1:])
			if err != nil {
				return nil, errors.New("[-] serialization -> tcpwrappedms 转换错误 ")
			}
		}
		if strings.HasPrefix(lines[i], "rarity ") {
			res.Rarity, err = strconv.Atoi(lines[i][len("rarity")+1:])
			if err != nil {
				return nil, errors.New("[-] serialization -> rarity 转换错误 ")
			}
		}
		if strings.HasPrefix(lines[i], "fallback ") {
			res.Fallback = lines[i][len("fallback")+1:]
		}
	}

	return &res, nil
}

func String2Match(s, name string) (*model.Match, error) {
	var res model.Match

	if name == "softmatch" {
		res.IsSoft = true
	}

	matchText := s[len(name)+1:]
	d := tools.NewDirective()
	directive, err := d.GetDirective(matchText)
	if err != nil {
		return nil, err
	}

	textSplinted := strings.Split(directive.DirectiveStr, directive.Delimiter)
	pattern, versionInfo := textSplinted[0], strings.Join(textSplinted[1:], "")

	patternUnescaped, _ := tools.DecodePattern(pattern)
	patternUnescapedStr := string([]rune(string(patternUnescaped)))
	// TODO: 部分规则 会报错, 无法获取Regexp, 需要后期优化
	patternCompiled, errReg := regexp.Compile(patternUnescapedStr)
	if errReg != nil {
		return nil, errors.New("解析正则表达式失败")
	}

	res.Service = directive.DirectiveName
	res.Pattern = pattern
	res.VersionInfo = versionInfo
	res.PatternCompiled = patternCompiled

	return &res, nil
}
