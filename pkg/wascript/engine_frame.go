// 版权 @2019 凹语言 作者。保留所有权利。

package wascript

import (
	"fmt"
	"go/types"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/waops"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

func (p *Engine) runFrame(fr *Frame) {
	for i := 0; i < len(fr.block.Instrs); i++ {
		switch instr := fr.block.Instrs[i].(type) {
		case *ssa.DebugRef:
			// no-op
		case *ssa.Go:
			panic("wa: donot support goroutine")
		case *ssa.MakeChan:
			panic("wa: donot support channel")
		case *ssa.Select:
			panic("wa do not support channel")
		case *ssa.Send:
			panic("wa: donot support channel")
		case *ssa.Panic:
			panic("wa: donot support panic")

		case *ssa.Alloc:
			var addr *watypes.Value
			if instr.Heap {
				addr = new(watypes.Value)
				fr.env[instr] = addr
			} else {
				addr = fr.env[instr].(*watypes.Value)
			}
			*addr = waops.Zero(waops.Deref(instr.Type()))

		case *ssa.Store:
			watypes.Store(waops.Deref(instr.Addr.Type()), p.getValue(fr, instr.Addr).(*watypes.Value), p.getValue(fr, instr.Val))

		case *ssa.Call:
			fn, args := p.prepareCall(fr, &instr.Call)
			fr.env[instr] = p.call(fr, fn, args)

		case *ssa.RunDefers:
			for i := len(fr.deferFn) - 1; i >= 0; i-- {
				p.call(fr, fr.deferFn[i], fr.deferArgs[i])
			}
			fr.deferFn = nil

		case *ssa.Return:
			switch len(instr.Results) {
			case 0:
			case 1:
				fr.result = p.getValue(fr, instr.Results[0])
			default:
				var res []watypes.Value
				for _, r := range instr.Results {
					res = append(res, p.getValue(fr, r))
				}
				fr.result = watypes.Tuple(res)
			}
			fr.block = nil
			return

		case *ssa.If:
			if p.getValue(fr, instr.Cond).(bool) {
				fr.prevBlock, fr.block = fr.block, fr.block.Succs[0] // true
				i = 0 - 1
				continue
			} else {
				fr.prevBlock, fr.block = fr.block, fr.block.Succs[1] // false
				i = 0 - 1
				continue
			}

		case *ssa.Jump:
			fr.prevBlock, fr.block = fr.block, fr.block.Succs[0]
			i = 0 - 1
			continue

		case *ssa.Phi:
			for i, pred := range instr.Block().Preds {
				if fr.prevBlock == pred {
					fr.env[instr] = p.getValue(fr, instr.Edges[i])
					break
				}
			}

		case *ssa.UnOp:
			fr.env[instr] = waops.UnOp(instr, p.getValue(fr, instr.X))
		case *ssa.BinOp:
			fr.env[instr] = waops.BinOp(instr.Op, instr.X.Type(), p.getValue(fr, instr.X), p.getValue(fr, instr.Y))

		case *ssa.MakeClosure:
			var bindings []watypes.Value
			for _, binding := range instr.Bindings {
				bindings = append(bindings, p.getValue(fr, binding))
			}
			fr.env[instr] = &watypes.Closure{
				Fn:  instr.Fn.(*ssa.Function),
				Env: bindings,
			}

		case *ssa.MakeInterface:
			fr.env[instr] = watypes.Iface{T: instr.X.Type(), V: p.getValue(fr, instr.X)}
		case *ssa.Lookup:
			fr.env[instr] = waops.Lookup(instr, p.getValue(fr, instr.X), p.getValue(fr, instr.Index))
		case *ssa.TypeAssert:
			fr.env[instr] = p.typeAssert(instr, p.getValue(fr, instr.X).(watypes.Iface))
		case *ssa.ChangeInterface:
			fr.env[instr] = p.getValue(fr, instr.X)
		case *ssa.ChangeType:
			fr.env[instr] = p.getValue(fr, instr.X) // (can't fail)
		case *ssa.Convert:
			fr.env[instr] = waops.Conv(instr.Type(), instr.X.Type(), p.getValue(fr, instr.X))
		case *ssa.Extract:
			fr.env[instr] = p.getValue(fr, instr.Tuple).(watypes.Tuple)[instr.Index]

		case *ssa.Range:
			fr.env[instr] = waops.RangeIter(p.getValue(fr, instr.X))
		case *ssa.Next:
			fr.env[instr] = p.getValue(fr, instr.Iter).(watypes.Iter).Next()

		case *ssa.Defer:
			fn, args := p.prepareCall(fr, &instr.Call)
			fr.deferFn = append(fr.deferFn, fn)
			fr.deferArgs = append(fr.deferArgs, args)

		case *ssa.MakeSlice:
			slice := make([]watypes.Value, waops.AsInt(p.getValue(fr, instr.Cap)))
			tElt := instr.Type().Underlying().(*types.Slice).Elem()
			for i := range slice {
				slice[i] = waops.Zero(tElt)
			}
			fr.env[instr] = slice[:waops.AsInt(p.getValue(fr, instr.Len))]

		case *ssa.Slice:
			fr.env[instr] = waops.Slice(p.getValue(fr, instr.X), p.getValue(fr, instr.Low), p.getValue(fr, instr.High), p.getValue(fr, instr.Max))

		case *ssa.Field:
			fr.env[instr] = p.getValue(fr, instr.X).(watypes.Struct)[instr.Field]
		case *ssa.FieldAddr:
			fr.env[instr] = &(*p.getValue(fr, instr.X).(*watypes.Value)).(watypes.Struct)[instr.Field]

		case *ssa.Index:
			fr.env[instr] = p.getValue(fr, instr.X).(watypes.Array)[waops.AsInt(p.getValue(fr, instr.Index))]
		case *ssa.IndexAddr:
			x := p.getValue(fr, instr.X)
			idx := p.getValue(fr, instr.Index)
			switch x := x.(type) {
			case []watypes.Value:
				fr.env[instr] = &x[waops.AsInt(idx)]
			case *watypes.Value: // *array
				fr.env[instr] = &(*x).(watypes.Array)[waops.AsInt(idx)]
			default:
				panic(fmt.Sprintf("unexpected x type in IndexAddr: %T", x))
			}

		case *ssa.MakeMap:
			reserve := 0
			if instr.Reserve != nil {
				reserve = waops.AsInt(p.getValue(fr, instr.Reserve))
			}
			fr.env[instr] = watypes.MakeMap(instr.Type().Underlying().(*types.Map).Key(), reserve)

		case *ssa.MapUpdate:
			m := p.getValue(fr, instr.Map)
			key := p.getValue(fr, instr.Key)
			v := p.getValue(fr, instr.Value)
			switch m := m.(type) {
			case map[watypes.Value]watypes.Value:
				m[key] = v
			case *watypes.Hashmap:
				m.Insert(key.(watypes.Hashable), v)
			default:
				panic(fmt.Sprintf("illegal map type: %T", m))
			}

		default:
			panic(fmt.Sprintf("unexpected instruction: %T", instr))
		}
	}
}
