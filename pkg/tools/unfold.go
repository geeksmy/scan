package tools

import (
	"io/ioutil"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/3th1nk/cidr"
)

// IP地址自增
func IncrIP(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}

func UnfoldIPs(targetIPs []string) (*[]string, error) {
	var res []string
	for i := 0; i < len(targetIPs); i++ {
		switch strings.Contains(targetIPs[i], "/") {
		case true:
			c, err := cidr.ParseCIDR(targetIPs[i])
			if err != nil {
				return nil, err
			}
			a, _, _ := net.ParseCIDR(targetIPs[i])
			for c.Contains(a.String()) {
				res = append(res, a.String())
				IncrIP(a)
			}

		case false:
			res = append(res, targetIPs[i])
		}

	}
	return &res, nil
}

func UnfoldPort(targetPorts []string) (*[]string, error) {
	var res []string

	for i := 0; i < len(targetPorts); i++ {
		switch strings.Contains(targetPorts[i], "-") {
		case true:
			s := strings.Split(targetPorts[i], "-")
			start, err := strconv.Atoi(s[0])
			if err != nil {
				return nil, err
			}
			end, err := strconv.Atoi(s[1])
			if err != nil {
				return nil, err
			}
			for j := start; j < (end + 1); j++ {
				res = append(res, strconv.Itoa(j))
			}
		case false:
			res = append(res, targetPorts[i])
		}
	}
	return &res, nil
}

func Shuffle(slice []string) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

// 判断是否包含端口
func IncludePort(targetPort string, ports string) bool {
	nPorts := strings.Split(ports, ",")
	nmapPorts, err := UnfoldPort(nPorts)
	if err != nil {
		return false
	}
	for _, port := range *nmapPorts {
		if targetPort == port {
			return true
		}
	}
	return false
}

func String2strings(s string) []string {
	s = s[1 : len(s)-1]
	return strings.Split(s, ",")
}

func GetFile2Strings(targetFile string) ([]string, error) {
	fileData, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return nil, err
	}
	ips := strings.Split(string(fileData), "\n")

	return ips, nil
}
