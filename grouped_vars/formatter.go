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

type group struct {
	pos   token.Pos
	end   token.Pos
	specs []ast.Spec
}

var Analyzer = &analysis.Analyzer{
	Name:     "grouped_vars_formatter",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.GenDecl)(nil)}
	files := map[string][]group{}
	oldFilename := ""
	oldLine := 0
	i := 0
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.GenDecl)

		filename := pass.Fset.Position(e.Pos()).Filename
		if filename != oldFilename {
			i = 0
			pos := e.Lparen
			if pos == token.NoPos {
				pos = e.Specs[0].Pos()
			}
			end := e.Rparen
			if end == token.NoPos {
				end = e.Specs[len(e.Specs)-1].Pos()
			}
			files[filename] = append(
				files[filename],
				group{
					pos:   pos,
					end:   end,
					specs: nil,
				},
			)
		}
		oldFilename = filename
		line := pass.Fset.Position(e.Pos()).Line

		if line-oldLine == 1 {
			files[filename][i].specs = append(files[filename][i].specs, e.Specs...)
			end := e.Rparen + 1
			if e.Rparen == token.NoPos {
				end = e.Specs[len(e.Specs)-1].End()
			}
			files[filename][i].end = end
		}
		if line-oldLine > 1 && oldLine != 0 {
			i++
			pos := e.Lparen
			if e.Lparen == token.NoPos {
				pos = e.Specs[0].Pos()
			}
			end := e.Rparen + 1
			if e.Rparen == token.NoPos {
				end = e.Specs[len(e.Specs)-1].End()
			}
			files[filename] = append(
				files[filename], group{
					pos:   pos,
					end:   end,
					specs: nil,
				},
			)
		}

		if len(files[filename][i].specs) == 0 {
			files[filename][i].specs = append(files[filename][i].specs, e.Specs...)
			end := e.Rparen + 1
			if e.Rparen == token.NoPos {
				end = e.Specs[len(e.Specs)-1].Pos()
			}
			files[filename][i].end = end
		}
		oldLine = pass.Fset.Position(e.End()).Line
	})
	for _, file := range files {
		for _, specs := range file {
			if len(specs.specs) == 1 {
				var b bytes.Buffer
				_ = printer.Fprint(&b, token.NewFileSet(), specs.specs[0])
				l1 := pass.Fset.Position(specs.pos).Line
				l2 := pass.Fset.Position(specs.end).Line
				if l1 == l2 {
					continue
				}
				out := b.String()
				pass.Report(analysis.Diagnostic{
					Pos:      specs.pos,
					End:      specs.end,
					Category: "names",
					Message:  "grouped vars",
					SuggestedFixes: []analysis.SuggestedFix{
						{
							Message: "grouped vars",
							TextEdits: []analysis.TextEdit{
								{
									Pos:     specs.pos,
									End:     specs.end,
									NewText: []byte(out),
								},
							},
						},
					},
					Related: nil,
				})
			}
			if len(specs.specs) < 2 {
				continue
			}
			var s []string
			for _, spec := range specs.specs {
				var b bytes.Buffer
				_ = printer.Fprint(&b, token.NewFileSet(), spec)
				s = append(s, strings.TrimSpace(b.String()))
			}
			out := "(\n" + strings.Join(s, "\n") + "\n)"
			pass.Report(analysis.Diagnostic{
				Pos:      specs.pos,
				End:      specs.end,
				Category: "names",
				Message:  "grouped vars",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "grouped vars",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     specs.pos,
								End:     specs.end,
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
