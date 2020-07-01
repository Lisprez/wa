// 版权 @2019 凹语言 作者。保留所有权利。

// 凹脚本语言解释器.
package wascript

import (
	"fmt"
	"go/types"
	"io"
	"sync"

	"github.com/wa-lang/wa/pkg/wa/ssa"
	"github.com/wa-lang/wa/pkg/wascript/waops"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

// 凹脚本语言解释器
type Engine struct {
	Args []string // 命令行参数
	Env  []string // 环境变量
	Dir  string   // 工作目录

	Stdin  io.Reader // 标准输入设备
	Stdout io.Writer // 标准输出设备
	Stderr io.Writer // 标准错误输出设置

	main     *ssa.Package
	args     []watypes.Value
	initOnce sync.Once

	// 全局变量(地址是固定的)
	globals map[string]*watypes.Value

	// 外部注入的函数
	externals map[string]UserFunc

	errorMethods map[string]*ssa.Function
	rtypeMethods map[string]*ssa.Function
}

// 函数调用帧
type Frame struct {
	fn     *ssa.Function
	result watypes.Value

	env       map[ssa.Value]watypes.Value
	locals    []watypes.Value
	deferFn   []watypes.Value
	deferArgs [][]watypes.Value

	block     *ssa.BasicBlock
	prevBlock *ssa.BasicBlock
}

// 外部函数类型
type UserFunc func(fr *Frame, args ...watypes.Value) watypes.Value

// 构造引擎
func NewEngine(mainpkg *ssa.Package, funcs map[string]UserFunc) *Engine {
	p := &Engine{
		main:      mainpkg,
		globals:   make(map[string]*watypes.Value),
		externals: make(map[string]UserFunc),
	}
	for k, fn := range funcs {
		p.externals[k] = fn
	}
	return p
}

// 读全局变量
func (p *Engine) getGlobal(key *ssa.Global) (v *watypes.Value, ok bool) {
	v, ok = p.globals[key.RelString(nil)]
	return
}

// 设置全局变量(初始化零值)
func (p *Engine) setGlobal(key *ssa.Global, v *watypes.Value) {
	p.globals[key.RelString(nil)] = v
}

// 读取值(nil/常量/全局变量/函数/局部变量)
// SSA形式已经不存在作用域概念, 因此每个分支都是唯一的结果, 查询的顺序并无关系
func (p *Engine) getValue(fr *Frame, key ssa.Value) watypes.Value {
	switch key := key.(type) {
	case nil:
		// Hack; simplifies handling of optional attributes
		// such as ssa.Slice.{Low,High}.
		return nil
	case *ssa.Const:
		return waops.ConstValue(key)
	case *ssa.Function, *ssa.Builtin:
		return key
	case *ssa.Global:
		if r, ok := p.getGlobal(key); ok {
			return r
		}
	}

	// 局部变量需要封装为方法
	if r, ok := fr.env[key]; ok {
		return r
	}

	panic(fmt.Sprintf("get: no value for %T: %v", key, key.Name()))
}

func (p *Engine) Run(args ...string) error {
	return p.initPkgs().run(args)
}

func (p *Engine) Invoke(fnName string, args watypes.Value) (watypes.Value, error) {
	return p.initPkgs().invoke(fnName, args)
}

func (p *Engine) initPkgs() *Engine {
	p.initOnce.Do(func() {
		for _, pkg := range p.main.Prog.AllPackages() {
			for _, m := range pkg.Members {
				switch v := m.(type) {
				case *ssa.Global:
					cell := waops.Zero(waops.Deref(v.Type()))
					p.setGlobal(v, &cell)
				}
			}
		}

		p.call(nil, p.main.Func("init"), nil)
	})
	return p
}

func (p *Engine) run(args []string) error {
	for _, arg := range args {
		p.args = append(p.args, arg)
	}

	mainFn := p.main.Func("main")
	if mainFn == nil {
		return fmt.Errorf("No main function.")
	}

	p.call(nil, mainFn, nil)
	return nil
}

func (p *Engine) invoke(fnName string, args watypes.Value) (watypes.Value, error) {
	return nil, fmt.Errorf("TODO")
}

func (p *Engine) typeAssert(instr *ssa.TypeAssert, itf watypes.Iface) watypes.Value {
	var v watypes.Value
	err := ""
	if itf.T == nil {
		err = fmt.Sprintf("interface conversion: interface is nil, not %s", instr.AssertedType)

	} else if idst, ok := instr.AssertedType.Underlying().(*types.Interface); ok {
		v = itf
		if meth, _ := types.MissingMethod(itf.T, idst, true); meth != nil {
			err = fmt.Sprintf(
				"interface conversion: %v is not %v: missing method %s",
				itf.T, idst, meth.Name(),
			)
		}

	} else if types.Identical(itf.T, instr.AssertedType) {
		v = itf.V // extract value

	} else {
		err = fmt.Sprintf("interface conversion: interface is %s, not %s", itf.T, instr.AssertedType)
	}

	if err != "" {
		if !instr.CommaOk {
			panic(err)
		}
		return watypes.Tuple{waops.Zero(instr.AssertedType), false}
	}
	if instr.CommaOk {
		return watypes.Tuple{v, true}
	}
	return v
}
