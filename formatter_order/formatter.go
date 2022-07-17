package formatter_order

import (
	"bytes"
	"fmt"
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

type fileData struct {
	groups       map[decl][]ast.Node
	positions    []*position
	lastPosition *position
}

func newFileData() *fileData {
	return &fileData{
		groups:       map[decl][]ast.Node{},
		positions:    []*position{},
		lastPosition: nil,
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	data := map[*token.File]*fileData{}

	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(
		[]ast.Node{(*ast.GenDecl)(nil), (*ast.FuncDecl)(nil)},
		func(n ast.Node) {
			currentFile := pass.Fset.File(n.Pos())
			if _, ok := data[currentFile]; !ok {
				data[currentFile] = newFileData()
			}

			switch e := n.(type) {
			case *ast.FuncDecl:
				data[currentFile].groups[funcDecl] = append(data[currentFile].groups[funcDecl], e)
				if data[currentFile].lastPosition != nil {
					data[currentFile].lastPosition.end = e.Pos() - 1
				}
				data[currentFile].lastPosition = &position{pos: e.Pos(), end: e.End()}
				data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
			case *ast.GenDecl:
				switch e.Tok {
				case token.CONST:
					data[currentFile].groups[constDecl] = append(data[currentFile].groups[constDecl], e)
					if data[currentFile].lastPosition != nil {
						data[currentFile].lastPosition.end = e.Pos() - 1
					}
					data[currentFile].lastPosition = &position{pos: e.Pos(), end: e.End()}
					data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
				case token.VAR:
					if currentFile.Position(e.Pos()).Column != 1 {
						return
					}
					data[currentFile].groups[varDecl] = append(data[currentFile].groups[varDecl], e)
					if data[currentFile].lastPosition != nil {
						data[currentFile].lastPosition.end = e.Pos() - 1
					}
					data[currentFile].lastPosition = &position{pos: e.Pos(), end: e.End()}
					data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
				case token.TYPE:
					data[currentFile].groups[typeDecl] = append(data[currentFile].groups[typeDecl], e)
					if data[currentFile].lastPosition != nil {
						data[currentFile].lastPosition.end = e.Pos() - 1
					}
					data[currentFile].lastPosition = &position{pos: e.Pos(), end: e.End()}
					data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
				case token.FUNC:
					data[currentFile].groups[funcDecl] = append(data[currentFile].groups[funcDecl], e)
					if data[currentFile].lastPosition != nil {
						data[currentFile].lastPosition.end = e.Pos() - 1
					}
					data[currentFile].lastPosition = &position{pos: e.Pos(), end: e.End()}
					data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
				}
			}
		},
	)

	for file := range data {
		f := pass.Fset.File(data[file].lastPosition.end)
		end := token.Pos(f.Base() + f.Size())
		data[file].lastPosition.end = end

		i := 0
		for _, decl := range orderDecl {
			if ok, k := reportGroup(pass, data[file].positions, i, data[file].groups, decl); ok {
				i = k
			}
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
	f := pass.Fset.File(pos)
	fmt.Println("GOV", f.Position(pos).Line, f.Position(pos).Column, f.Position(end).Line, f.Position(end).Column, string(text))
}
