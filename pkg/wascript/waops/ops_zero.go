// 版权 @2019 凹语言 作者。保留所有权利。

package waops

import (
	"fmt"
	"go/types"
	"unsafe"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

// 生成零值
func Zero(t types.Type) watypes.Value {
	return zero(t)
}

func zero(t types.Type) watypes.Value {
	switch t := t.(type) {
	case *types.Basic:
		if t.Kind() == types.UntypedNil {
			panic("untyped nil has no zero value")
		}
		if t.Info()&types.IsUntyped != 0 {
			// TODO(adonovan): make it an invariant that
			// this is unreachable.  Currently some
			// constants have 'untyped' types when they
			// should be defaulted by the typechecker.
			t = types.Default(t).(*types.Basic)
		}
		switch t.Kind() {
		case types.Bool:
			return false
		case types.Int:
			return int(0)
		case types.Int8:
			return int8(0)
		case types.Int16:
			return int16(0)
		case types.Int32:
			return int32(0)
		case types.Int64:
			return int64(0)
		case types.Uint:
			return uint(0)
		case types.Uint8:
			return uint8(0)
		case types.Uint16:
			return uint16(0)
		case types.Uint32:
			return uint32(0)
		case types.Uint64:
			return uint64(0)
		case types.Uintptr:
			return uintptr(0)
		case types.Float32:
			return float32(0)
		case types.Float64:
			return float64(0)
		case types.Complex64:
			return complex64(0)
		case types.Complex128:
			return complex128(0)
		case types.String:
			return ""
		case types.UnsafePointer:
			return unsafe.Pointer(nil)
		default:
			panic(fmt.Sprint("zero for unexpected type:", t))
		}
	case *types.Pointer:
		return (*watypes.Value)(nil)
	case *types.Array:
		a := make(watypes.Array, t.Len())
		for i := range a {
			a[i] = Zero(t.Elem())
		}
		return a
	case *types.Named:
		return Zero(t.Underlying())
	case *types.Interface:
		return watypes.Iface{} // nil type, methodset and value
	case *types.Slice:
		return []watypes.Value(nil)
	case *types.Struct:
		s := make(watypes.Struct, t.NumFields())
		for i := range s {
			s[i] = Zero(t.Field(i).Type())
		}
		return s
	case *types.Tuple:
		if t.Len() == 1 {
			return Zero(t.At(0).Type())
		}
		s := make(watypes.Tuple, t.Len())
		for i := range s {
			s[i] = Zero(t.At(i).Type())
		}
		return s
	case *types.Chan:
		panic("wa: donot support channel")
	case *types.Map:
		if watypes.UsesBuiltinMap(t.Key()) {
			return map[watypes.Value]watypes.Value(nil)
		}
		return (*watypes.Hashmap)(nil)
	case *types.Signature:
		return (*ssa.Function)(nil)
	}
	panic(fmt.Sprint("zero: unexpected ", t))
}
