// 版权 @2019 凹语言 作者。保留所有权利。

package wa

import (
	"fmt"
	"io/ioutil"

	"github.com/wa-lang/wa/pkg/parser"
	"github.com/wa-lang/wa/pkg/vfs"
	"github.com/wa-lang/wa/pkg/wascript"
	"github.com/wa-lang/wa/pkg/wascript/watypes"
)

// 凹脚本引擎
type WaScript struct {
	opt options
	ctx *parser.Context
	pkg *parser.Package

	funcs map[string]wascript.UserFunc
}

// 凹脚本函数帧类型
type Frame = wascript.Frame

// 凹值类型
type Value = watypes.Value

// 创建新的凹脚本引擎
func NewScript(opts ...Options) *WaScript {
	p := &WaScript{
		funcs: make(map[string]wascript.UserFunc),
	}
	p.opt.applyOptions(opts...)
	if p.ctx == nil {
		p.ctx = parser.DefaultContext()
	}
	return p
}

// 释放脚本
func (p *WaScript) Close() error {
	return nil
}

// 定义外部函数
func (p *WaScript) DefineFunc(name string, fn wascript.UserFunc) *WaScript {
	p.funcs[name] = fn
	return p
}

// 读取脚本
// 支持读取单个文件
// 支持加载整个子文件系统(类似mount行为)
func (p *WaScript) Load(path string, _data ...interface{}) error {
	if path == "" {
		return fmt.Errorf("invalid path: %s", path)
	}

	pkgpath := path
	if len(_data) > 0 {
		switch x := _data[0].(type) {
		case string:
			p.ctx.FileSystem[path] = map[string]string{
				path: x,
			}
		case []byte:
			p.ctx.FileSystem[path] = map[string]string{
				path: string(x),
			}
		case map[string]string:
			p.ctx.FileSystem[path] = x
		case vfs.Files:
			p.ctx.FileSystem[path] = x
		case map[string]map[string]string:
			for k, v := range x {
				p.ctx.FileSystem[k] = v
			}
		case vfs.FileSystem:
			for k, v := range x {
				p.ctx.FileSystem[k] = v
			}
		}
	} else {
		if isFile(path) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			p.ctx.FileSystem[pkgpath] = map[string]string{
				"path": string(data),
			}
		} else {
			fs, err := vfs.LoadFileSystem(path, nil)
			if err != nil {
				return err
			}
			for k, v := range fs {
				p.ctx.FileSystem[k] = v
			}
		}
	}

	pkg, err := parser.LoadPackage(p.ctx, pkgpath)
	if err != nil {
		return err
	}

	p.pkg = pkg.BuildSSA()
	return nil
}

// 读取脚本, 链式操作, 失败抛异常
func (p *WaScript) MustLoad(path string, data ...interface{}) *WaScript {
	if err := p.Load(path, data...); err != nil {
		panic(err)
	}
	return p
}

// 调用脚本内的函数
func (p *WaScript) Call(fn string, args ...Value) (Value, error) {
	return nil, fmt.Errorf("TODO")
}

// 执行脚本(main.main)
func (p *WaScript) Run(args ...string) error {
	return wascript.NewEngine(p.pkg.SSA, p.funcs).Run(args...)
}
