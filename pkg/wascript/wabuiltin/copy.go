// 版权 @2019 凹语言 作者。保留所有权利。

package wabuiltin

import (
	"go/types"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/waops"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

// copy([]T, []T) int or copy([]byte, string) int
func Copy(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	src := args[1]
	if _, ok := src.(string); ok {
		params := fn.Type().(*types.Signature).Params()
		src = waops.Conv(params.At(0).Type(), params.At(1).Type(), src)
	}
	return copy(args[0].([]watypes.Value), src.([]watypes.Value))
}
