package blasting

import (
	"time"

	"github.com/jlaffaye/ftp"
)

func NewConnFTP(addr, user, pass string, timeout int) bool {
	client, err := ftp.Dial(addr, ftp.DialWithTimeout(time.Duration(timeout)*time.Second))
	if err != nil {
		return false
	}
	if err = client.Login(user, pass); err != nil {
		return false
	}

	_ = client.Logout()
	return true
}
