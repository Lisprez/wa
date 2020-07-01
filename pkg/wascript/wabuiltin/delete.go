// 版权 @2019 凹语言 作者。保留所有权利。

package wabuiltin

import (
	"fmt"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

func Delete(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	switch m := args[0].(type) {
	case map[watypes.Value]watypes.Value:
		delete(m, args[1])
	case *watypes.Hashmap:
		m.Delete(args[1].(watypes.Hashable))
	default:
		panic(fmt.Sprintf("illegal map type: %T", m))
	}
	return nil
}
