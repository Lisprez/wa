// 版权 @2019 凹语言 作者。保留所有权利。

package wabuiltin

import (
	"fmt"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

func Cap(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	switch x := args[0].(type) {
	case watypes.Array:
		return cap(x)
	case *watypes.Value:
		return cap((*x).(watypes.Array))
	case []watypes.Value:
		return cap(x)
	default:
		panic(fmt.Sprintf("cap: illegal operand: %T", x))
	}
}
