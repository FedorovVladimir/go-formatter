package decl_to_groups

import (
	"bytes"
	"go/ast"
	"go/token"
	"io/ioutil"

	"github.com/FedorovVladimir/go-formatter/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const Doc = `decl_to_groups`

var Analyzer = &analysis.Analyzer{
	Name:     "decl_to_groups",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

type declGroup struct {
	groupType string
	specsText []byte
	groupPos  token.Pos
	groupEnd  token.Pos
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		var groups []declGroup
		if len(file.Decls) == 0 {
			continue
		}

		currentFile := pass.Fset.File(file.Decls[0].Pos())
		fileBytes, err := ioutil.ReadFile(currentFile.Name())
		if err != nil {
			return nil, err
		}

		for _, n := range file.Decls {
			group, ok := n.(*ast.GenDecl)
			if !ok {
				continue
			}

			if len(groups) == 0 {
				groups = append(groups, declGroup{})
			}
			if groups[len(groups)-1].groupPos != 0 {
				lastGroupLastLine := pass.Fset.Position(groups[len(groups)-1].groupEnd).Line
				currentLine := pass.Fset.Position(group.Pos()).Line
				if currentLine-lastGroupLastLine > 1 {
					groups = append(groups, declGroup{})
				}
			}
			if groups[len(groups)-1].groupPos == 0 {
				groups[len(groups)-1].groupPos = group.Pos()
			}

			for _, spec := range group.Specs {
				text := groups[len(groups)-1].specsText
				if len(text) > 0 {
					text = append(text, []byte("\n")...)
				}
				text = append(text, []byte("\t")...)
				text = append(text, getText(fileBytes, currentFile, spec)...)
				groups[len(groups)-1].specsText = text
			}

			groupEnd := getGroupEnd(group)
			groups[len(groups)-1].groupEnd = groupEnd
			groups[len(groups)-1].groupType = getGroupType(group)
		}
		for _, group := range groups {
			oldText := utils.CutTextFromFile(fileBytes, currentFile, group.groupPos, group.groupEnd)
			text := append(append([]byte((" (\n")), group.specsText...), []byte("\n)")...)
			text = append([]byte(group.groupType), text...)
			if !bytes.Equal(oldText, text) {
				utils.Report(pass, group.groupPos, group.groupEnd, text, "incorrect single declaration style")
			}
		}
	}

	return nil, nil
}

func getText(fileBytes []byte, currentFile *token.File, spec ast.Spec) []byte {
	return utils.CutTextFromFile(fileBytes, currentFile, spec.Pos(), getSpecEnd(spec))
}

func getSpecEnd(spec ast.Spec) token.Pos {
	end := spec.End()
	switch s := spec.(type) {
	case *ast.ValueSpec: // const and var
		if s.Comment != nil {
			end = s.Comment.End()
		}
	case *ast.TypeSpec:
		if s.Comment != nil {
			end = s.Comment.End()
		}
	case *ast.ImportSpec:
		if s.Comment != nil {
			end = s.Comment.End()
		}
	default:
		panic("spec not support")
	}
	return end
}

func getGroupEnd(group *ast.GenDecl) token.Pos {
	if group.Lparen == 0 {
		return getSpecEnd(group.Specs[0])
	}
	return group.Rparen + 1
}

func getGroupType(group *ast.GenDecl) string {
	switch group.Tok {
	case token.IMPORT:
		return "import"
	case token.VAR:
		return "var"
	case token.CONST:
		return "const"
	case token.TYPE:
		return "type"
	}
	panic("tok not support")
}
