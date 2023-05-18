package staticlint

import (
	"fmt"

	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var ExitErrCheck = &analysis.Analyzer{
	Name: "ExitErrCheck",
	Doc:  "An analyzer that disables the use of the os.Exit() function in the main block",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {

	for _, f := range pass.Files {
		ast.Inspect(f, func(node ast.Node) bool {
			if x, ok := node.(*ast.CallExpr); ok {
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					if s.Sel.Name == "Exit" {
						fmt.Printf("%v: called not allowed function os.%s()\n", f.Name.Name, s.Sel.Name)
						return false
					}
				}
			}

			return true
		})
	}

	return nil, nil
}
