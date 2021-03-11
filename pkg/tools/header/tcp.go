package header

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"scan/pkg/util"
)

type TCPService interface {
	SetSrcPort(srcPort int)
	SetDstPort(dstPort int)
	SetSeq(seq int)
	SetAck(ack int)
	SetHdrLen(hdrLen int)
	SetFlags(flags int)
	SetWindowSizeValue(windowSizeValue int)
	SetChecksum(checksum int)
	SetUrgentPointer(urgentPointer int)
	SetOptions(options []byte)
	SetOptionsMSS()
	Marshal() ([]byte, error)
	Show()
	Hex(tcpHeader []byte)
	GetSrcPort() int
	GetDstPort() int
	GetSeq() int
	GetAck() int
	GetHdrLen() int
	GetFlags() int
	GetWindowSizeValue() int
	GetChecksum() int
	GetUrgentPointer() int
}

type TCPHeader struct {
	SrcPort         int    // 源端口 随机数 1023-65535
	DstPort         int    // 目的端口
	Seq             int    // 序号 随机数 0-65535*65535
	Ack             int    // 确认号 随机数 0-65535*65535
	HdrLen          int    // Header长度
	Flags           int    // 标志位
	WindowSizeValue int    // 窗口
	Checksum        int    // 校验和
	UrgentPointer   int    // 紧急指针
	Options         []byte // 选项
}

func NewTCPHeader() TCPService {
	return &TCPHeader{
		SrcPort:         getSrcPort(),
		Seq:             getSeq(),
		Ack:             0,
		HdrLen:          util.HeaderLen,
		Flags:           util.SockSyn,
		WindowSizeValue: 1024,
		Checksum:        0,
		UrgentPointer:   0,
	}
}

func (h *TCPHeader) SetSrcPort(srcPort int) {
	h.SrcPort = srcPort
}

func (h *TCPHeader) SetDstPort(dstPort int) {
	h.DstPort = dstPort
}

func (h *TCPHeader) SetSeq(seq int) {
	h.Seq = seq
}

func (h *TCPHeader) SetAck(ack int) {
	h.Ack = ack
}

func (h *TCPHeader) SetHdrLen(hdrLen int) {
	h.HdrLen = hdrLen
}

func (h *TCPHeader) SetFlags(flags int) {
	h.Flags = flags
}

func (h *TCPHeader) SetWindowSizeValue(windowSizeValue int) {
	h.WindowSizeValue = windowSizeValue
}

func (h *TCPHeader) SetChecksum(checksum int) {
	h.Checksum = checksum
}

func (h *TCPHeader) SetUrgentPointer(urgentPointer int) {
	h.UrgentPointer = urgentPointer
}

func (h *TCPHeader) SetOptions(options []byte) {
	h.Options = append(h.Options, options...)
}

func (h *TCPHeader) SetOptionsMSS() {
	mss := make([]byte, 4)
	optionKind := 2
	optionLen := 4
	mssVal := 1460
	binary.BigEndian.PutUint16(mss[:2], uint16(optionKind<<8|optionLen))
	binary.BigEndian.PutUint16(mss[2:], uint16(mssVal))
	h.Options = append(h.Options, mss...)
}

func (h *TCPHeader) Marshal() ([]byte, error) {
	if h == nil {
		return nil, errors.New("[-] TCPHeader -> 空header")
	}

	hdrLen := util.HeaderLen + len(h.Options)
	fmt.Printf("%d", hdrLen)
	b := make([]byte, hdrLen)

	binary.BigEndian.PutUint16(b[0:2], uint16(h.SrcPort))
	binary.BigEndian.PutUint16(b[2:4], uint16(h.DstPort))
	binary.BigEndian.PutUint32(b[4:8], uint32(h.Seq))
	binary.BigEndian.PutUint32(b[8:12], uint32(h.Ack))
	binary.BigEndian.PutUint16(b[12:14], uint16(hdrLen<<10|h.Flags))
	binary.BigEndian.PutUint16(b[14:16], uint16(h.WindowSizeValue))
	binary.BigEndian.PutUint16(b[16:18], uint16(h.Checksum))
	binary.BigEndian.PutUint16(b[18:20], uint16(h.UrgentPointer))
	if len(h.Options) > 0 {
		copy(b[util.HeaderLen:], h.Options)
	}

	return b, nil
}

func (h *TCPHeader) Show() {
	fmt.Println("\n==========TCP Header==========")
	fmt.Printf(`源端口[SrcPort]: %d
目的端口[DstPort]: %d
序号[Seq]: %d
确认号[Ack]: %d
Header长度[HdrLen]: %d
标志位[Flags]: %s
窗口[Win]: %d
校验和[Checksum]: %d
紧急指针[UrgentPointer]: %d
选项[Options]: %b`, h.SrcPort, h.DstPort, h.Seq, h.Ack, h.HdrLen, getFlags2String(h.Flags), h.WindowSizeValue, h.Checksum, h.UrgentPointer, h.Options)
}

func (h *TCPHeader) Hex(tcpHeader []byte) {
	fmt.Println("\n==========TCP Header Hex==========")
	fmt.Printf("%s\n", hex.Dump(tcpHeader))
}

func (h *TCPHeader) GetSrcPort() int {
	return h.SrcPort
}

func (h *TCPHeader) GetDstPort() int {
	return h.DstPort
}

func (h *TCPHeader) GetSeq() int {
	return h.Seq
}

func (h *TCPHeader) GetAck() int {
	return h.Ack
}

func (h *TCPHeader) GetHdrLen() int {
	return h.HdrLen
}

func (h *TCPHeader) GetFlags() int {
	return h.Flags
}

func (h *TCPHeader) GetWindowSizeValue() int {
	return h.WindowSizeValue
}

func (h *TCPHeader) GetChecksum() int {
	return h.Checksum
}

func (h *TCPHeader) GetUrgentPointer() int {
	return h.UrgentPointer
}

func getFlags2String(flags int) string {
	switch flags {
	case util.SockSyn:
		return util.SockSynS
	case util.SockFin:
		return util.SockFinS
	case util.SockRst:
		return util.SockRstS
	case util.SockPush:
		return util.SockPushS
	case util.SockAck:
		return util.SockAckS
	case util.SockUrg:
		return util.SockUrgS
	case util.SockEcn:
		return util.SockEcnS
	case util.SockCwr:
		return util.SockCwrS
	default:
		return ""
	}
}

func getSrcPort() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1<<16-1)%16383 + 49152
}

func getSeq() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(1<<32 - 1)
}
