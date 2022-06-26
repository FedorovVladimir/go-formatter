package context_first_parameter

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `context first parameter in function

# Example:
func a(a int, ctx context.Context) {
} 
-> 
func a(ctx context.Context, a int) {}`

var Analyzer = &analysis.Analyzer{
	Name:     "context_first_parameter_formatter",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		var result string
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
				if x.Name+"."+sel.Name == "context.Context" {
					result = parameters[i].Names[0].Name + " " + x.Name + "." + sel.Name + ", " + result
				} else {
					result = result + " " + parameters[i].Names[0].Name + " " + x.Name + "." + sel.Name + ", "
				}
			case *ast.Ident:
				result = result + parameters[i].Names[0].Name + " " + p.Name + ", "
			default:
				panic("wow u a lox")
			}

		}
		result = strings.TrimSuffix(result, ", ")
		result = "(" + result + ")"
		pass.Report(analysis.Diagnostic{
			Pos:      e.Body.Pos(),
			End:      e.Body.End(),
			Category: "func",
			Message:  "context first parameter",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "context first parameter",
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
