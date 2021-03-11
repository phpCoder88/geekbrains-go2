package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println(getGoCount("./gen/async.go", "Counter"))
}

func getGoCount(srcFileName string, funcName string) (int, error) {
	var v visitor

	fSet := token.NewFileSet()
	astFile, err := parser.ParseFile(fSet, srcFileName, nil, 0)
	if err != nil {
		return 0, err
	}

	for _, decl := range astFile.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Name.Name != funcName {
			continue
		}
		spew.Dump(decl)
		ast.Walk(&v, decl)

		break
	}

	return v.counter, nil
}

type visitor struct {
	counter int
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}

	switch node.(type) {
	case *ast.GoStmt:
		v.counter++
		return nil
	}

	return v
}
