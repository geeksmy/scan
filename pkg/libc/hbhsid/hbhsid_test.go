package hbhsid

import (
	"encoding/base32"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	_ = Hide.SetUint32(big.NewInt(1500450271))
	_ = Hide.SetXor(big.NewInt(1414751965602410773))
}

func TestNewUint64(t *testing.T) {
	cases := [][2]interface{}{
		{uint32(103902621), "22F5RBIM"},
		{uint32(1<<32 - 1), "WSDNFLAG"},
	}

	for _, tc := range cases {
		i := tc[0].(uint32)
		expect := tc[1].(string)

		id := New(i)
		assert.Equal(t, i, id.Origin())
		assert.Equal(t, expect, id.String())

		_id := ID{}
		err := _id.FromString(id.String())
		assert.NoError(t, err)
		assert.Equal(t, i, _id.Origin())
	}
}

func TestID_Compare(t *testing.T) {
	cases := [][3]int64{
		{100, 99, 1},
		{99, 199, -1},
		{101, 101, 0},
	}

	for i, tc := range cases {
		id1 := New(uint32(tc[0]))
		id2 := New(uint32(tc[1]))

		cmp := id1.Compare(id2)
		assert.Equal(t, cmp, int8(tc[2]),
			fmt.Sprintf("testcase[%d] %q compare %q expect %d actual %d", i, id1, id2, tc[2], cmp))
	}
}

func TestID_Equal(t *testing.T) {
	cases := [][3]interface{}{
		{uint32(100), uint32(99), false},
		{uint32(99), uint32(199), false},
		{uint32(101), uint32(101), true},
	}

	for i, tc := range cases {
		id1 := New(tc[0].(uint32))
		id2 := New(tc[1].(uint32))

		eq := id1.Equal(id2)
		assert.True(t, eq == tc[2].(bool),
			fmt.Sprintf("testcase[%d]: %q eq %q expect[%v] return %v", i, id1, id2, tc[2].(bool), eq))
	}
}

func TestParseBytes(t *testing.T) {
	id := New(8964)
	bs := []byte(id.String())

	_id, err := ParseBytes(bs)
	assert.NoError(t, err)
	assert.True(t, id.Equal(_id))
}

func TestParseBytes_ErrLowercase(t *testing.T) {
	id := New(2346)
	bs := []byte(strings.ToLower(id.String()))
	t.Log("lowercase bytes: ", bs)

	_id, err := ParseBytes(bs)
	assert.Error(t, err)
	if err != nil {
		assert.IsType(t, base32.CorruptInputError(0), errors.Unwrap(err))
	}
	assert.True(t, _id.Equal(Zero))
}
