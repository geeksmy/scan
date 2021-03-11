package header

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"syscall"
	"time"

	"scan/pkg/util"
)

type IPV4Service interface {
	SetVersion(version int)
	SetHdrLen(hdrLen int)
	SetDsField(dsField int)
	SetLen(len int)
	SetId(id int)
	SetFlags(flags int)
	SetFragOffset(fragOffset int)
	SetTTL(ttl int)
	SetProto(proto int)
	SetChecksum(checksum int)
	SetSrc(src string)
	SetDst(dst string)
	Marshal() ([]byte, error)
	Show()
	Hex(ipv4Header []byte)
	GetVersion() int
	GetHdrLen() int
	GetDsField() int
	GetLen() int
	GetId() int
	GetFlags() int
	GetFragOffset() int
	GetTTL() int
	GetProto() int
	GetChecksum() int
	GetSrcAddr() net.IP
	GetDstAddr() net.IP
}

type IPV4Header struct {
	Version    int    // 版本 4bit
	HdrLen     int    // header长度 4bit
	DsField    int    // 区分服务 8bit
	Len        int    // 总长度 16bit
	Id         int    // 标识 16bit
	Flags      int    // 标志
	FragOffset int    // 片偏移
	TTL        int    // 生存时间
	Proto      int    // 协议
	Checksum   int    // 校验和
	Src        net.IP // 源地址
	Dst        net.IP // 目的地址
}

func NewIPV4Header() IPV4Service {
	return &IPV4Header{
		Version:    4,
		HdrLen:     20,
		DsField:    0,
		Len:        20,
		Flags:      0,
		FragOffset: 0,
		TTL:        64,
		Checksum:   0,
		Id:         random(0, 65535),
	}
}

func (h *IPV4Header) SetVersion(version int) {
	h.Version = version
}

func (h *IPV4Header) SetHdrLen(hdrLen int) {
	h.HdrLen = hdrLen
}

func (h *IPV4Header) SetDsField(dsField int) {
	h.DsField = dsField
}

func (h *IPV4Header) SetLen(len int) {
	h.Len = len
}

func (h *IPV4Header) SetId(id int) {
	h.Id = id
}

func (h *IPV4Header) SetFlags(flags int) {
	h.Flags = flags
}

func (h *IPV4Header) SetFragOffset(fragOffset int) {
	h.FragOffset = fragOffset
}

func (h *IPV4Header) SetTTL(ttl int) {
	h.TTL = ttl
}

func (h *IPV4Header) SetProto(proto int) {
	h.Proto = proto
}

func (h *IPV4Header) SetChecksum(checksum int) {
	h.Checksum = checksum
}

func (h *IPV4Header) SetSrc(src string) {
	ip, _ := str2int(src)
	h.Src = ip
}

func (h *IPV4Header) SetDst(dst string) {
	ip, _ := str2int(dst)
	h.Dst = ip
}

func (h *IPV4Header) Marshal() ([]byte, error) {
	if h == nil {
		return nil, errors.New("[-] IPV4Header -> 空header")
	}
	if h.Len < util.HeaderLen {
		return nil, errors.New("[-] IPV4Header -> header最小长度是20")
	}
	b := make([]byte, util.HeaderLen)
	b[0] = byte(h.Version<<4 | (util.HeaderLen >> 2 & 0x0f))
	b[1] = byte(h.DsField)
	flagsAndFragOff := (h.FragOffset & 0x1fff) | h.Flags<<13
	binary.BigEndian.PutUint16(b[2:4], uint16(h.Len))
	binary.BigEndian.PutUint16(b[4:6], uint16(h.Id))
	binary.BigEndian.PutUint16(b[6:8], uint16(flagsAndFragOff))
	b[8] = byte(h.TTL)
	b[9] = byte(h.Proto)
	binary.BigEndian.PutUint16(b[10:12], uint16(h.Checksum))
	if ip := h.Src.To4(); ip != nil {
		copy(b[12:16], ip[:net.IPv4len])
	}
	if ip := h.Dst.To4(); ip != nil {
		copy(b[16:20], ip[:net.IPv4len])
	}
	return b, nil
}

func (h *IPV4Header) Show() {
	fmt.Println("\n==========IPV4 Header==========")
	fmt.Printf(`版本[Version]: %d
Header长度[HdrLen]: %d
区分服务[DsField]: %d
总长度[Len]: %d
标识[Id]: %d
标志[Flags]: %d
片偏移[FragOffset]: %d
生存时间[TTL]: %d
协议[Proto]: %s
校验和[Checksum]: %d
源地址[SrcAddr]: %s
目的地址[DstAddr]: %s`, h.Version, h.HdrLen, h.DsField, h.Len, h.Id, h.Flags, h.FragOffset, h.TTL, getProto2String(h.Proto), h.Checksum, h.Src.String(), h.Dst.String())
}

func (h *IPV4Header) Hex(ipv4Header []byte) {
	fmt.Println("\n==========IPV4 Header Hex==========")
	fmt.Printf("%s\n", hex.Dump(ipv4Header))
}

func (h *IPV4Header) GetVersion() int {
	return h.Version
}

func (h *IPV4Header) GetHdrLen() int {
	return h.HdrLen
}

func (h *IPV4Header) GetDsField() int {
	return h.DsField
}

func (h *IPV4Header) GetLen() int {
	return h.Len
}

func (h *IPV4Header) GetId() int {
	return h.Id
}

func (h *IPV4Header) GetFlags() int {
	return h.Flags
}

func (h *IPV4Header) GetFragOffset() int {
	return h.FragOffset
}

func (h *IPV4Header) GetTTL() int {
	return h.TTL
}

func (h *IPV4Header) GetProto() int {
	return h.Proto
}

func (h *IPV4Header) GetChecksum() int {
	return h.Checksum
}

func (h *IPV4Header) GetSrcAddr() net.IP {
	return h.Src
}

func (h *IPV4Header) GetDstAddr() net.IP {
	return h.Dst
}

func getProto2String(proto int) string {
	switch proto {
	case syscall.IPPROTO_TCP:
		return "TCP"
	case syscall.IPPROTO_UDP:
		return "UDP"
	case syscall.IPPROTO_ICMP:
		return "ICMP"
	case syscall.IPPROTO_IP:
		return "IP"
	default:
		return ""
	}
}

func str2int(s string) ([]byte, error) {
	ipS := strings.Split(s, ".")
	ip := make([]byte, len(ipS))
	a, err := strconv.Atoi(ipS[0])
	if err != nil {
		return nil, err
	}
	ip[0] = byte(a)
	b, err := strconv.Atoi(ipS[1])
	if err != nil {
		return nil, err
	}
	ip[1] = byte(b)
	c, err := strconv.Atoi(ipS[2])
	if err != nil {
		return nil, err
	}
	ip[2] = byte(c)
	d, err := strconv.Atoi(ipS[3])
	if err != nil {
		return nil, err
	}
	ip[3] = byte(d)
	return ip, nil
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
