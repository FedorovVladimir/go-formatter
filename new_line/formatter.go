package new_line

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `new_line`

var Analyzer = &analysis.Analyzer{
	Name:     "new_line",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.FuncDecl)
		pass.Report(analysis.Diagnostic{
			Pos:      e.Pos(),
			End:      e.Pos(),
			Category: "func",
			Message:  "new_line",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "many arguments",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     e.Pos(),
							End:     e.Pos(),
							NewText: []byte("\n"),
						},
					},
				},
			},
			Related: nil,
		})
	})
	return nil, nil
}
