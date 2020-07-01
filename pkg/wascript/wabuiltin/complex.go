// 版权 @2019 凹语言 作者。保留所有权利。

package wabuiltin

import (
	"fmt"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

func Complex(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	switch f := args[0].(type) {
	case float32:
		return complex(f, args[1].(float32))
	case float64:
		return complex(f, args[1].(float64))
	default:
		panic(fmt.Sprintf("complex: illegal operand: %T", f))
	}
}

func Real(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	switch c := args[0].(type) {
	case complex64:
		return real(c)
	case complex128:
		return real(c)
	default:
		panic(fmt.Sprintf("real: illegal operand: %T", c))
	}
}

func Imag(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	switch c := args[0].(type) {
	case complex64:
		return imag(c)
	case complex128:
		return imag(c)
	default:
		panic(fmt.Sprintf("imag: illegal operand: %T", c))
	}
}

//complex
