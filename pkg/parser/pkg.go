// 版权 @2019 凹语言 作者。保留所有权利。

package parser

import (
	"fmt"
	"go/ast"
	"go/types"
	"io/ioutil"

	"github.com/wa-lang/wa/pkg/importer"
	"github.com/wa-lang/wa/pkg/vfs"
	"github.com/wa-lang/wa/pkg/wa/ssa"
)

// 单个包对象
type Package struct {
	Ctx   *Context       // 环境
	Pkg   *types.Package // 类型检查后的包
	Info  *types.Info    // 包的类型检查信息
	Files []*ast.File    // AST语法树
	SSA   *ssa.Package   // SSA
}

// 加载文件(不支持导入其他包)
func LoadFiles(ctx *Context, files ...string) (*Package, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no wa file")
	}

	var pkgfiles = make(vfs.Files)
	for _, filename := range files {
		src, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		pkgfiles[filename] = string(src)
	}
	if len(pkgfiles) == 0 {
		return nil, fmt.Errorf("no wa file")
	}

	pkgpath := files[0]
	ctx.FileSystem[pkgpath] = pkgfiles

	importer := importer.New(ctx.Fset, ctx.FileSystem, ctx.GetSizes())
	if _, err := importer.Import(pkgpath); err != nil {
		return nil, err
	}

	wpkg := &Package{
		Ctx: ctx,

		Pkg:   importer.Pkgs[pkgpath],
		Info:  importer.Info[pkgpath],
		Files: importer.Files[pkgpath],
	}

	return wpkg, nil
}

// 加载包(不支持导入其他包)
func LoadPackage(ctx *Context, pkgpath string) (*Package, error) {
	importer := importer.New(ctx.Fset, ctx.FileSystem, ctx.GetSizes())
	if _, err := importer.Import(pkgpath); err != nil {
		return nil, err
	}

	wpkg := &Package{
		Ctx: ctx,

		Pkg:   importer.Pkgs[pkgpath],
		Info:  importer.Info[pkgpath],
		Files: importer.Files[pkgpath],
	}

	return wpkg, nil
}

func (p *Package) BuildSSA() *Package {
	prog := ssa.NewProgram(p.Ctx.Fset, ssa.SanityCheckFunctions)
	p.SSA = prog.CreatePackage(p.Pkg, p.Files, p.Info, true)
	p.SSA.Build()
	return p
}
