package arguments_form

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `arguments form in function

# Example:
func a(a,b int) {
} 
-> 
func a(a int, b int) {}`

var Analyzer = &analysis.Analyzer{
	Name:     "arguments_form_formatter",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		result := ""
		e := n.(*ast.FuncDecl)
		if e.Type.Params.NumFields() < 2 {
			return
		}
		startPos := e.Type.Params.Pos()
		endPos := e.Type.Params.End()
		parameters := e.Type.Params.List
		for i := 0; i < len(parameters); i++ {
			for j := 0; j < len(parameters[i].Names); j++ {
				t, ok := parameters[i].Type.(*ast.Ident)
				if ok {
					result = result + parameters[i].Names[j].Name + " " + t.Name + ","
				}
			}
		}
		result = "(" + result + ")"
		pass.Report(analysis.Diagnostic{
			Pos:      e.Body.Pos(),
			End:      e.Body.End(),
			Category: "func",
			Message:  "arguments form",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "arguments form",
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
