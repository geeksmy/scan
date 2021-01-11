// hbhsid `hide base32 hex string id`
// 工作原理: 使用 base32 将 uin32 转换成字符串输出, 使用 hide 来隐藏原始的 id
// 该 id 实现为只能支持 uint32
// 考虑目前应用规模 2^32 的空间足够使用
// example:
//  func init() {
//	    _ = Hide.SetUint32(big.NewInt(1500450271))
//	    _ = Hide.SetXor(big.NewInt(1500450271))
// }
// id := New(14147519656024107973)
// id.String() // FUCI3QIY7LWEI
//
package hbhsid

import (
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/c2h5oh/hide"
)

type WrongByteErr int

func (e WrongByteErr) Error() string {
	return fmt.Sprintf("read wrong bytes %d", e)
}

var (
	b32encoding = base32.StdEncoding.WithPadding(base32.NoPadding)
	varintLen   = binary.MaxVarintLen32

	// caller 需要初始化 hide prime
	Hide hide.Hide
	Zero ID
)

func New(i uint32) ID {
	h := Hide.Uint32Obfuscate(i)

	bs := make([]byte, varintLen)
	binary.PutUvarint(bs, uint64(h))
	return ID{src: bs, orig: i, hide: h}
}

func ParseFromString(s string) (ID, error) {
	s = strings.ToUpper(s)
	return ParseBytes([]byte(s))
}

func ParseBytes(src []byte) (ID, error) {
	dstBuf := make([]byte, 2*varintLen)
	n, err := b32encoding.Decode(dstBuf, src)
	if err != nil {
		return Zero, fmt.Errorf("decode wrong bytes %w", err)
	}

	dstBuf = dstBuf[:n]
	h, rn := binary.Uvarint(dstBuf)
	// rn == 0: buf to small
	// rn < 0: value overflow,
	// see: go/src/encoding/binary/varint.go:60
	if rn <= 0 {
		return Zero, WrongByteErr(rn)
	}

	o := Hide.Uint32Deobfuscate(uint32(h))

	m := ID{
		src:  dstBuf,
		orig: o,
		hide: uint32(h),
	}

	return m, nil
}

type ID struct {
	src  []byte
	orig uint32
	hide uint32
}

func (id ID) Origin() uint32 {
	return id.orig
}
func (id ID) String() string {
	return b32encoding.EncodeToString(id.src)
}

func (id *ID) FromString(s string) error {
	_id, err := ParseFromString(s)
	if err != nil {
		return err
	}
	*id = _id
	return err
}

func (id *ID) Equal(id2 ID) bool {
	return id.Compare(id2) == 0
}

// Compare 比较 ID 大小, 使用 origin 进行比较
// @return -1 - id < id2; 0: id == id2; 1 - id > id2
func (id ID) Compare(id2 ID) int8 {
	if id.orig == id2.orig {
		return 0
	}

	if id.orig < id2.orig {
		return -1
	}

	return 1
}
