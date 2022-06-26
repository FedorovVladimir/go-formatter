package with

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `rename functions with 'With' path

# Example:
CarWithColor -> WithCarColor`

var Analyzer = &analysis.Analyzer{
	Name:     "with_formatter",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.FuncDecl)
		name := e.Name.Name
		if strings.HasPrefix(name, "With") {
			return
		}
		if !strings.Contains(name, "With") {
			return
		}
		if name[0] <= 'z' && name[0] >= 'a' {
			name = "with" + strings.Title(strings.Replace(name, "With", "", 1))
		} else {
			name = "With" + strings.Replace(name, "With", "", 1)
		}
		pass.Report(analysis.Diagnostic{
			Pos:      e.Name.Pos(),
			End:      e.Name.End(),
			Category: "names",
			Message:  "rename func",
			SuggestedFixes: []analysis.SuggestedFix{
				{
					Message: "rename func",
					TextEdits: []analysis.TextEdit{
						{
							Pos:     e.Name.Pos(),
							End:     e.Name.End(),
							NewText: []byte(name),
						},
					},
				},
			},
			Related: nil,
		})
	})
	return nil, nil
}
