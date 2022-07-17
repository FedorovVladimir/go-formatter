package formatter_order

import (
	"go/ast"
	"go/token"
	"io"
	"os"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `formatter_order`

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

type decl int

type position struct {
	pos      token.Pos
	end      token.Pos
	filename string
}

type fileData struct {
	groups       map[decl][]*position
	lastNode     *position
	positions    []*position
	lastPosition *position
}

func newFileData() *fileData {
	return &fileData{
		groups:       map[decl][]*position{},
		lastNode:     nil,
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
				if data[currentFile].lastNode != nil {
					data[currentFile].lastNode.end = e.Pos() - 1
				}
				data[currentFile].lastNode = &position{pos: e.Pos(), end: e.End(), filename: currentFile.Name()}
				data[currentFile].groups[funcDecl] = append(data[currentFile].groups[funcDecl], data[currentFile].lastNode)

				if data[currentFile].lastPosition != nil {
					data[currentFile].lastPosition.end = e.Pos() - 1
				}
				data[currentFile].lastPosition = &position{pos: e.Pos(), end: e.End()}
				data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
			case *ast.GenDecl:
				switch e.Tok {
				case token.CONST:
					if data[currentFile].lastNode != nil {
						data[currentFile].lastNode.end = e.Pos() - 1
					}
					data[currentFile].lastNode = &position{pos: e.Pos(), end: e.End(), filename: currentFile.Name()}
					data[currentFile].groups[constDecl] = append(data[currentFile].groups[constDecl], data[currentFile].lastNode)

					if data[currentFile].lastPosition != nil {
						data[currentFile].lastPosition.end = e.Pos() - 1
					}
					data[currentFile].lastPosition = &position{pos: e.Pos(), end: e.End()}
					data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
				case token.VAR:
					if currentFile.Position(e.Pos()).Column != 1 {
						return
					}

					if data[currentFile].lastNode != nil {
						data[currentFile].lastNode.end = e.Pos() - 1
					}
					data[currentFile].lastNode = &position{pos: e.Pos(), end: e.End(), filename: currentFile.Name()}
					data[currentFile].groups[varDecl] = append(data[currentFile].groups[varDecl], data[currentFile].lastNode)

					if data[currentFile].lastPosition != nil {
						data[currentFile].lastPosition.end = e.Pos() - 1
					}
					data[currentFile].lastPosition = &position{pos: e.Pos(), end: e.End()}
					data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
				case token.TYPE:
					if data[currentFile].lastNode != nil {
						data[currentFile].lastNode.end = e.Pos() - 1
					}
					data[currentFile].lastNode = &position{pos: e.Pos(), end: e.End(), filename: currentFile.Name()}
					data[currentFile].groups[typeDecl] = append(data[currentFile].groups[typeDecl], data[currentFile].lastNode)

					if data[currentFile].lastPosition != nil {
						data[currentFile].lastPosition.end = e.Pos() - 1
					}
					data[currentFile].lastPosition = &position{pos: e.Pos(), end: e.End()}
					data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
				case token.FUNC:
					if data[currentFile].lastNode != nil {
						data[currentFile].lastNode.end = e.Pos() - 1
					}
					data[currentFile].lastNode = &position{pos: e.Pos(), end: e.End(), filename: currentFile.Name()}
					data[currentFile].groups[funcDecl] = append(data[currentFile].groups[funcDecl], data[currentFile].lastNode)

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
		f := pass.Fset.File(data[file].lastPosition.pos)
		end := token.Pos(f.Base() + f.Size())
		data[file].lastPosition.end = end - 1
		data[file].lastNode.end = end

		i := 0
		for _, decl := range orderDecl {
			if ok, k := reportGroup(pass, data[file].positions, i, data[file].groups, decl, f); ok {
				i = k
			}
		}
	}

	return nil, nil
}

func reportGroup(pass *analysis.Pass, positions []*position, i int, groups map[decl][]*position, decl decl, f *token.File) (bool, int) {
	if nodes, ok := groups[decl]; ok {
		for _, node := range nodes {
			if node.pos == positions[i].pos {
				i++
				continue
			}

			node.pos = token.Pos(int(node.pos) - f.Base())
			node.end = token.Pos(int(node.end) - f.Base())
			d, _ := readFile(node.filename)
			text := d[node.pos:node.end]

			report(pass, positions[i].pos, positions[i].end, text, "formatter_order")
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

func readFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]byte, 64)
	str := ""
	for {
		fl, err := file.Read(data)
		if err == io.EOF {
			break
		}
		str += string(data[:fl])
	}

	return []byte(str), nil
}
