// 版权 @2019 凹语言 作者。保留所有权利。

package waops

import (
	"fmt"

	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

func AsInt(x watypes.Value) int {
	return asInt(x)
}

func AsUint64(x watypes.Value) uint64 {
	return asUint64(x)
}

func asInt(x watypes.Value) int {
	switch x := x.(type) {
	case int:
		return x
	case int8:
		return int(x)
	case int16:
		return int(x)
	case int32:
		return int(x)
	case int64:
		return int(x)
	case uint:
		return int(x)
	case uint8:
		return int(x)
	case uint16:
		return int(x)
	case uint32:
		return int(x)
	case uint64:
		return int(x)
	case uintptr:
		return int(x)
	}
	panic(fmt.Sprintf("cannot convert %T to int", x))
}

func asUint64(x watypes.Value) uint64 {
	switch x := x.(type) {
	case uint:
		return uint64(x)
	case uint8:
		return uint64(x)
	case uint16:
		return uint64(x)
	case uint32:
		return uint64(x)
	case uint64:
		return x
	case uintptr:
		return uint64(x)
	}
	panic(fmt.Sprintf("cannot convert %T to uint64", x))
}
