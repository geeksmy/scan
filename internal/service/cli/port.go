package cli

type ScanPort interface {
	SetData(data []byte)
}

type ScanPortService interface {
	InitPackageArgs() ScanPort
}

// 原始套接字参数
type ScanPortArgs struct {
	Data []byte
}

func (s *ScanPortArgs) SetData(data []byte) {
	s.Data = data
}

/**
 * 1. 发送数据包
 *  1.1 整理需要发送数据包的参数
 *  1.2 发送数据包
 * 2. 指纹识别
 * 3. 输出打印
 */

// // 整理数据包参数
// func (p *Port) initPackageArgs() {
//
// }
//
// // 数据包发送
// func (p *Port) packageSend() {
//
// }
//
// // 指纹识别
// func (p *Port) fingerprintRecognition() {
//
// }
//
// // 输出打印
// func (p *Port) outputPrinting() {
//
// }
