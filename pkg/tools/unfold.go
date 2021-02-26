package tools

import (
	"net"
	"strconv"
	"strings"

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
