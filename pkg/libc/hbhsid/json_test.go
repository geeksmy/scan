package hbhsid

import (
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestID_MarshalJSON(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Uint32()

	id := New(n)
	s := id.String()
	t.Logf("id: %s", s)

	expectOutput, _ := json.Marshal(s)

	marshalOutput, err := json.Marshal(id)
	if err != nil {
		t.Fatalf("json marshal error: %s", err)
	}

	assert.Equal(t, expectOutput, marshalOutput)

	id2 := ID{}
	if err := json.Unmarshal(marshalOutput, &id2); err != nil {
		t.Fatalf("json unmarshal error: %s", err)
	}

	assert.Equal(t, id.orig, id2.orig)
}
