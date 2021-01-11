// gomock 实现 gomock 的 Matcher interface
package hbhsid

import (
	"fmt"
)

func (id ID) Matches(x interface{}) bool {
	if val, ok := x.(ID); ok {
		return id.Equal(val)
	}

	// if val, ok := x.(*ID); ok {
	// 	return id.Equal(*val)
	// }

	return false
}

func EQMatcher(x ID) eqMatcher {
	return eqMatcher{x: x}
}

type eqMatcher struct {
	x ID
}

func (e eqMatcher) Matches(x interface{}) bool {
	if val, ok := x.(ID); ok {
		return e.x.Equal(val)
	}

	if val, ok := x.(*ID); ok {
		return e.x.Equal(*val)
	}

	return false
}

func (e eqMatcher) String() string {
	return fmt.Sprintf("is equal to %q", e.x)
}
