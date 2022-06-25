package many_arguments

import (
	"go/ast"
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
	var result string = "\n"
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
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
					result = result + parameters[i].Names[0].Name + " " + x.Name + "." + sel.Name + " ,\n"
				} else {
					result = result + " " + parameters[i].Names[0].Name + " " + x.Name + "." + sel.Name + ",\n"
				}
			case *ast.Ident:
				if i != len(parameters)-1 {
					result = result + parameters[i].Names[0].Name + " " + p.Name + ",\n" + " "
				} else {
					result = result + parameters[i].Names[0].Name + " " + p.Name + ","
				}
			default:
				panic("wow u a lox")
			}

		}
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
