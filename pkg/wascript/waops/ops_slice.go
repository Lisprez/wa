// 版权 @2019 凹语言 作者。保留所有权利。

package waops

import (
	"fmt"

	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

// 切片操作
func Slice(x, lo, hi, max watypes.Value) watypes.Value {
	return slice(x, lo, hi, max)
}

func slice(x, lo, hi, max watypes.Value) watypes.Value {
	var Len, Cap int
	switch x := x.(type) {
	case string:
		Len = len(x)
	case []watypes.Value:
		Len = len(x)
		Cap = cap(x)
	case *watypes.Value: // *array
		a := (*x).(watypes.Array)
		Len = len(a)
		Cap = cap(a)
	}

	l := 0
	if lo != nil {
		l = AsInt(lo)
	}

	h := Len
	if hi != nil {
		h = AsInt(hi)
	}

	m := Cap
	if max != nil {
		m = AsInt(max)
	}

	switch x := x.(type) {
	case string:
		return x[l:h]
	case []watypes.Value:
		return x[l:h:m]
	case *watypes.Value: // *array
		a := (*x).(watypes.Array)
		return []watypes.Value(a)[l:h:m]
	}
	panic(fmt.Sprintf("slice: unexpected X type: %T", x))
}
