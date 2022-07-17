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
	groups := map[decl][]ast.Node{}
	var positions []*position

	var lastPosition *position
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(
		[]ast.Node{(*ast.GenDecl)(nil), (*ast.FuncDecl)(nil)},
		func(n ast.Node) {
			switch e := n.(type) {
			case *ast.FuncDecl:
				groups[funcDecl] = append(groups[funcDecl], e)
				if lastPosition != nil {
					lastPosition.end = e.Pos()
				}
				lastPosition = &position{pos: e.Pos(), end: e.End()}
				positions = append(positions, lastPosition)
			case *ast.GenDecl:
				switch e.Tok {
				case token.CONST:
					groups[constDecl] = append(groups[constDecl], e)
					if lastPosition != nil {
						lastPosition.end = e.Pos()
					}
					lastPosition = &position{pos: e.Pos(), end: e.End()}
					positions = append(positions, lastPosition)
				case token.VAR:
					groups[varDecl] = append(groups[varDecl], e)
					if lastPosition != nil {
						lastPosition.end = e.Pos()
					}
					lastPosition = &position{pos: e.Pos(), end: e.End()}
					positions = append(positions, lastPosition)
				case token.TYPE:
					groups[typeDecl] = append(groups[typeDecl], e)
					if lastPosition != nil {
						lastPosition.end = e.Pos()
					}
					lastPosition = &position{pos: e.Pos(), end: e.End()}
					positions = append(positions, lastPosition)
				case token.FUNC:
					groups[funcDecl] = append(groups[funcDecl], e)
					if lastPosition != nil {
						lastPosition.end = e.Pos()
					}
					lastPosition = &position{pos: e.Pos(), end: e.End()}
					positions = append(positions, lastPosition)
				}
			}
		},
	)
	end := token.Pos(pass.Fset.File(lastPosition.end).Size())
	lastPosition.end = end

	i := 0
	for _, decl := range orderDecl {
		if ok, k := reportGroup(pass, positions, i, groups, decl); ok {
			i = k
		}
	}

	return nil, nil
}

func reportGroup(pass *analysis.Pass, positions []*position, i int, groups map[decl][]ast.Node, decl decl) (bool, int) {
	if nodes, ok := groups[decl]; ok {
		for _, node := range nodes {
			var b bytes.Buffer
			_ = printer.Fprint(&b, token.NewFileSet(), node)
			b.Write([]byte("\n"))
			report(pass, positions[i].pos, positions[i].end, b.Bytes(), "formatter_order")
			i++
		}
		return true, i
	}
	return false, i
}

func report(pass *analysis.Pass, pos token.Pos, end token.Pos, text []byte, msg string) {
	pass.Report(analysis.Diagnostic{
		Pos:      pos,
		End:      end,
		Category: msg,
		Message:  msg,
		SuggestedFixes: []analysis.SuggestedFix{
			{
				Message: msg,
				TextEdits: []analysis.TextEdit{
					{
						Pos:     pos,
						End:     end,
						NewText: text,
					},
				},
			},
		},
		Related: nil,
	})
}
