// 版权 @2019 凹语言 作者。保留所有权利。

package wabuiltin

import (
	"fmt"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

func Len(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	switch x := args[0].(type) {
	case string:
		return len(x)
	case watypes.Array:
		return len(x)
	case *watypes.Value:
		return len((*x).(watypes.Array))
	case []watypes.Value:
		return len(x)
	case map[watypes.Value]watypes.Value:
		return len(x)
	case *watypes.Hashmap:
		return x.Len()
	default:
		panic(fmt.Sprintf("len: illegal operand: %T", x))
	}
}
