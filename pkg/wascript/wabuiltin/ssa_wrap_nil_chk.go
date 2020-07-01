// 版权 @2019 凹语言 作者。保留所有权利。

package wabuiltin

import (
	"fmt"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

func SSA_CheckNil(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	recv := args[0]
	if recv.(*watypes.Value) == nil {
		recvType := args[1]
		methodName := args[2]
		panic(fmt.Sprintf("value method (%s).%s called using nil *%s pointer",
			recvType, methodName, recvType))
	}
	return recv
}
