package blasting

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func NewConnTomcat(addr, user, pass, tomcatPath string, timeout int) bool {
	url := addr + tomcatPath
	payload := strings.NewReader("")
	userPass := fmt.Sprintf("%s:%s", user, pass)
	val := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(userPass)))
	client := http.Client{Timeout: time.Duration(timeout)}
	req, err := http.NewRequest("GET", url, payload)
	if err != nil {
		return false
	}
	req.Header.Add("Authorization", val)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return false
	}

	defer res.Body.Close()
	switch res.StatusCode {
	case 200:
		return true
	case 301:
		return true
	case 302:
		return true
	default:
		return false
	}
}
