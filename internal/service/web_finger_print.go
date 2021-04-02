package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"scan/config"
	"scan/pkg/tools"

	"github.com/axgle/mahonia"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type WebFingerPrintSVC interface {
	/**
	 * InitArgs 初始化参数
	 * @param cmd 命令行传入参数结构体
	 */
	InitCmdArgs(cmd *cobra.Command) (*WebFingerPrintCmdArgs, error)
	/**
	 * InitFingerPrintFile 初始化指纹文件
	 */
	InitFingerPrintFile() error
	/**
	 * InitPackageArgs 初始化发送参数
	 */
	InitRequestArgs()
	/**
	 * SendRequest 发送请求
	 * @param responses 返回参数管道
	 * @param mainWG 计数器
	 */
	SendRequest(responses chan WebFingerPrintResponse, mainWG *sync.WaitGroup)
	/**
	 * IdentifyResponse 识别指纹
	 * @param responses 返回参数管道
	 * @param results 指纹识别后返回值管道
	 * @param mainWG 计数器
	 */
	IdentifyResponse(responses chan WebFingerPrintResponse, results chan WebFingerPrintResult, mainWG *sync.WaitGroup)
	/**
	 * OutputPrinting 输出打印
	 * @param results 指纹识别后返回值管道
	 */
	OutputPrinting(results chan WebFingerPrintResult)
}

type WebFingerPrintResult struct {
	Url         string
	StateCode   int
	Server      string
	Title       string
	FingerPrint string
	Retry       int
}

type WebFingerPrintResponse struct {
	Url       string
	StateCode int
	Body      []byte
	Hash      string
	Server    string
	Title     string
	Header    http.Header
}

type JsonFile struct {
	Fingerprint []WebFingerPrintArgs `json:"fingerprint"`
}

type WebFingerPrintArgs struct {
	Name     string   `json:"cms"`
	Method   string   `json:"method"`
	Location string   `json:"location"`
	Keyword  []string `json:"keyword"`
}

type WebFingerPrintCmdArgs struct {
	TargetIPs   *[]string
	TargetPorts *[]string
	FileName    string
	Thread      int
	Timeout     int
	Retry       int
	OutPut      string
	OutFileName string
}

type WebFingerPrint struct {
	// db     *gorm.DB
	logger *zap.Logger

	Args            WebFingerPrintCmdArgs
	KeyWordArgs     []WebFingerPrintArgs
	FaviconHashArgs []WebFingerPrintArgs
	Urls            []string
}

func NewWebFingerPrintSVC(logger *zap.Logger) WebFingerPrintSVC {
	return &WebFingerPrint{
		logger: logger,
	}
}

func (svc *WebFingerPrint) InitCmdArgs(cmd *cobra.Command) (*WebFingerPrintCmdArgs, error) {
	conf := config.C

	targetFileName, _ := cmd.Flags().GetString("target-urls")
	switch targetFileName {
	case "":
		targetUrls, err := tools.GetFile2Strings(conf.WebFingerprint.TargetUrls)
		if err != nil {
			return nil, err
		}
		ips, err := tools.UnfoldIPs(targetUrls)
		if err != nil {
			return nil, err
		}
		tools.Shuffle(*ips)
		svc.Args.TargetIPs = ips
	default:
		targetUrls, err := tools.GetFile2Strings(targetFileName)
		if err != nil {
			return nil, err
		}
		ips, err := tools.UnfoldIPs(targetUrls)
		if err != nil {
			return nil, err
		}
		tools.Shuffle(*ips)
		svc.Args.TargetIPs = ips
	}

	targetPorts, _ := cmd.Flags().GetStringArray("target-ports")
	switch len(targetPorts) {
	case 1:
		ports, err := tools.UnfoldPort(tools.String2strings(targetPorts[0]))
		if err != nil {
			return nil, err
		}
		svc.Args.TargetPorts = ports
	case 0:
		// 如果不传命令行参数使用配置文件的配置
		ports, err := tools.UnfoldPort(conf.WebFingerprint.TargetPorts)
		if err != nil {
			return nil, err
		}
		svc.Args.TargetPorts = ports
	default:
		ports, err := tools.UnfoldPort(targetPorts)
		if err != nil {
			return nil, err
		}
		svc.Args.TargetPorts = ports
	}

	thread, _ := cmd.Flags().GetInt("thread")
	switch thread {
	case 0:
		svc.Args.Thread = conf.WebFingerprint.Thread
	default:
		svc.Args.Thread = thread
	}

	if svc.Args.Thread <= 1 {
		svc.Args.Thread = 1
	}

	retry, _ := cmd.Flags().GetInt("retry")
	switch retry {
	case 0:
		svc.Args.Retry = conf.WebFingerprint.Retry
	default:
		svc.Args.Retry = retry
	}

	if svc.Args.Retry <= 1 {
		svc.Args.Retry = 1
	}

	timeout, _ := cmd.Flags().GetInt("timeout")
	switch timeout {
	case 0:
		svc.Args.Timeout = conf.WebFingerprint.Timeout
	default:
		svc.Args.Timeout = timeout
	}

	if svc.Args.Timeout <= 1 {
		svc.Args.Timeout = 1
	}

	fileName, _ := cmd.Flags().GetString("fingerprint-file")
	switch fileName {
	case "":
		svc.Args.FileName = conf.WebFingerprint.FingerprintName
	default:
		svc.Args.FileName = fileName
	}

	outFile, _ := cmd.Flags().GetString("out-file")
	switch outFile {
	case "":
		svc.Args.OutPut = "print"
	default:
		svc.Args.OutPut = "file"
		svc.Args.OutFileName = outFile
	}

	return &svc.Args, nil
}

func (svc *WebFingerPrint) InitFingerPrintFile() error {
	fileData, err := ioutil.ReadFile(svc.Args.FileName)
	if err != nil {
		return err
	}

	var jsonFile JsonFile
	err = json.Unmarshal(fileData, &jsonFile)
	if err != nil {
		return err
	}
	for _, args := range jsonFile.Fingerprint {
		switch args.Method {
		case "keyword":
			svc.KeyWordArgs = append(svc.KeyWordArgs, args)
		case "faviconhash":
			svc.FaviconHashArgs = append(svc.FaviconHashArgs, args)
		default:
			continue
		}
	}
	return nil
}

func (svc *WebFingerPrint) InitRequestArgs() {
	for _, port := range *svc.Args.TargetPorts {
		for _, ip := range *svc.Args.TargetIPs {
			switch port {
			case "443", "3443", "4443", "5443", "6443", "7443", "8443", "9443", "10443", "4430":
				url := fmt.Sprintf("https://%s:%s", ip, port)
				svc.Urls = append(svc.Urls, url)
			default:
				url := fmt.Sprintf("http://%s:%s", ip, port)
				svc.Urls = append(svc.Urls, url)
			}
		}
	}
}

/**
 * sendRequestWork 发送请求线程池
 * @param responses 返回参数管道
 * @param wg 计数器
 */
func sendRequestWork(urls <-chan string, responses chan<- WebFingerPrintResponse, wg *sync.WaitGroup, timeout int) {
	for url := range urls {
		var insecureSkipVerify bool
		if strings.Contains(url, "https") {
			insecureSkipVerify = true
		}

		body, stateCode, server, header, err := tools.NewHTTPClient(url, insecureSkipVerify, timeout)
		if err != nil {
			wg.Done()
			continue
		}

		response := WebFingerPrintResponse{
			Url:       url,
			StateCode: stateCode,
			Body:      body,
			Server:    server,
			Header:    header,
		}

		re, errRe := regexp.Compile(`<title>(.*?)</title>`)
		if errRe != nil {
			response.Title = ""
		} else {
			title := re.FindAllStringSubmatch(string(body), -1)
			if title != nil {
				if title[0] != nil {
					response.Title = title[0][1]
					if len(title[0][1]) > 45 {
						s := title[0][1][:40] + "..."
						response.Title = s
					}
				}
			}
		}

		var (
			charset string
			codeRe  *regexp.Regexp
		)
		codeRe, errRe = regexp.Compile(`charset=(.*?)["|']`)
		if strings.Contains(string(body), `charset=["|']`) {
			codeRe, errRe = regexp.Compile(`charset="(.*?)["|']`)
		}
		if errRe != nil {
			charset = "utf-8"
		} else {
			charsetByte := codeRe.FindAllStringSubmatch(string(body), -1)
			if charsetByte != nil {
				if charsetByte[0] != nil {
					charset = charsetByte[0][1]
				}
			}
		}

		// if charset == "" && !strings.Contains(string(body), "charset=") {
		// 	charset = "gbk"
		// }

		switch strings.ToLower(charset) {
		case strings.ToLower("GB2312"), strings.ToLower("Big5"), strings.ToLower("GB18030"), strings.ToLower("GBK"):
			charset = "GBK"
		default:
			charset = "utf-8"
		}

		if strings.ToLower(charset) != "utf-8" {
			decoder := mahonia.NewDecoder(charset)
			response.Title = decoder.ConvertString(response.Title)
		}

		faviconUrl := fmt.Sprintf("%s/favicon.ico", url)
		hashBody, _, _, _, err := tools.NewHTTPClient(faviconUrl, insecureSkipVerify, timeout)
		if err != nil {
			responses <- response
			wg.Done()
			continue
		}

		response.Hash = tools.Mmh3Hash32(tools.StandBase64(hashBody))
		responses <- response

		wg.Done()
	}
}

func (svc *WebFingerPrint) SendRequest(responses chan WebFingerPrintResponse, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(responses)

	var wg sync.WaitGroup
	urls := make(chan string, svc.Args.Thread)

	for i := 0; i < cap(urls); i++ {
		go sendRequestWork(urls, responses, &wg, svc.Args.Timeout)
	}

	for _, url := range svc.Urls {
		wg.Add(1)
		urls <- url
	}

	wg.Wait()
	close(urls)
}

func identifyResponseWork(responses <-chan WebFingerPrintResponse, results chan<- WebFingerPrintResult, wg *sync.WaitGroup,
	keyword, faviconHash []WebFingerPrintArgs, retry int) {
	for response := range responses {
		result := WebFingerPrintResult{
			Url:       response.Url,
			StateCode: response.StateCode,
			Server:    response.Server,
			Title:     response.Title,
			Retry:     0,
		}

		for i := 0; i < retry; i++ {
			if result.Retry < retry && result.FingerPrint == "" {
				for _, args := range keyword {
					switch args.Location {
					case "body":
						for i := 0; i < len(args.Keyword); i++ {
							if strings.Contains(string(response.Body), args.Keyword[i]) {
								// result.FingerPrint = args.Name
								break
							}
						}
					case "header":
						for i := 0; i < len(args.Keyword); i++ {
							for _, v := range response.Header {
								for j := 0; j < len(v); j++ {
									if strings.Contains(v[i], args.Keyword[i]) {
										// result.FingerPrint = args.Name
										break
									}
								}
							}
						}
					default:
						continue
					}

				}

				if response.Hash != "" && result.FingerPrint != "" {
					for _, hash := range faviconHash {
						for _, v := range hash.Keyword {
							if response.Hash == v {
								// result.FingerPrint = hash.Name
								break
							}
						}
					}
				}

				result.Retry += 1
			}
		}

		results <- result
		wg.Done()
	}

}

func (svc *WebFingerPrint) IdentifyResponse(responses chan WebFingerPrintResponse, results chan WebFingerPrintResult, mainWG *sync.WaitGroup) {
	defer mainWG.Done()
	defer close(results)

	var wg sync.WaitGroup
	res := make(chan WebFingerPrintResponse, svc.Args.Thread)

	for i := 0; i < cap(res); i++ {
		go identifyResponseWork(res, results, &wg, svc.KeyWordArgs, svc.FaviconHashArgs, svc.Args.Retry)
	}

	for response := range responses {
		wg.Add(1)
		res <- response
	}

	wg.Wait()
	close(res)
}

func (svc *WebFingerPrint) OutputPrinting(results chan WebFingerPrintResult) {

	_, err := os.Stat(svc.Args.OutFileName)
	if err == nil {
		// 如果文件存在
		_ = os.Remove(svc.Args.OutFileName)
	}

	file, _ := os.Create(svc.Args.OutFileName)

	switch svc.Args.OutPut {
	case "file":
		_, _ = file.WriteString(fmt.Sprintf("%-30s%-30.25s%-15s%-20s%-50s\n", "url", "server", "state_code", "fingerprint", "title"))
		fmt.Printf("%-30s%-30.25s%-15s%-20s%-50s\n", "url", "server", "state_code", "fingerprint", "title")
		for result := range results {
			fmt.Printf("%-30s%-30.25s%-3d%-12s%-20s%-50s\n", result.Url, result.Server, result.StateCode, "", result.FingerPrint, result.Title)
			_, _ = file.WriteString(fmt.Sprintf("%-30s%-30.25s%-3d%-12s%-20s%-50.45s\n", result.Url, result.Server, result.StateCode, "", result.FingerPrint, result.Title))
		}
	default:
		fmt.Printf("%-30s%-30.25s%-15s%-20s%-50s\n", "url", "server", "state_code", "fingerprint", "title")
		for result := range results {
			fmt.Printf("%-30s%-30.25s%-3d%-12s%-20s%-50s\n", result.Url, result.Server, result.StateCode, "", result.FingerPrint, result.Title)
		}
	}
}
