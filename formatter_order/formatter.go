package formatter_order

import (
	"go/ast"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/FedorovVladimir/go-formatter/utils"
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
	currentFile  *token.File
	comments     []*position
}

func newFileData(currentFile *token.File) *fileData {
	return &fileData{
		groups:       map[decl][]*position{},
		lastNode:     nil,
		positions:    []*position{},
		lastPosition: nil,
		currentFile:  currentFile,
	}
}

func (data *fileData) getComments(pass *analysis.Pass, file *ast.File) {
	data.comments = make([]*position, 0, len(file.Comments))
	for _, c := range file.Comments {
		if pass.Fset.Position(c.Pos()).Column != 1 {
			continue
		}
		data.comments = append(data.comments, &position{pos: c.Pos(), end: c.End()})
	}
}

func (data *fileData) getPrevEnd() token.Pos {
	if data.lastNode != nil {
		return data.lastNode.end
	}
	return token.Pos(data.currentFile.Base())
}

func (data *fileData) getPos(n ast.Node, end token.Pos) token.Pos {
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
	for _, comment := range data.comments {
		if end < comment.pos && comment.pos < pos {
			pos = comment.pos
		}
	}
	return pos
}

func (data *fileData) appendFuncDecl(pos token.Pos, e *ast.FuncDecl) {
	data.lastNode = &position{pos: pos, end: e.End(), filename: data.currentFile.Name()}
	d := selectDeclForFunc(e.Name)
	data.groups[d] = append(data.groups[d], data.lastNode)

	data.lastPosition = &position{pos: pos, end: e.End()}
	data.positions = append(data.positions, data.lastPosition)
}

func (data *fileData) appendGenDecl(pos token.Pos, e *ast.GenDecl) {
	switch e.Tok {
	case token.IMPORT:
		data.lastNode = &position{pos: pos, end: e.End(), filename: data.currentFile.Name()}
		data.lastPosition = &position{pos: pos, end: e.End()}
	case token.CONST, token.VAR, token.TYPE:
		if data.currentFile.Position(e.Pos()).Column != 1 {
			return
		}

		data.lastNode = &position{pos: pos, end: e.End(), filename: data.currentFile.Name()}
		data.groups[tokenToDecl[e.Tok]] = append(data.groups[tokenToDecl[e.Tok]], data.lastNode)

		data.lastPosition = &position{pos: pos, end: e.End()}
		data.positions = append(data.positions, data.lastPosition)
	}
}

func (data *fileData) reportGroup(pass *analysis.Pass, i int, decl decl) (int, error) {
	nodes, ok := data.groups[decl]
	if !ok {
		return i, nil
	}
	for _, node := range nodes {
		if node.pos == data.positions[i].pos {
			i++
			continue
		}

		node.pos = utils.GetPosInFile(data.currentFile, node.pos)
		node.end = utils.GetPosInFile(data.currentFile, node.end)
		fileBytes, err := ioutil.ReadFile(node.filename)
		if err != nil {
			return i, err
		}
		text := fileBytes[node.pos:node.end]

		utils.Report(pass, data.positions[i].pos, data.positions[i].end, text, "incorrect declaration order")
		i++
	}
	return i, nil
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if len(file.Decls) == 0 {
			continue
		}

		data := newFileData(pass.Fset.File(file.Decls[0].Pos()))
		data.getComments(pass, file)

		for _, n := range file.Decls {
			prevEnd := data.getPrevEnd()
			pos := data.getPos(n, prevEnd)

			if data.lastNode != nil {
				data.lastNode.end = pos - 1
			}
			if data.lastPosition != nil {
				data.lastPosition.end = pos - 1
			}

			switch e := n.(type) {
			case *ast.FuncDecl:
				data.appendFuncDecl(pos, e)
			case *ast.GenDecl:
				data.appendGenDecl(pos, e)
			}
		}

		end := token.Pos(data.currentFile.Base() + data.currentFile.Size())
		data.lastPosition.end = end - 1
		data.lastNode.end = end

		i := 0
		for _, decl := range orderDecl {
			k, err := data.reportGroup(pass, i, decl)
			if err != nil {
				return nil, err
			}
			i = k
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
