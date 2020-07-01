// 版权 @2019 凹语言 作者。保留所有权利。

// 凹对应的SSA相关的指令操作.
package waops

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"
	"unsafe"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

// 查找 map 或 string 中的元素.
func Lookup(instr *ssa.Lookup, x, idx watypes.Value) watypes.Value {
	switch x := x.(type) { // map or string
	case string:
		return x[AsInt(idx)]
	case map[watypes.Value]watypes.Value:
		v, ok := x[idx]
		if !ok {
			v = Zero(instr.X.Type().Underlying().(*types.Map).Elem())
		}
		if instr.CommaOk {
			v = watypes.Tuple{v, ok}
		}
		return v
	case *watypes.Hashmap:
		v := x.Lookup(idx.(watypes.Hashable))
		ok := v != nil
		if !ok {
			v = Zero(instr.X.Type().Underlying().(*types.Map).Elem())
		}
		if instr.CommaOk {
			v = watypes.Tuple{v, ok}
		}
		return v
	}
	panic(fmt.Sprintf("unexpected x type in Lookup: %T", x))
}

// 迭代 map 和 string (不支持 slice 和 array ?)
func RangeIter(x watypes.Value) watypes.Iter {
	switch x := x.(type) {
	case map[watypes.Value]watypes.Value:
		return &watypes.MapIter{Iter: reflect.ValueOf(x).MapRange()}
	case *watypes.Hashmap:
		return &watypes.HashmapIter{Iter: reflect.ValueOf(x.Entries()).MapRange()}
	case string:
		return &watypes.StringIter{Reader: strings.NewReader(x)}
	}
	panic(fmt.Sprintf("cannot range over %T", x))
}

// 将基础类型扩展对最大精度类型, 比如浮点数扩展为 float64 等.
func Widen(x watypes.Value) watypes.Value {
	switch y := x.(type) {
	case bool, int64, uint64, float64, complex128, string, unsafe.Pointer:
		return x
	case int:
		return int64(y)
	case int8:
		return int64(y)
	case int16:
		return int64(y)
	case int32:
		return int64(y)
	case uint:
		return uint64(y)
	case uint8:
		return uint64(y)
	case uint16:
		return uint64(y)
	case uint32:
		return uint64(y)
	case uintptr:
		return uint64(y)
	case float32:
		return float64(y)
	case complex64:
		return complex128(y)
	}
	panic(fmt.Sprintf("cannot widen %T", x))
}
