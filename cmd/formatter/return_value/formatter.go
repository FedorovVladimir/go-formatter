package return_value

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `return value in function

# Example:
func a(a,b int) (c,d int) {} 
-> 
func a(a int, b int) (c int, d int) {}`

var Analyzer = &analysis.Analyzer{
	Name:     "return_value_formatter",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	var result string
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.FuncDecl)
		if e.Type.Params.NumFields() < 2 {
			return
		}
		startPos := e.Type.Results.Pos()
		endPos := e.Type.Results.End()
		results := e.Type.Results.List
		for i := 0; i < len(results); i++ {
			result = result + results[i].Type.(*ast.Ident).Name + ","
		}
		result = "(" + result + ")"
		pass.Report(analysis.Diagnostic{
			Pos:      e.Body.Pos(),
			End:      e.Body.End(),
			Category: "func",
			Message:  "return value",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "return value",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     startPos,
							End:     endPos,
							NewText: []byte(result),
						},
					},
				},
			},
			Related: nil,
		})
	})
	return nil, nil
}
