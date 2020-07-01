// 版权 @2019 凹语言 作者。保留所有权利。

package waops

import (
	"fmt"
	"go/types"
	"unsafe"

	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

// 强制转换类型
func Conv(t_dst, t_src types.Type, x watypes.Value) watypes.Value {
	// 是否可以通过 ssa.Convert 处理?
	return conv(t_dst, t_src, x)
}

func conv(t_dst, t_src types.Type, x watypes.Value) watypes.Value {
	ut_src := t_src.Underlying()
	ut_dst := t_dst.Underlying()

	// Destination type is not an "untyped" type.
	if b, ok := ut_dst.(*types.Basic); ok && b.Info()&types.IsUntyped != 0 {
		panic("oops: conversion to 'untyped' type: " + b.String())
	}

	// Nor is it an interface type.
	if _, ok := ut_dst.(*types.Interface); ok {
		if _, ok := ut_src.(*types.Interface); ok {
			panic("oops: Convert should be ChangeInterface")
		} else {
			panic("oops: Convert should be MakeInterface")
		}
	}

	// Remaining conversions:
	//    + untyped string/number/bool constant to a specific
	//      representation.
	//    + conversions between non-complex numeric types.
	//    + conversions between complex numeric types.
	//    + integer/[]byte/[]rune -> string.
	//    + string -> []byte/[]rune.
	//
	// All are treated the same: first we extract the value to the
	// widest representation (int64, uint64, float64, complex128,
	// or string), then we convert it to the desired type.

	switch ut_src := ut_src.(type) {
	case *types.Pointer:
		switch ut_dst := ut_dst.(type) {
		case *types.Basic:
			// *value to unsafe.Pointer?
			if ut_dst.Kind() == types.UnsafePointer {
				return unsafe.Pointer(x.(*watypes.Value))
			}
		}

	case *types.Slice:
		// []byte or []rune -> string
		// TODO(adonovan): fix: type B byte; conv([]B -> string).
		switch ut_src.Elem().(*types.Basic).Kind() {
		case types.Byte:
			x := x.([]watypes.Value)
			b := make([]byte, 0, len(x))
			for i := range x {
				b = append(b, x[i].(byte))
			}
			return string(b)

		case types.Rune:
			x := x.([]watypes.Value)
			r := make([]rune, 0, len(x))
			for i := range x {
				r = append(r, x[i].(rune))
			}
			return string(r)
		}

	case *types.Basic:
		x = Widen(x)

		// integer -> string?
		// TODO(adonovan): fix: test integer -> named alias of string.
		if ut_src.Info()&types.IsInteger != 0 {
			if ut_dst, ok := ut_dst.(*types.Basic); ok && ut_dst.Kind() == types.String {
				return string(AsInt(x))
			}
		}

		// string -> []rune, []byte or string?
		if s, ok := x.(string); ok {
			switch ut_dst := ut_dst.(type) {
			case *types.Slice:
				var res []watypes.Value
				// TODO(adonovan): fix: test named alias of rune, byte.
				switch ut_dst.Elem().(*types.Basic).Kind() {
				case types.Rune:
					for _, r := range []rune(s) {
						res = append(res, r)
					}
					return res
				case types.Byte:
					for _, b := range []byte(s) {
						res = append(res, b)
					}
					return res
				}
			case *types.Basic:
				if ut_dst.Kind() == types.String {
					return x.(string)
				}
			}
			break // fail: no other conversions for string
		}

		// unsafe.Pointer -> *value
		if ut_src.Kind() == types.UnsafePointer {
			// TODO(adonovan): this is wrong and cannot
			// really be fixed with the current design.
			//
			// return (*value)(x.(unsafe.Pointer))
			// creates a new pointer of a different
			// type but the underlying interface value
			// knows its "true" type and so cannot be
			// meaningfully used through the new pointer.
			//
			// To make this work, the interpreter needs to
			// simulate the memory layout of a real
			// compiled implementation.
			//
			// To at least preserve type-safety, we'll
			// just return the zero value of the
			// destination type.
			return Zero(t_dst)
		}

		// Conversions between complex numeric types?
		if ut_src.Info()&types.IsComplex != 0 {
			switch ut_dst.(*types.Basic).Kind() {
			case types.Complex64:
				return complex64(x.(complex128))
			case types.Complex128:
				return x.(complex128)
			}
			break // fail: no other conversions for complex
		}

		// Conversions between non-complex numeric types?
		if ut_src.Info()&types.IsNumeric != 0 {
			kind := ut_dst.(*types.Basic).Kind()
			switch x := x.(type) {
			case int64: // signed integer -> numeric?
				switch kind {
				case types.Int:
					return int(x)
				case types.Int8:
					return int8(x)
				case types.Int16:
					return int16(x)
				case types.Int32:
					return int32(x)
				case types.Int64:
					return int64(x)
				case types.Uint:
					return uint(x)
				case types.Uint8:
					return uint8(x)
				case types.Uint16:
					return uint16(x)
				case types.Uint32:
					return uint32(x)
				case types.Uint64:
					return uint64(x)
				case types.Uintptr:
					return uintptr(x)
				case types.Float32:
					return float32(x)
				case types.Float64:
					return float64(x)
				}

			case uint64: // unsigned integer -> numeric?
				switch kind {
				case types.Int:
					return int(x)
				case types.Int8:
					return int8(x)
				case types.Int16:
					return int16(x)
				case types.Int32:
					return int32(x)
				case types.Int64:
					return int64(x)
				case types.Uint:
					return uint(x)
				case types.Uint8:
					return uint8(x)
				case types.Uint16:
					return uint16(x)
				case types.Uint32:
					return uint32(x)
				case types.Uint64:
					return uint64(x)
				case types.Uintptr:
					return uintptr(x)
				case types.Float32:
					return float32(x)
				case types.Float64:
					return float64(x)
				}

			case float64: // floating point -> numeric?
				switch kind {
				case types.Int:
					return int(x)
				case types.Int8:
					return int8(x)
				case types.Int16:
					return int16(x)
				case types.Int32:
					return int32(x)
				case types.Int64:
					return int64(x)
				case types.Uint:
					return uint(x)
				case types.Uint8:
					return uint8(x)
				case types.Uint16:
					return uint16(x)
				case types.Uint32:
					return uint32(x)
				case types.Uint64:
					return uint64(x)
				case types.Uintptr:
					return uintptr(x)
				case types.Float32:
					return float32(x)
				case types.Float64:
					return float64(x)
				}
			}
		}
	}

	panic(fmt.Sprintf("unsupported conversion: %s  -> %s, dynamic type %T", t_src, t_dst, x))
}
