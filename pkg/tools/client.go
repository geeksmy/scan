package tools

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"

	browser "github.com/EDDYCJY/fake-useragent"
)

func NewHTTPClient(url string, insecureSkipVerify bool, timeout int) ([]byte, int, string, http.Header, error) {
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecureSkipVerify},
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, "", nil, err
	}

	userAgent := browser.Random()
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		// fmt.Println("失败: ", url, err.Error())
		return nil, 0, "", nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, "", nil, err
	}

	server := resp.Header.Get("Server")

	return body, resp.StatusCode, server, resp.Header, nil
}
