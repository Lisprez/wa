// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watypes

// Values
//
// All interpreter values are "boxed" in the empty interface, value.
// The range of possible dynamic types within value are:
//
// - bool
// - numbers (all built-in int/float/complex types are distinguished)
// - string
// - map[value]value --- maps for which  usesBuiltinMap(keyType)
//   *hashmap        --- maps for which !usesBuiltinMap(keyType)
// - []value --- slices
// - Iface --- interfaces.
// - Struct --- structs.  Fields are ordered and accessed by numeric indices.
// - array --- arrays.
// - *value --- pointers.  Careful: *value is a distinct type from *array etc.
// - *ssa.Function \
//   *ssa.Builtin   } --- functions.  A nil 'func' is always of type *ssa.Function.
//   *closure      /
// - tuple --- as returned by Return, Next, "value,ok" modes, etc.
// - iter --- iterators from 'range' over map or string.
// - bad --- a poison pill for locals that have gone out of scope.
//
// Note that nil is not on this list.
//
// Pay close attention to whether or not the dynamic type is a pointer.
// The compiler cannot help you since value is an empty interface.

import (
	"bytes"
	"fmt"
	"go/types"
	"sync"
	"unsafe"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wa/types/typeutil"
)

type Value interface{}

type Tuple []Value

type Array []Value

type Iface struct {
	T types.Type // never an "untyped" type
	V Value
}

type Struct []Value

// For map, array, *array, slice, string or channel.
type Iter interface {
	// next returns a Tuple (key, value, ok).
	// key and value are unaliased, e.g. copies of the sequence element.
	Next() Tuple
}

type Closure struct {
	Fn  *ssa.Function
	Env []Value
}

type Bad struct{}

// Hash functions and equivalence relation:

// hashString computes the FNV hash of s.
func HashString(s string) int {
	var h uint32
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return int(h)
}

var (
	mu     sync.Mutex
	hasher = typeutil.MakeHasher()
)

// hashType returns a hash for t such that
// types.Identical(x, y) => hashType(x) == hashType(y).
func HashType(t types.Type) int {
	mu.Lock()
	h := int(hasher.Hash(t))
	mu.Unlock()
	return h
}

// usesBuiltinMap returns true if the built-in hash function and
// equivalence relation for type t are consistent with those of the
// interpreter's representation of type t.  Such types are: all basic
// types (bool, numbers, string), pointers and channels.
//
// usesBuiltinMap returns false for types that require a custom map
// implementation: interfaces, arrays and structs.
//
// Panic ensues if t is an invalid map key type: function, map or slice.
func UsesBuiltinMap(t types.Type) bool {
	switch t := t.(type) {
	case *types.Basic, *types.Chan, *types.Pointer:
		return true
	case *types.Named:
		return UsesBuiltinMap(t.Underlying())
	case *types.Interface, *types.Array, *types.Struct:
		return false
	}
	panic(fmt.Sprintf("invalid map key type: %T", t))
}

func (x Array) Eq(t types.Type, _y interface{}) bool {
	y := _y.(Array)
	tElt := t.Underlying().(*types.Array).Elem()
	for i, xi := range x {
		if !Equals(tElt, xi, y[i]) {
			return false
		}
	}
	return true
}

func (x Array) Hash(t types.Type) int {
	h := 0
	tElt := t.Underlying().(*types.Array).Elem()
	for _, xi := range x {
		h += Hash(tElt, xi)
	}
	return h
}

func (x Struct) Eq(t types.Type, _y interface{}) bool {
	y := _y.(Struct)
	tStruct := t.Underlying().(*types.Struct)
	for i, n := 0, tStruct.NumFields(); i < n; i++ {
		if f := tStruct.Field(i); !f.Anonymous() {
			if !Equals(f.Type(), x[i], y[i]) {
				return false
			}
		}
	}
	return true
}

func (x Struct) Hash(t types.Type) int {
	tStruct := t.Underlying().(*types.Struct)
	h := 0
	for i, n := 0, tStruct.NumFields(); i < n; i++ {
		if f := tStruct.Field(i); !f.Anonymous() {
			h += Hash(f.Type(), x[i])
		}
	}
	return h
}

// nil-tolerant variant of types.Identical.
func SameType(x, y types.Type) bool {
	if x == nil {
		return y == nil
	}
	return y != nil && types.Identical(x, y)
}

func (x Iface) Eq(t types.Type, _y interface{}) bool {
	y := _y.(Iface)
	return SameType(x.T, y.T) && (x.T == nil || Equals(x.T, x.V, y.V))
}

func (x Iface) Hash(_ types.Type) int {
	return HashType(x.T)*8581 + Hash(x.T, x.V)
}

// equals returns true iff x and y are equal according to Go's
// linguistic equivalence relation for type t.
// In a well-typed program, the dynamic types of x and y are
// guaranteed equal.
func Equals(t types.Type, x, y Value) bool {
	switch x := x.(type) {
	case bool:
		return x == y.(bool)
	case int:
		return x == y.(int)
	case int8:
		return x == y.(int8)
	case int16:
		return x == y.(int16)
	case int32:
		return x == y.(int32)
	case int64:
		return x == y.(int64)
	case uint:
		return x == y.(uint)
	case uint8:
		return x == y.(uint8)
	case uint16:
		return x == y.(uint16)
	case uint32:
		return x == y.(uint32)
	case uint64:
		return x == y.(uint64)
	case uintptr:
		return x == y.(uintptr)
	case float32:
		return x == y.(float32)
	case float64:
		return x == y.(float64)
	case complex64:
		return x == y.(complex64)
	case complex128:
		return x == y.(complex128)
	case string:
		return x == y.(string)
	case *Value:
		return x == y.(*Value)
	case Struct:
		return x.Eq(t, y)
	case Array:
		return x.Eq(t, y)
	case Iface:
		return x.Eq(t, y)
	}

	// Since map, func and slice don't support comparison, this
	// case is only reachable if one of x or y is literally nil
	// (handled in eqnil) or via interface{} values.
	panic(fmt.Sprintf("comparing uncomparable type %s", t))
}

// Returns an integer hash of x such that equals(x, y) => hash(x) == hash(y).
func Hash(t types.Type, x Value) int {
	switch x := x.(type) {
	case bool:
		if x {
			return 1
		}
		return 0
	case int:
		return x
	case int8:
		return int(x)
	case int16:
		return int(x)
	case int32:
		return int(x)
	case int64:
		return int(x)
	case uint:
		return int(x)
	case uint8:
		return int(x)
	case uint16:
		return int(x)
	case uint32:
		return int(x)
	case uint64:
		return int(x)
	case uintptr:
		return int(x)
	case float32:
		return int(x)
	case float64:
		return int(x)
	case complex64:
		return int(real(x))
	case complex128:
		return int(real(x))
	case string:
		return HashString(x)
	case *Value:
		return int(uintptr(unsafe.Pointer(x)))
	case Struct:
		return x.Hash(t)
	case Array:
		return x.Hash(t)
	case Iface:
		return x.Hash(t)
	}
	panic(fmt.Sprintf("%T is unhashable", x))
}

// reflect.Value struct values don't have a fixed shape, since the
// payload can be a scalar or an aggregate depending on the instance.
// So store (and load) can't simply use recursion over the shape of the
// rhs value, or the lhs, to copy the value; we need the static type
// information.  (We can't make reflect.Value a new basic data type
// because its "structness" is exposed to Go programs.)

// load returns the value of type T in *addr.
func Load(T types.Type, addr *Value) Value {
	switch T := T.Underlying().(type) {
	case *types.Struct:
		v := (*addr).(Struct)
		a := make(Struct, len(v))
		for i := range a {
			a[i] = Load(T.Field(i).Type(), &v[i])
		}
		return a
	case *types.Array:
		v := (*addr).(Array)
		a := make(Array, len(v))
		for i := range a {
			a[i] = Load(T.Elem(), &v[i])
		}
		return a
	default:
		return *addr
	}
}

// store stores value v of type T into *addr.
func Store(T types.Type, addr *Value, v Value) {
	switch T := T.Underlying().(type) {
	case *types.Struct:
		lhs := (*addr).(Struct)
		rhs := v.(Struct)
		for i := range lhs {
			Store(T.Field(i).Type(), &lhs[i], rhs[i])
		}
	case *types.Array:
		lhs := (*addr).(Array)
		rhs := v.(Array)
		for i := range lhs {
			Store(T.Elem(), &lhs[i], rhs[i])
		}
	default:
		*addr = v
	}
}

// Prints in the style of built-in println.
// (More or less; in gc println is actually a compiler intrinsic and
// can distinguish println(1) from println(interface{}(1)).)
func WriteValue(buf *bytes.Buffer, v Value) {
	switch v := v.(type) {
	case nil, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, complex64, complex128, string:
		fmt.Fprintf(buf, "%v", v)

	case map[Value]Value:
		buf.WriteString("map[")
		sep := ""
		for k, e := range v {
			buf.WriteString(sep)
			sep = " "
			WriteValue(buf, k)
			buf.WriteString(":")
			WriteValue(buf, e)
		}
		buf.WriteString("]")

	case *Hashmap:
		buf.WriteString("map[")
		sep := " "
		for _, e := range v.Entries() {
			for e != nil {
				buf.WriteString(sep)
				sep = " "
				WriteValue(buf, e.Key)
				buf.WriteString(":")
				WriteValue(buf, e.Value)
				e = e.Next
			}
		}
		buf.WriteString("]")

	case *Value:
		if v == nil {
			buf.WriteString("<nil>")
		} else {
			fmt.Fprintf(buf, "%p", v)
		}

	case Iface:
		fmt.Fprintf(buf, "(%s, ", v.T)
		WriteValue(buf, v.V)
		buf.WriteString(")")

	case Struct:
		buf.WriteString("{")
		for i, e := range v {
			if i > 0 {
				buf.WriteString(" ")
			}
			WriteValue(buf, e)
		}
		buf.WriteString("}")

	case Array:
		buf.WriteString("[")
		for i, e := range v {
			if i > 0 {
				buf.WriteString(" ")
			}
			WriteValue(buf, e)
		}
		buf.WriteString("]")

	case []Value:
		buf.WriteString("[")
		for i, e := range v {
			if i > 0 {
				buf.WriteString(" ")
			}
			WriteValue(buf, e)
		}
		buf.WriteString("]")

	case *ssa.Function, *ssa.Builtin, *Closure:
		fmt.Fprintf(buf, "%p", v) // (an address)

	case Tuple:
		// Unreachable in well-formed Go programs
		buf.WriteString("(")
		for i, e := range v {
			if i > 0 {
				buf.WriteString(", ")
			}
			WriteValue(buf, e)
		}
		buf.WriteString(")")

	default:
		fmt.Fprintf(buf, "<%T>", v)
	}
}

// Implements printing of Go values in the style of built-in println.
func ToString(v Value) string {
	var b bytes.Buffer
	WriteValue(&b, v)
	return b.String()
}
