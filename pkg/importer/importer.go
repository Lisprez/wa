// 版权 @2019 凹语言 作者。保留所有权利。

package importer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"github.com/wa-lang/wa/pkg/vfs"
	"github.com/wa-lang/wa/pkg/wa/astcheck"
)

type Importer struct {
	vfs.FileSystem
	Fset  *token.FileSet
	Sizes types.Sizes

	Files map[string][]*ast.File
	Info  map[string]*types.Info
	Pkgs  map[string]*types.Package
}

func New(fset *token.FileSet, fs vfs.FileSystem, sizes types.Sizes) *Importer {
	p := &Importer{
		FileSystem: fs,
		Fset:       fset,
		Sizes:      sizes,
	}

	if p.Fset == nil {
		p.Fset = token.NewFileSet()
	}

	if p.Files == nil {
		p.Files = make(map[string][]*ast.File)
	}
	if p.Info == nil {
		p.Info = make(map[string]*types.Info)
	}
	if p.Pkgs == nil {
		p.Pkgs = make(map[string]*types.Package)
	}

	return p
}

func (p *Importer) Import(pkgpath string) (*types.Package, error) {
	if p.Fset == nil {
		p.Fset = token.NewFileSet()
	}
	if p.Files == nil {
		p.Files = make(map[string][]*ast.File)
	}
	if p.Info == nil {
		p.Info = make(map[string]*types.Info)
	}

	if pkg, ok := p.Pkgs[pkgpath]; ok {
		return pkg, nil
	}

	var files []*ast.File
	for filename, src := range p.FileSystem[pkgpath] {
		f, err := parser.ParseFile(p.Fset, filename, src, parser.AllErrors)
		if err != nil {
			return nil, err
		}
		files = append(files, f)
	}
	if err := astcheck.CheckAST(p.Fset, files...); err != nil {
		return nil, err
	}

	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}

	conf := types.Config{
		Importer: p,
		Sizes:    p.Sizes,
	}
	pkg, err := conf.Check(pkgpath, p.Fset, files, info)
	if err != nil {
		return nil, err
	}

	// OK
	p.Files[pkgpath] = files
	p.Info[pkgpath] = info
	p.Pkgs[pkgpath] = pkg
	return pkg, nil
}
