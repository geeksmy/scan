package util

const (
	// 常量
	HeaderLen    = 20
	Version      = 4
	MaxHeaderLen = 60

	// FLAGS
	SockFin  = 0x01
	SockSyn  = 0x02
	SockRst  = 0x04
	SockPush = 0x08
	SockAck  = 0x10
	SockUrg  = 0x20
	SockEcn  = 0x40
	SockCwr  = 0x80

	// FLAGS string
	SockFinS  = "FIN"
	SockSynS  = "SYN"
	SockRstS  = "RST"
	SockPushS = "PUSH"
	SockAckS  = "ACK"
	SockUrgS  = "URG"
	SockEcnS  = "ECN"
	SockCwrS  = "CWR"
)
