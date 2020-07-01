// 版权 @2019 凹语言 作者。保留所有权利。

package parser

import (
	"github.com/wa-lang/wa/pkg/importer"
	"github.com/wa-lang/wa/pkg/wa/ssa"
)

// 程序对象
// 包含程序需要的全部信息
type Program struct {
	Ctx  *Context     // 环境
	Pkgs []*Package   // 全部包
	SSA  *ssa.Program // SSA
}

// 加载程序
func LoadProgram(ctx *Context, mainPkgPath string) (*Program, error) {
	p := &Program{Ctx: ctx}

	importer := importer.New(p.Ctx.Fset, p.Ctx.FileSystem, p.Ctx.GetSizes())
	if _, err := importer.Import(mainPkgPath); err != nil {
		return nil, err
	}

	for pkgpath, pkg := range importer.Pkgs {
		p.Pkgs = append(p.Pkgs, &Package{
			Ctx:   ctx,
			Pkg:   pkg,
			Info:  importer.Info[pkgpath],
			Files: importer.Files[pkgpath],
		})
	}

	return p, nil
}

func (p *Program) MainPath() string {
	for _, pkg := range p.Pkgs {
		if pkg.Pkg.Name() == "main" {
			return pkg.Pkg.Path()
		}
	}
	return ""
}

func (p *Program) MainPkg() *Package {
	for _, pkg := range p.Pkgs {
		if pkg.Pkg.Name() == "main" {
			return pkg
		}
	}
	return nil
}

func (p *Program) TestMainPath() string {
	for _, pkg := range p.Pkgs {
		if pkg.Pkg.Name() == "testmain" {
			return pkg.Pkg.Path()
		}
	}
	return ""
}

func (p *Program) BuildSSA() *Program {
	prog := ssa.NewProgram(p.Ctx.Fset, ssa.SanityCheckFunctions)
	for _, pkg := range p.Pkgs {
		prog.CreatePackage(pkg.Pkg, pkg.Files, pkg.Info, true)
	}
	p.SSA.Build()
	p.SSA = prog
	return p
}
