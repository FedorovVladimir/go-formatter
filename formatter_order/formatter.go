package formatter_order

import (
	"bytes"
	"go/ast"
	"go/printer"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `formatter_order`

var Analyzer = &analysis.Analyzer{
	Name:     "formatter_order",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type position struct {
	pos token.Pos
	end token.Pos
}

func run(pass *analysis.Pass) (interface{}, error) {
	nodeFilter := []ast.Node{(*ast.GenDecl)(nil), (*ast.FuncDecl)(nil)}
	m := map[token.Token][]ast.Node{}
	var positions []position
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		switch e := n.(type) {
		case *ast.FuncDecl:
			m[token.FUNC] = append(m[token.FUNC], e)
			positions = append(positions, position{pos: e.Pos(), end: e.End()})
		case *ast.GenDecl:
			switch e.Tok {
			case token.CONST, token.VAR, token.TYPE, token.FUNC:
				m[e.Tok] = append(m[e.Tok], e)
				positions = append(positions, position{pos: e.Pos(), end: e.End()})
			}
		}
	})
	i := 0
	if ok, k := work(pass, positions, i, m, token.CONST); ok {
		i = k
	}
	if ok, k := work(pass, positions, i, m, token.VAR); ok {
		i = k
	}
	if ok, k := work(pass, positions, i, m, token.TYPE); ok {
		i = k
	}
	if ok, k := work(pass, positions, i, m, token.FUNC); ok {
		i = k
	}
	return nil, nil
}

func work(pass *analysis.Pass, positions []position, i int, m map[token.Token][]ast.Node, t token.Token) (bool, int) {
	if c, ok := m[t]; ok {
		for _, node := range c {
			var b bytes.Buffer
			_ = printer.Fprint(&b, token.NewFileSet(), node)
			s := b.String()
			pass.Report(analysis.Diagnostic{
				Pos:      positions[i].pos,
				End:      positions[i].end,
				Category: "func",
				Message:  "formatter_order",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "formatter_order",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     positions[i].pos,
								End:     positions[i].end,
								NewText: []byte(s),
							},
						},
					},
				},
				Related: nil,
			})
			i++
		}
		return true, i
	}
	return false, i
}
