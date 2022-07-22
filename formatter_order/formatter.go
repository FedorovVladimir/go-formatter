package formatter_order

import (
	"go/ast"
	"go/token"
	"io"
	"os"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
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
	privateConstructorFuncDecl,
	publicFuncDecl,
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
	for _, file := range pass.Files {
		data := newFileData()

		comments := getComments(pass, file)

		for _, n := range file.Decls {
			currentFile := pass.Fset.File(n.Pos())

			var prevEnd token.Pos
			if data.lastNode != nil {
				prevEnd = data.lastNode.end
			}
			if prevEnd == token.NoPos {
				prevEnd = token.Pos(currentFile.Base())
			}
			pos := getPos(n, comments, prevEnd)
			if data.lastNode != nil {
				data.lastNode.end = pos - 1
			}
			if data.lastPosition != nil {
				data.lastPosition.end = pos - 1
			}

			switch e := n.(type) {
			case *ast.FuncDecl:
				data.lastNode = &position{pos: pos, end: e.End(), filename: currentFile.Name()}
				d := selectDeclForFunc(e.Name)
				data.groups[d] = append(data.groups[d], data.lastNode)

				data.lastPosition = &position{pos: pos, end: e.End()}
				data.positions = append(data.positions, data.lastPosition)
			case *ast.GenDecl:
				switch e.Tok {
				case token.IMPORT:
					data.lastNode = &position{pos: pos, end: e.End(), filename: currentFile.Name()}
					data.lastPosition = &position{pos: pos, end: e.End()}
				case token.CONST, token.VAR, token.TYPE:
					if currentFile.Position(e.Pos()).Column != 1 {
						continue
					}

					data.lastNode = &position{pos: pos, end: e.End(), filename: currentFile.Name()}
					data.groups[tokenToDecl[e.Tok]] = append(data.groups[tokenToDecl[e.Tok]], data.lastNode)

					data.lastPosition = &position{pos: pos, end: e.End()}
					data.positions = append(data.positions, data.lastPosition)
				}
			}
		}

		if len(data.positions) == 0 {
			continue
		}

		f := pass.Fset.File(data.lastPosition.pos)
		end := token.Pos(f.Base() + f.Size())
		data.lastPosition.end = end - 1
		data.lastNode.end = end

		i := 0
		for _, decl := range orderDecl {
			k, err := reportGroup(pass, data.positions, i, data.groups, decl, f)
			if err != nil {
				return nil, err
			}
			i = k
		}
	}

	return nil, nil
}

func getComments(pass *analysis.Pass, file *ast.File) []*position {
	comments := make([]*position, 0, len(file.Comments))
	for _, c := range file.Comments {
		if pass.Fset.Position(c.Pos()).Column != 1 {
			continue
		}
		comments = append(comments, &position{pos: c.Pos(), end: c.End()})
	}
	return comments
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

func reportGroup(pass *analysis.Pass, positions []*position, i int, groups map[decl][]*position, decl decl, f *token.File) (int, error) {
	if nodes, ok := groups[decl]; ok {
		for _, node := range nodes {
			if node.pos == positions[i].pos {
				i++
				continue
			}

			node.pos = token.Pos(int(node.pos) - f.Base())
			node.end = token.Pos(int(node.end) - f.Base())
			fileBytes, err := readFile(node.filename)
			if err != nil {
				return i, err
			}
			text := fileBytes[node.pos:node.end]

			report(pass, positions[i].pos, positions[i].end, text, "formatter_order")
			i++
		}
		return i, nil
	}
	return i, nil
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

	data := make([]byte, 64)
	str := ""
	for {
		fl, err := file.Read(data)
		if err == io.EOF {
			break
		}
		str += string(data[:fl])
	}

	return []byte(str), file.Close()
}
