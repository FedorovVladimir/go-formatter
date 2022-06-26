package start_enums_at_one

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `start_enums_at_one`

var Analyzer = &analysis.Analyzer{
	Name:     "start_enums_at_one",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.GenDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.GenDecl)
		// e.Specs[0]
		s := e.Specs[0]
		v, ok := s.(*ast.ValueSpec)
		if !ok {
			return
		}
		if len(v.Values) == 1 {
			ident, ok := v.Values[0].(*ast.Ident)
			if !ok {
				return
			}
			if ident.Name == "iota" {
				pass.Report(analysis.Diagnostic{
					Pos:      ident.Pos(),
					End:      ident.End(),
					Category: "names",
					Message:  "start_enums_at_one",
					SuggestedFixes: []analysis.SuggestedFix{
						{
							Message: "rm_ignore_vars",
							TextEdits: []analysis.TextEdit{
								{
									Pos:     ident.Pos(),
									End:     ident.End(),
									NewText: []byte("iota + 1"),
								},
							},
						},
					},
					Related: nil,
				})
			}
		}
	})
	return nil, nil
}
