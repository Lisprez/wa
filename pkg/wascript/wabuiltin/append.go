// 版权 @2019 凹语言 作者。保留所有权利。

package wabuiltin

import (
	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

func Append(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	if len(args) == 1 {
		return args[0]
	}

	if s, ok := args[1].(string); ok {
		// append([]byte, ...string) []byte
		arg0 := args[0].([]watypes.Value)
		for i := 0; i < len(s); i++ {
			arg0 = append(arg0, s[i])
		}
		return arg0
	}

	// append([]T, ...[]T) []T
	return append(args[0].([]watypes.Value), args[1].([]watypes.Value)...)
}
