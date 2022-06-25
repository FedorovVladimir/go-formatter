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
	files := map[string][][]ast.Spec{}
	oldFilename := ""
	oldLine := 0
	i := 0
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.GenDecl)
		if len(e.Specs) > 1 {
			return
		}

		filename := pass.Fset.Position(e.Pos()).Filename
		if filename != oldFilename {
			i = 0
			files[filename] = append(files[filename], []ast.Spec{})
		}
		oldFilename = filename
		line := pass.Fset.Position(e.Pos()).Line

		if line-oldLine == 1 {
			files[filename][i] = append(files[filename][i], e.Specs...)
		}
		if line-oldLine > 1 && oldLine != 0 {
			i++
			files[filename] = append(files[filename], []ast.Spec{})
		}

		if len(files[filename][i]) == 0 {
			files[filename][i] = append(files[filename][i], e.Specs...)
		}
		oldLine = line
	})
	for _, file := range files {
		for _, specs := range file {
			if len(specs) < 2 {
				return nil, nil
			}
			var s []string
			for _, spec := range specs {
				var b bytes.Buffer
				_ = printer.Fprint(&b, token.NewFileSet(), spec)
				s = append(s, strings.TrimSpace(b.String()))
			}
			out := "(\n" + strings.Join(s, "\n") + "\n)"
			pass.Report(analysis.Diagnostic{
				Pos:      specs[0].Pos(),
				End:      specs[len(specs)-1].End(),
				Category: "names",
				Message:  "grouped vars",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "grouped vars",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     specs[0].Pos(),
								End:     specs[len(specs)-1].End(),
								NewText: []byte(out),
							},
						},
					},
				},
				Related: nil,
			})
		}
	}
	return nil, nil
}
