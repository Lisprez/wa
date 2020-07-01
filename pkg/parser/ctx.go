// 版权 @2019 凹语言 作者。保留所有权利。

package parser

import (
	"go/token"
	"go/types"
	"os"
	"runtime"

	"github.com/wa-lang/wa/pkg/vfs"
)

// 构建环境
type Context struct {
	vfs.FileSystem
	Fset *token.FileSet

	WaOS   string
	WaArch string
	WaRoot string

	Sizes     types.Sizes
	DebugMode bool
}

func DefaultContext() *Context {
	p := &Context{
		Fset: token.NewFileSet(),
	}

	if p.WaOS == "" {
		if s := os.Getenv("WaOS"); s != "" {
			p.WaOS = s
		} else {
			p.WaOS = runtime.GOOS
		}
	}
	if p.WaArch == "" {
		if s := os.Getenv("WaArch"); s != "" {
			p.WaArch = s
		} else {
			p.WaArch = runtime.GOARCH
		}
	}
	if p.WaRoot == "" {
		if s := os.Getenv("WaRoot"); s != "" {
			p.WaOS = s
		}
	}

	if p.WaRoot != "" {
		p.FileSystem, _ = vfs.LoadFileSystem(p.WaRoot, vfs.Filter(p.WaOS, p.WaArch))
	}
	if p.FileSystem == nil {
		p.FileSystem = make(vfs.FileSystem)
	}

	p.Sizes = types.SizesFor("gc", p.WaArch)
	return p
}

func (p *Context) GetSizes() types.Sizes {
	if p == nil || p.Sizes == nil {
		return types.SizesFor("gc", p.WaArch)
	} else {
		return p.Sizes
	}
}
