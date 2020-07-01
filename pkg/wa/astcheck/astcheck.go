// 版权 @2019 凹语言 作者。保留所有权利。

package astcheck

import (
	"fmt"
	"go/ast"
	"go/token"
)

func CheckAST(fset *token.FileSet, files ...*ast.File) error {
	var err error
	for _, f := range files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.GoStmt:
				err = fmt.Errorf("%v: wa do not support goroutine", fset.Position(x.Pos()))
				return false
			case *ast.ChanType:
				err = fmt.Errorf("%v: wa do not support channel", fset.Position(x.Pos()))
				return false
			case *ast.SelectStmt:
				err = fmt.Errorf("%v: wa do not support channel", fset.Position(x.Pos()))
				return false
			case *ast.SendStmt:
				err = fmt.Errorf("%v: wa do not support channel", fset.Position(x.Pos()))
				return false
			}
			return true
		})
		if err != nil {
			return err
		}
	}
	return nil
}
