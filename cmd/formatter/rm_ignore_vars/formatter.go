package rm_ignore_vars

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `rm_ignore_vars`

var Analyzer = &analysis.Analyzer{
	Name:     "rm_ignore_vars",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.GenDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.GenDecl)
		v, ok := e.Specs[0].(*ast.ValueSpec)
		if !ok {
			return
		}
		needRm := true
		for _, name := range v.Names {
			if name.Name != "_" {
				needRm = false
			}
		}
		if needRm {
			pass.Report(analysis.Diagnostic{
				Pos:      e.Pos(),
				End:      e.End(),
				Category: "names",
				Message:  "rm_ignore_vars",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "rm_ignore_vars",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     e.Pos(),
								End:     e.End(),
								NewText: []byte(""),
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
