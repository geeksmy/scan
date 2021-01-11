package hbhsid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestID_Matches(t *testing.T) {
	id1 := New(1)
	id2 := New(1)
	id3 := New(3)

	assert.True(t, id1.Matches(id2))
	assert.False(t, id1.Matches(id3))
}

func TestEqMatcher_Matches(t *testing.T) {
	id1 := New(1)
	id2 := New(1)
	id3 := New(2)

	if !EQMatcher(id1).Matches(id2) {
		t.Fatalf("expect equal")
	}

	if !EQMatcher(id1).Matches(&id2) {
		t.Fatalf("expect equal")
	}

	if EQMatcher(id1).Matches(id3) {
		t.Fatalf("expect not equal")
	}
}
