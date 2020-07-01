// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package watypes

// Custom hashtable atop map.
// For use when the key's equivalence relation is not consistent with ==.

// The Go specification doesn't address the atomicity of map operations.
// The FAQ states that an implementation is permitted to crash on
// concurrent map access.

import (
	"go/types"
	"reflect"
)

// A hashtable atop the built-in map.  Since each bucket contains
// exactly one hash value, there's no need to perform hash-equality
// tests when walking the linked list.  Rehashing is done by the
// underlying map.
type Hashmap struct {
	KeyType types.Type
	Table   map[int]*HashmapEntry
	Length  int // number of entries in map
}

type HashmapEntry struct {
	Key   Hashable
	Value Value
	Next  *HashmapEntry
}

type Hashable interface {
	Hash(t types.Type) int
	Eq(t types.Type, x interface{}) bool
}

// makeMap returns an empty initialized map of key type kt,
// preallocating space for reserve elements.
func MakeMap(kt types.Type, reserve int) Value {
	if UsesBuiltinMap(kt) {
		return make(map[Value]Value, reserve)
	}
	return &Hashmap{KeyType: kt, Table: make(map[int]*HashmapEntry, reserve)}
}

// delete removes the association for key k, if any.
func (m *Hashmap) Delete(k Hashable) {
	if m != nil {
		hash := k.Hash(m.KeyType)
		head := m.Table[hash]
		if head != nil {
			if k.Eq(m.KeyType, head.Key) {
				m.Table[hash] = head.Next
				m.Length--
				return
			}
			prev := head
			for e := head.Next; e != nil; e = e.Next {
				if k.Eq(m.KeyType, e.Key) {
					prev.Next = e.Next
					m.Length--
					return
				}
				prev = e
			}
		}
	}
}

// lookup returns the value associated with key k, if present, or
// value(nil) otherwise.
func (m *Hashmap) Lookup(k Hashable) Value {
	if m != nil {
		hash := k.Hash(m.KeyType)
		for e := m.Table[hash]; e != nil; e = e.Next {
			if k.Eq(m.KeyType, e.Key) {
				return e.Value
			}
		}
	}
	return nil
}

// insert updates the map to associate key k with value v.  If there
// was already an association for an eq() (though not necessarily ==)
// k, the previous key remains in the map and its associated value is
// updated.
func (m *Hashmap) Insert(k Hashable, v Value) {
	hash := k.Hash(m.KeyType)
	head := m.Table[hash]
	for e := head; e != nil; e = e.Next {
		if k.Eq(m.KeyType, e.Key) {
			e.Value = v
			return
		}
	}
	m.Table[hash] = &HashmapEntry{
		Key:   k,
		Value: v,
		Next:  head,
	}
	m.Length++
}

// len returns the number of key/value associations in the map.
func (m *Hashmap) Len() int {
	if m != nil {
		return m.Length
	}
	return 0
}

// entries returns a rangeable map of entries.
func (m *Hashmap) Entries() map[int]*HashmapEntry {
	if m != nil {
		return m.Table
	}
	return nil
}

type MapIter struct {
	Iter *reflect.MapIter
	Ok   bool
}

func (it *MapIter) Next() Tuple {
	it.Ok = it.Iter.Next()
	if !it.Ok {
		return []Value{false, nil, nil}
	}
	k, v := it.Iter.Key().Interface(), it.Iter.Value().Interface()
	return []Value{true, k, v}
}

type HashmapIter struct {
	Iter *reflect.MapIter
	Ok   bool
	Cur  *HashmapEntry
}

func (it *HashmapIter) Next() Tuple {
	for {
		if it.Cur != nil {
			k, v := it.Cur.Key, it.Cur.Value
			it.Cur = it.Cur.Next
			return []Value{true, k, v}
		}
		it.Ok = it.Iter.Next()
		if !it.Ok {
			return []Value{false, nil, nil}
		}
		it.Cur = it.Iter.Value().Interface().(*HashmapEntry)
	}
}
