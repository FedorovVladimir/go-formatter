package decl_to_groups

import (
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
	specsText []byte
	groupPos  token.Pos
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
			if group.Lparen != token.NoPos {
				continue
			}
			if len(group.Specs) > 1 {
				continue
			}

			if len(groups) == 0 || groups[len(groups)-1].groupPos != 0 {
				groups = append(groups, declGroup{})
			}

			spec := group.Specs[0]

			groups[len(groups)-1].specsText = getText(fileBytes, currentFile, spec)
			groups[len(groups)-1].groupPos = group.Pos()
			utils.Report(pass, group.Pos(), getSpecEnd(spec), []byte{}, "rm decl")
		}
		for _, group := range groups {
			text := append(append([]byte(("var (\n")), group.specsText...), []byte("\n)")...)
			utils.Report(pass, group.groupPos, group.groupPos, text, "incorrect single declaration style")
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
