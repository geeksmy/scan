package blasting

import (
	"net"
	"strconv"
	"strings"
	"time"
)

type TelnetClient struct {
	addr             string
	IsAuthentication bool
	Username         string
	Password         string
}

const (
	TimeDelayAfterWrite = 300 // 300ms
)

func NewConnTelnet(addr, user, pass string, timeout int) bool {
	client := new(TelnetClient)
	client.addr = addr
	client.Username = user
	client.Password = pass
	client.IsAuthentication = true
	action := []string{"w_ping", "r_4096"}
	e := client.Telnet(action, timeout)
	switch e {
	case false:
		return false
	default:
		return true
	}
}

func (t *TelnetClient) PortIsOpen(timeout int) bool {
	conn, err := net.DialTimeout("tcp", t.addr, time.Duration(timeout)*time.Second)
	if nil != err {
		return false
	}
	defer conn.Close()
	return true
}

func (t *TelnetClient) Telnet(action []string, timeout int) bool {
	conn, err := net.DialTimeout("tcp", t.addr, time.Duration(timeout)*time.Second)
	if nil != err {
		return false
	}
	defer conn.Close()
	if false == t.telnetProtocolHandshake(conn) {
		return false
	}
	//	conn.SetReadDeadline(time.Now().Add(time.Second * 30))
	for _, v := range action {
		actSlice := strings.SplitN(v, "_", 2)
		if 2 > len(actSlice) {
			return false
		}
		switch actSlice[0] {
		case "r":
			var n int
			n, err = strconv.Atoi(actSlice[1])
			if nil != err {
				return false
			}
			p := make([]byte, n)
			//	p := make([]byte, 0, n)
			n, err = conn.Read(p[0:])
			if nil != err {
				return false
			}
			var buf []byte
			buf = append(buf, p[0:n]...)
			// fmt.Println("read data length:", n)
			// fmt.Println(string(p[0:n]) + "\n\n")
			// fmt.Println(string(buf))
			if n > 100 {
				return true
			}
		case "w":
			_, err = conn.Write([]byte(actSlice[1] + "\n"))
			if nil != err {
				return false
			}
			// fmt.Println("wirte:", actSlice[1])
			time.Sleep(time.Millisecond * TimeDelayAfterWrite)
		}
	}
	return false
}

func (t *TelnetClient) telnetProtocolHandshake(conn net.Conn) bool {
	var buf [4096]byte
	n, err := conn.Read(buf[0:])
	if nil != err {
		return false
	}
	// fmt.Println(string(buf[0:n]))
	// fmt.Println((buf[0:n]))

	buf[1] = 252
	buf[4] = 252
	buf[7] = 252
	buf[10] = 252
	// fmt.Println((buf[0:n]))
	n, err = conn.Write(buf[0:n])
	if nil != err {
		return false
	}
	time.Sleep(time.Millisecond * TimeDelayAfterWrite)

	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}
	// fmt.Println(string(buf[0:n]))
	// fmt.Println((buf[0:n]))

	buf[1] = 252
	buf[4] = 251
	buf[7] = 252
	buf[10] = 254
	buf[13] = 252
	// fmt.Println((buf[0:n]))
	n, err = conn.Write(buf[0:n])
	if nil != err {
		return false
	}
	time.Sleep(time.Millisecond * TimeDelayAfterWrite)

	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}
	// fmt.Println(string(buf[0:n]))
	// fmt.Println((buf[0:n]))

	if false == t.IsAuthentication {
		return true
	}

	n, err = conn.Write([]byte(t.Username + "\n"))
	if nil != err {
		return false
	}
	time.Sleep(time.Millisecond * TimeDelayAfterWrite)

	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}
	// fmt.Println(string(buf[0:n]))

	n, err = conn.Write([]byte(t.Password + "\n"))
	if nil != err {
		return false
	}
	time.Sleep(time.Millisecond * TimeDelayAfterWrite)

	n, err = conn.Read(buf[0:])
	if nil != err {
		return false
	}
	// fmt.Println(string(buf[0:n]))
	return true
}
