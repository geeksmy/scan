package tools

import (
	"fmt"
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
	// s = s[1 : len(s)-1]
	return strings.Split(s, ",")
}

func GetFile2Strings(targetFile string) ([]string, error) {
	fileData, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return nil, err
	}
	strS := strings.Split(string(fileData), "\n")
	if strings.Contains(string(fileData), "\r\n") {
		strS = strings.Split(string(fileData), "\r\n")
	}

	s := removeDuplicatesAndEmpty(strS)

	return s, nil
}

// 数组去重 去空
func removeDuplicatesAndEmpty(a []string) (ret []string) {
	aLen := len(a)
	for i := 0; i < aLen; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

// 生成内网网段
func GenerateIntranet(s []string, d []string) []string {
	var ips []string
	for i := 0; i < len(s); i++ {
		x := strings.Split(s[i], ".")
		switch len(x) {
		case 2:
			for z := 1; z < 256; z++ {
				for c := 1; c < 256; c++ {
					for v := 0; v < len(d); v++ {
						ips = append(ips, fmt.Sprintf("%s.%d.%d.%s", x[0], z, c, d[v]))
					}
				}
			}
		case 3:
			for z := 1; z < 256; z++ {
				for c := 0; c < len(d); c++ {
					ips = append(ips, fmt.Sprintf("%s.%s.%d.%s", x[0], x[1], z, d[c]))
				}
			}
		}
	}

	return ips
}
