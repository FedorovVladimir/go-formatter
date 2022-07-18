package formatter_order

import (
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `formatter_order`

const (
	constDecl decl = iota + 1
	varDecl
	typeDecl
	publicConstructorFuncDecl
	publicFuncDecl
	privateConstructorFuncDecl
	privateFuncDecl
)

var orderDecl = []decl{
	constDecl,
	varDecl,
	typeDecl,
	publicConstructorFuncDecl,
	publicFuncDecl,
	privateConstructorFuncDecl,
	privateFuncDecl,
}

var tokenToDecl = map[token.Token]decl{
	token.CONST: constDecl,
	token.VAR:   varDecl,
	token.TYPE:  typeDecl,
}

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

	var comments []*position
	for _, file := range pass.Files {
		for _, c := range file.Comments {
			if pass.Fset.Position(c.Pos()).Column != 1 {
				continue
			}
			comments = append(comments, &position{pos: c.Pos(), end: c.End()})
		}
	}

	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(
		[]ast.Node{(*ast.GenDecl)(nil), (*ast.FuncDecl)(nil)},
		func(n ast.Node) {
			currentFile := pass.Fset.File(n.Pos())
			if _, ok := data[currentFile]; !ok {
				data[currentFile] = newFileData()
			}

			var end token.Pos
			if data[currentFile].lastNode != nil {
				end = data[currentFile].lastNode.end
			}
			if end == token.NoPos {
				end = token.Pos(currentFile.Base())
			}
			pos := getPos(n, comments, end)
			if data[currentFile].lastNode != nil {
				data[currentFile].lastNode.end = pos - 1
			}
			if data[currentFile].lastPosition != nil {
				data[currentFile].lastPosition.end = pos - 1
			}

			switch e := n.(type) {
			case *ast.FuncDecl:
				data[currentFile].lastNode = &position{pos: pos, end: e.End(), filename: currentFile.Name()}
				d := selectDeclForFunc(e.Name)
				data[currentFile].groups[d] = append(data[currentFile].groups[d], data[currentFile].lastNode)

				data[currentFile].lastPosition = &position{pos: pos, end: e.End()}
				data[currentFile].positions = append(data[currentFile].positions, data[currentFile].lastPosition)
			case *ast.GenDecl:
				switch e.Tok {
				case token.CONST, token.VAR, token.TYPE:
					if currentFile.Position(e.Pos()).Column != 1 {
						return
					}

					data[currentFile].lastNode = &position{pos: pos, end: e.End(), filename: currentFile.Name()}
					data[currentFile].groups[tokenToDecl[e.Tok]] = append(data[currentFile].groups[tokenToDecl[e.Tok]], data[currentFile].lastNode)

					data[currentFile].lastPosition = &position{pos: pos, end: e.End()}
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

func selectDeclForFunc(name *ast.Ident) decl {
	n := name.Name
	if strings.HasPrefix(n, "New") {
		return publicConstructorFuncDecl
	}
	if strings.HasPrefix(n, "new") {
		return privateConstructorFuncDecl
	}
	if name.IsExported() {
		return publicFuncDecl
	}
	return privateFuncDecl
}

func getPos(n ast.Node, comments []*position, end token.Pos) token.Pos {
	pos := n.Pos()
	switch e := n.(type) {
	case *ast.FuncDecl:
		if e.Doc != nil {
			pos = e.Doc.Pos()
		}
	case *ast.GenDecl:
		if e.Doc != nil {
			pos = e.Doc.Pos()
		}
	}
	for _, comment := range comments {
		if end < comment.pos && comment.pos < pos {
			pos = comment.pos
		}
	}
	return pos
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
			fileBytes, _ := readFile(node.filename)
			text := fileBytes[node.pos:node.end]

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
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

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
