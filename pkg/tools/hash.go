package tools

import (
	"bytes"
	"encoding/base64"
	"fmt"

	"github.com/twmb/murmur3"
)

func Mmh3Hash32(raw []byte) string {
	var h32 = murmur3.New32()
	_, _ = h32.Write(raw)
	return fmt.Sprintf("%d", int32(h32.Sum32()))
}

func StandBase64(body []byte) []byte {
	base64Body := base64.StdEncoding.EncodeToString(body)
	var buffer bytes.Buffer
	for i := 0; i < len(base64Body); i++ {
		ch := base64Body[i]
		buffer.WriteByte(ch)
		if (i+1)%76 == 0 {
			buffer.WriteByte('\n')
		}
	}
	buffer.WriteByte('\n')
	return buffer.Bytes()

}
