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

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
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

			spec := group.Specs[0]
			specEnd := spec.End()
			var comment []byte

			switch s := spec.(type) {
			case *ast.ValueSpec: // const and var
				if s.Comment != nil {
					specEnd = s.Comment.End()
					comment = utils.CutTextFromFile(fileBytes, currentFile, s.Comment.Pos(), s.Comment.End())
				}
			case *ast.TypeSpec:
				if s.Comment != nil {
					specEnd = s.Comment.End()
					comment = utils.CutTextFromFile(fileBytes, currentFile, s.Comment.Pos(), s.Comment.End())
				}
			case *ast.ImportSpec:
				if s.Comment != nil {
					specEnd = s.Comment.End()
					comment = utils.CutTextFromFile(fileBytes, currentFile, s.Comment.Pos(), s.Comment.End())
				}
			}

			text := utils.CutTextFromFile(fileBytes, currentFile, spec.Pos(), spec.End())
			text = append([]byte("(\n"), text...)
			text = append(text, comment...)
			text = append(text, []byte("\n)")...)
			utils.Report(pass, spec.Pos(), specEnd, text, "incorrect single declaration style")
		}
	}
	return nil, nil
}
