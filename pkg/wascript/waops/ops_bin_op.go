// 版权 @2019 凹语言 作者。保留所有权利。

package waops

import (
	"fmt"
	"go/token"
	"go/types"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

// 实现全部的二元运算
func BinOp(op token.Token, t types.Type, x, y watypes.Value) watypes.Value {
	return binop(op, t, x, y)
}

// 判断 x 和 y 是否相等, 对于引用类型可以和 nil 比较.
func eqNil(t types.Type, x, y watypes.Value) bool {
	switch t.Underlying().(type) {
	case *types.Map, *types.Signature, *types.Slice:
		switch x := x.(type) {
		case *watypes.Hashmap:
			return (x != nil) == (y.(*watypes.Hashmap) != nil)
		case map[watypes.Value]watypes.Value:
			return (x != nil) == (y.(map[watypes.Value]watypes.Value) != nil)
		case *ssa.Function:
			switch y := y.(type) {
			case *ssa.Function:
				return (x != nil) == (y != nil)
			case *watypes.Closure:
				return true
			}
		case *watypes.Closure:
			return (x != nil) == (y.(*ssa.Function) != nil)
		case []watypes.Value:
			return (x != nil) == (y.([]watypes.Value) != nil)
		}

		panic(fmt.Sprintf("eqnil(%s): illegal dynamic type: %T", t, x))
	}

	return watypes.Equals(t, x, y)
}

func binop(op token.Token, t types.Type, x, y watypes.Value) watypes.Value {
	switch op {
	case token.EQL:
		return eqNil(t, x, y)
	case token.NEQ:
		return !eqNil(t, x, y)

	case token.ADD:
		switch x.(type) {
		case int:
			return x.(int) + y.(int)
		case int8:
			return x.(int8) + y.(int8)
		case int16:
			return x.(int16) + y.(int16)
		case int32:
			return x.(int32) + y.(int32)
		case int64:
			return x.(int64) + y.(int64)
		case uint:
			return x.(uint) + y.(uint)
		case uint8:
			return x.(uint8) + y.(uint8)
		case uint16:
			return x.(uint16) + y.(uint16)
		case uint32:
			return x.(uint32) + y.(uint32)
		case uint64:
			return x.(uint64) + y.(uint64)
		case uintptr:
			return x.(uintptr) + y.(uintptr)
		case float32:
			return x.(float32) + y.(float32)
		case float64:
			return x.(float64) + y.(float64)
		case complex64:
			return x.(complex64) + y.(complex64)
		case complex128:
			return x.(complex128) + y.(complex128)
		case string:
			return x.(string) + y.(string)
		}

	case token.SUB:
		switch x.(type) {
		case int:
			return x.(int) - y.(int)
		case int8:
			return x.(int8) - y.(int8)
		case int16:
			return x.(int16) - y.(int16)
		case int32:
			return x.(int32) - y.(int32)
		case int64:
			return x.(int64) - y.(int64)
		case uint:
			return x.(uint) - y.(uint)
		case uint8:
			return x.(uint8) - y.(uint8)
		case uint16:
			return x.(uint16) - y.(uint16)
		case uint32:
			return x.(uint32) - y.(uint32)
		case uint64:
			return x.(uint64) - y.(uint64)
		case uintptr:
			return x.(uintptr) - y.(uintptr)
		case float32:
			return x.(float32) - y.(float32)
		case float64:
			return x.(float64) - y.(float64)
		case complex64:
			return x.(complex64) - y.(complex64)
		case complex128:
			return x.(complex128) - y.(complex128)
		}

	case token.MUL:
		switch x.(type) {
		case int:
			return x.(int) * y.(int)
		case int8:
			return x.(int8) * y.(int8)
		case int16:
			return x.(int16) * y.(int16)
		case int32:
			return x.(int32) * y.(int32)
		case int64:
			return x.(int64) * y.(int64)
		case uint:
			return x.(uint) * y.(uint)
		case uint8:
			return x.(uint8) * y.(uint8)
		case uint16:
			return x.(uint16) * y.(uint16)
		case uint32:
			return x.(uint32) * y.(uint32)
		case uint64:
			return x.(uint64) * y.(uint64)
		case uintptr:
			return x.(uintptr) * y.(uintptr)
		case float32:
			return x.(float32) * y.(float32)
		case float64:
			return x.(float64) * y.(float64)
		case complex64:
			return x.(complex64) * y.(complex64)
		case complex128:
			return x.(complex128) * y.(complex128)
		}

	case token.QUO:
		switch x.(type) {
		case int:
			return x.(int) / y.(int)
		case int8:
			return x.(int8) / y.(int8)
		case int16:
			return x.(int16) / y.(int16)
		case int32:
			return x.(int32) / y.(int32)
		case int64:
			return x.(int64) / y.(int64)
		case uint:
			return x.(uint) / y.(uint)
		case uint8:
			return x.(uint8) / y.(uint8)
		case uint16:
			return x.(uint16) / y.(uint16)
		case uint32:
			return x.(uint32) / y.(uint32)
		case uint64:
			return x.(uint64) / y.(uint64)
		case uintptr:
			return x.(uintptr) / y.(uintptr)
		case float32:
			return x.(float32) / y.(float32)
		case float64:
			return x.(float64) / y.(float64)
		case complex64:
			return x.(complex64) / y.(complex64)
		case complex128:
			return x.(complex128) / y.(complex128)
		}

	case token.REM:
		switch x.(type) {
		case int:
			return x.(int) % y.(int)
		case int8:
			return x.(int8) % y.(int8)
		case int16:
			return x.(int16) % y.(int16)
		case int32:
			return x.(int32) % y.(int32)
		case int64:
			return x.(int64) % y.(int64)
		case uint:
			return x.(uint) % y.(uint)
		case uint8:
			return x.(uint8) % y.(uint8)
		case uint16:
			return x.(uint16) % y.(uint16)
		case uint32:
			return x.(uint32) % y.(uint32)
		case uint64:
			return x.(uint64) % y.(uint64)
		case uintptr:
			return x.(uintptr) % y.(uintptr)
		}

	case token.AND:
		switch x.(type) {
		case int:
			return x.(int) & y.(int)
		case int8:
			return x.(int8) & y.(int8)
		case int16:
			return x.(int16) & y.(int16)
		case int32:
			return x.(int32) & y.(int32)
		case int64:
			return x.(int64) & y.(int64)
		case uint:
			return x.(uint) & y.(uint)
		case uint8:
			return x.(uint8) & y.(uint8)
		case uint16:
			return x.(uint16) & y.(uint16)
		case uint32:
			return x.(uint32) & y.(uint32)
		case uint64:
			return x.(uint64) & y.(uint64)
		case uintptr:
			return x.(uintptr) & y.(uintptr)
		}

	case token.OR:
		switch x.(type) {
		case int:
			return x.(int) | y.(int)
		case int8:
			return x.(int8) | y.(int8)
		case int16:
			return x.(int16) | y.(int16)
		case int32:
			return x.(int32) | y.(int32)
		case int64:
			return x.(int64) | y.(int64)
		case uint:
			return x.(uint) | y.(uint)
		case uint8:
			return x.(uint8) | y.(uint8)
		case uint16:
			return x.(uint16) | y.(uint16)
		case uint32:
			return x.(uint32) | y.(uint32)
		case uint64:
			return x.(uint64) | y.(uint64)
		case uintptr:
			return x.(uintptr) | y.(uintptr)
		}

	case token.XOR:
		switch x.(type) {
		case int:
			return x.(int) ^ y.(int)
		case int8:
			return x.(int8) ^ y.(int8)
		case int16:
			return x.(int16) ^ y.(int16)
		case int32:
			return x.(int32) ^ y.(int32)
		case int64:
			return x.(int64) ^ y.(int64)
		case uint:
			return x.(uint) ^ y.(uint)
		case uint8:
			return x.(uint8) ^ y.(uint8)
		case uint16:
			return x.(uint16) ^ y.(uint16)
		case uint32:
			return x.(uint32) ^ y.(uint32)
		case uint64:
			return x.(uint64) ^ y.(uint64)
		case uintptr:
			return x.(uintptr) ^ y.(uintptr)
		}

	case token.AND_NOT:
		switch x.(type) {
		case int:
			return x.(int) &^ y.(int)
		case int8:
			return x.(int8) &^ y.(int8)
		case int16:
			return x.(int16) &^ y.(int16)
		case int32:
			return x.(int32) &^ y.(int32)
		case int64:
			return x.(int64) &^ y.(int64)
		case uint:
			return x.(uint) &^ y.(uint)
		case uint8:
			return x.(uint8) &^ y.(uint8)
		case uint16:
			return x.(uint16) &^ y.(uint16)
		case uint32:
			return x.(uint32) &^ y.(uint32)
		case uint64:
			return x.(uint64) &^ y.(uint64)
		case uintptr:
			return x.(uintptr) &^ y.(uintptr)
		}

	case token.SHL:
		y := AsUint64(y)
		switch x.(type) {
		case int:
			return x.(int) << y
		case int8:
			return x.(int8) << y
		case int16:
			return x.(int16) << y
		case int32:
			return x.(int32) << y
		case int64:
			return x.(int64) << y
		case uint:
			return x.(uint) << y
		case uint8:
			return x.(uint8) << y
		case uint16:
			return x.(uint16) << y
		case uint32:
			return x.(uint32) << y
		case uint64:
			return x.(uint64) << y
		case uintptr:
			return x.(uintptr) << y
		}

	case token.SHR:
		y := AsUint64(y)
		switch x.(type) {
		case int:
			return x.(int) >> y
		case int8:
			return x.(int8) >> y
		case int16:
			return x.(int16) >> y
		case int32:
			return x.(int32) >> y
		case int64:
			return x.(int64) >> y
		case uint:
			return x.(uint) >> y
		case uint8:
			return x.(uint8) >> y
		case uint16:
			return x.(uint16) >> y
		case uint32:
			return x.(uint32) >> y
		case uint64:
			return x.(uint64) >> y
		case uintptr:
			return x.(uintptr) >> y
		}

	case token.LSS:
		switch x.(type) {
		case int:
			return x.(int) < y.(int)
		case int8:
			return x.(int8) < y.(int8)
		case int16:
			return x.(int16) < y.(int16)
		case int32:
			return x.(int32) < y.(int32)
		case int64:
			return x.(int64) < y.(int64)
		case uint:
			return x.(uint) < y.(uint)
		case uint8:
			return x.(uint8) < y.(uint8)
		case uint16:
			return x.(uint16) < y.(uint16)
		case uint32:
			return x.(uint32) < y.(uint32)
		case uint64:
			return x.(uint64) < y.(uint64)
		case uintptr:
			return x.(uintptr) < y.(uintptr)
		case float32:
			return x.(float32) < y.(float32)
		case float64:
			return x.(float64) < y.(float64)
		case string:
			return x.(string) < y.(string)
		}

	case token.LEQ:
		switch x.(type) {
		case int:
			return x.(int) <= y.(int)
		case int8:
			return x.(int8) <= y.(int8)
		case int16:
			return x.(int16) <= y.(int16)
		case int32:
			return x.(int32) <= y.(int32)
		case int64:
			return x.(int64) <= y.(int64)
		case uint:
			return x.(uint) <= y.(uint)
		case uint8:
			return x.(uint8) <= y.(uint8)
		case uint16:
			return x.(uint16) <= y.(uint16)
		case uint32:
			return x.(uint32) <= y.(uint32)
		case uint64:
			return x.(uint64) <= y.(uint64)
		case uintptr:
			return x.(uintptr) <= y.(uintptr)
		case float32:
			return x.(float32) <= y.(float32)
		case float64:
			return x.(float64) <= y.(float64)
		case string:
			return x.(string) <= y.(string)
		}

	case token.GTR:
		switch x.(type) {
		case int:
			return x.(int) > y.(int)
		case int8:
			return x.(int8) > y.(int8)
		case int16:
			return x.(int16) > y.(int16)
		case int32:
			return x.(int32) > y.(int32)
		case int64:
			return x.(int64) > y.(int64)
		case uint:
			return x.(uint) > y.(uint)
		case uint8:
			return x.(uint8) > y.(uint8)
		case uint16:
			return x.(uint16) > y.(uint16)
		case uint32:
			return x.(uint32) > y.(uint32)
		case uint64:
			return x.(uint64) > y.(uint64)
		case uintptr:
			return x.(uintptr) > y.(uintptr)
		case float32:
			return x.(float32) > y.(float32)
		case float64:
			return x.(float64) > y.(float64)
		case string:
			return x.(string) > y.(string)
		}

	case token.GEQ:
		switch x.(type) {
		case int:
			return x.(int) >= y.(int)
		case int8:
			return x.(int8) >= y.(int8)
		case int16:
			return x.(int16) >= y.(int16)
		case int32:
			return x.(int32) >= y.(int32)
		case int64:
			return x.(int64) >= y.(int64)
		case uint:
			return x.(uint) >= y.(uint)
		case uint8:
			return x.(uint8) >= y.(uint8)
		case uint16:
			return x.(uint16) >= y.(uint16)
		case uint32:
			return x.(uint32) >= y.(uint32)
		case uint64:
			return x.(uint64) >= y.(uint64)
		case uintptr:
			return x.(uintptr) >= y.(uintptr)
		case float32:
			return x.(float32) >= y.(float32)
		case float64:
			return x.(float64) >= y.(float64)
		case string:
			return x.(string) >= y.(string)
		}
	}
	panic(fmt.Sprintf("invalid binary op: %T %s %T", x, op, y))
}
