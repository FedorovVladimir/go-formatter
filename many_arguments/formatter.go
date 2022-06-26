package many_arguments

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `context many arguments in function

# Example:
func a(a int, b int, s string, b bool) {
} 
-> 
func a(a int, 
	b int, 
	s string, 
	b bool) {
}`

var Analyzer = &analysis.Analyzer{
	Name:     "many_arguments_formatter",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		result := "\n"
		e := n.(*ast.FuncDecl)
		if e.Type.Params.NumFields() < 2 {
			return
		}
		startPos := e.Type.Params.Pos()
		endPos := e.Type.Params.End()
		parameters := e.Type.Params.List
		for i := 0; i < len(parameters); i++ {
			switch p := parameters[i].Type.(type) {
			case *ast.SelectorExpr:
				x, ok := p.X.(*ast.Ident)
				if !ok {
					continue
				}
				sel := p.Sel
				result = result + parameters[i].Names[0].Name + " " + x.Name + "." + sel.Name + ",\n"
			case *ast.Ident:
				result = result + parameters[i].Names[0].Name + " " + p.Name + ",\n"
			default:
				panic("wow u a lox")
			}
		}
		result = strings.TrimSuffix(result, "\n")
		result = "(" + result + "\n)"
		pass.Report(analysis.Diagnostic{
			Pos:      e.Body.Pos(),
			End:      e.Body.End(),
			Category: "func",
			Message:  "many arguments",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "many arguments",
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
