package grouped_vars

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `grouped vars`

var Analyzer = &analysis.Analyzer{
	Name:     "grouped_vars_formatter",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.GenDecl)(nil)}
	startPos := token.NoPos
	endPos := token.NoPos
	var arr []ast.Spec
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.GenDecl)
		if len(e.Specs) > 1 {
			return
		}
		if startPos == token.NoPos {
			startPos = e.Specs[0].(*ast.ValueSpec).Names[0].Pos()
		}
		endPos = e.Specs[0].End()
		arr = append(arr, e.Specs...)
	})
	if len(arr) == 0 {
		return nil, nil
	}
	var s []string
	for _, spec := range arr {
		var b bytes.Buffer
		_ = printer.Fprint(&b, token.NewFileSet(), spec)
		s = append(s, strings.TrimSpace(b.String()))
	}
	out := "(\n" + strings.Join(s, "\n") + "\n)"
	pass.Report(analysis.Diagnostic{
		Pos:      startPos,
		End:      endPos,
		Category: "names",
		Message:  "grouped vars",
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: "grouped vars",
				TextEdits: []analysis.TextEdit{
					{
						Pos:     startPos,
						End:     endPos,
						NewText: []byte(out),
					},
				},
			},
		},
		Related: nil,
	})
	return nil, nil
}
