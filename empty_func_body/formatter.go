package empty_func_body

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `empty functions body

# Example:
func a() {
} 
-> 
func a() {}`

var Analyzer = &analysis.Analyzer{
	Name:     "empty_func_body_formatter",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.FuncDecl)
		if e.Body.End()-e.Body.Pos() == 2 {
			return
		}
		if e.Body.List == nil {
			pass.Report(analysis.Diagnostic{
				Pos:      e.Body.Pos(),
				End:      e.Body.End(),
				Category: "func",
				Message:  "empty func",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "empty func",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     e.Body.Pos(),
								End:     e.Body.End(),
								NewText: []byte("{}"),
							},
						},
					},
				},
				Related: nil,
			})
		}
	})
	return nil, nil
}
