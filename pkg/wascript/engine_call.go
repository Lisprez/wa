// 版权 @2019 凹语言 作者。保留所有权利。

package wascript

import (
	"fmt"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/wabuiltin"
	"github.com/wa-lang/wa/pkg/wascript/waops"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

// 用于生成普通调用或defer调用时的参数
func (p *Engine) prepareCall(fr *Frame, call *ssa.CallCommon) (fn watypes.Value, args []watypes.Value) {
	v := p.getValue(fr, call.Value)

	// 普通函数调用
	if call.Method == nil {
		// Function call.
		fn = v

		// 转换参数, getValue 是核心方法
		for _, arg := range call.Args {
			args = append(args, p.getValue(fr, arg))
		}

		// OK
		return
	}

	// 接口方法调用
	recv := v.(watypes.Iface)
	if recv.T == nil {
		panic("method invoked on nil interface")
	}

	// 查找方法函数
	fn = p.main.Prog.LookupMethod(recv.T, call.Method.Pkg(), call.Method.Name())
	if fn == nil {
		panic(fmt.Sprintf("method set for dynamic type %v does not contain %s", recv.T, call.Method))
	}

	// 第一个参数是接收者信息
	args = append(args, recv.V)

	// 转换参数, getValue 是核心方法
	for _, arg := range call.Args {
		args = append(args, p.getValue(fr, arg))
	}
	return
}

func (i *Engine) call(caller *Frame, fn watypes.Value, args []watypes.Value) watypes.Value {
	// 调用 builtin 函数
	if fn, ok := fn.(*ssa.Builtin); ok {
		return callBuiltin(fn, args)
	}

	// 展开SSA函数对象
	var fnClosure = func() watypes.Closure {
		if x, ok := fn.(*ssa.Function); ok {
			return watypes.Closure{Fn: x}
		} else if closure, ok := fn.(*watypes.Closure); ok {
			return *closure
		} else {
			panic(fmt.Sprintf("cannot call %T", fn))
		}
	}()

	fr := &Frame{fn: fnClosure.Fn}

	// 只有全局函数才可以声明为外部函数
	if fnClosure.Fn.Parent() == nil {
		name := fnClosure.Fn.String()

		// 外部定义的函数
		// 注意外部函数的参数和返回值是如何处理的
		if ext := i.externals[name]; ext != nil {
			return ext(fr, args)
		}

		// 缺少函数body
		// 可以在执行之前进行全局检查, 执行时不需要
		if fnClosure.Fn.Blocks == nil {
			panic("no code for function: " + name)
		}
	}

	// 为当前要调用的函数构造上下文环境
	fr.env = make(map[ssa.Value]watypes.Value)
	fr.block = fnClosure.Fn.Blocks[0]
	fr.locals = make([]watypes.Value, len(fnClosure.Fn.Locals))

	// 将SSA函数的局部变量填充到 frame 中
	for i, l := range fnClosure.Fn.Locals {
		// 先将局部变量零值初始化
		fr.locals[i] = waops.Zero(waops.Deref(l.Type()))

		// 记录到上下文环境
		fr.env[l] = &fr.locals[i]
	}

	// 函数的参数添加到上下文环境
	for i, p := range fnClosure.Fn.Params {
		fr.env[p] = args[i]
	}

	// 要释放的变量也添加到上下文环境
	// 这些可能是闭包函数捕获的外层变量
	for i, fv := range fnClosure.Fn.FreeVars {
		fr.env[fv] = fnClosure.Env[i]
	}

	// 执行函数
	for fr.block != nil {
		i.runFrame(fr) // 核心逻辑
	}

	// 清除布局变量(是否可以省略?), 避免被无意使用
	for i := range fnClosure.Fn.Locals {
		fr.locals[i] = watypes.Bad{}
	}

	// 返回结果
	return fr.result
}

func callBuiltin(fn *ssa.Builtin, args []watypes.Value) watypes.Value {
	switch fn.Name() {
	case "panic", "recover":
		panic("wa do not support panic")
	case "close": // close(chan T)
		panic("wa: donot support channel")

	case "print", "println": // print(any, ...)
		return wabuiltin.Print(fn, args)

	case "len":
		return wabuiltin.Len(fn, args)
	case "cap":
		return wabuiltin.Cap(fn, args)

	case "complex":
		return wabuiltin.Complex(fn, args)
	case "real":
		return wabuiltin.Real(fn, args)
	case "imag":
		return wabuiltin.Imag(fn, args)

	case "append":
		return wabuiltin.Append(fn, args)
	case "copy": // copy([]T, []T) int or copy([]byte, string) int
		return wabuiltin.Copy(fn, args)
	case "delete": // delete(map[K]value, K)
		return wabuiltin.Delete(fn, args)

	case "ssa:wrapnilchk":
		return wabuiltin.SSA_CheckNil(fn, args)
	}

	panic("unknown built-in: " + fn.Name())
}
