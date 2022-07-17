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

type decl int

const (
	constDecl decl = iota + 1
	varDecl
	typeDecl
	funcDecl
)

var orderDecl = []decl{constDecl, varDecl, typeDecl, funcDecl}

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
	m := map[decl][]ast.Node{}
	var positions []position
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		switch e := n.(type) {
		case *ast.FuncDecl:
			m[funcDecl] = append(m[funcDecl], e)
			positions = append(positions, position{pos: e.Pos(), end: e.End()})
		case *ast.GenDecl:
			switch e.Tok {
			case token.CONST:
				m[constDecl] = append(m[constDecl], e)
				positions = append(positions, position{pos: e.Pos(), end: e.End()})
			case token.VAR:
				m[varDecl] = append(m[varDecl], e)
				positions = append(positions, position{pos: e.Pos(), end: e.End()})
			case token.TYPE:
				m[typeDecl] = append(m[typeDecl], e)
				positions = append(positions, position{pos: e.Pos(), end: e.End()})
			case token.FUNC:
				m[funcDecl] = append(m[funcDecl], e)
				positions = append(positions, position{pos: e.Pos(), end: e.End()})
			}
		}
	})

	i := 0
	for _, decl := range orderDecl {
		if ok, k := work(pass, positions, i, m, decl); ok {
			i = k
		}
	}

	return nil, nil
}

func work(pass *analysis.Pass, positions []position, i int, m map[decl][]ast.Node, decl decl) (bool, int) {
	if c, ok := m[decl]; ok {
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
