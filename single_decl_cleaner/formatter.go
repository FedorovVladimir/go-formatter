package single_decl_cleaner

import (
	"go/ast"
	"go/token"
	"io/ioutil"

	"github.com/FedorovVladimir/go-formatter/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const Doc = `single_decl_cleaner`

var Analyzer = &analysis.Analyzer{
	Name:     "single_decl_cleaner",
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
			group := n.(*ast.GenDecl)
			if group.Lparen == token.NoPos {
				continue
			}
			if len(group.Specs) > 1 {
				continue
			}

			spec := group.Specs[0].(*ast.ValueSpec)

			if spec.Doc != nil {
				text := utils.CutTextFromFile(fileBytes, currentFile, spec.Doc.Pos(), spec.Doc.End())
				text = append(text, []byte("\n")...)
				utils.Report(pass, group.Pos(), group.Pos(), text, "incorrect single declaration style for doc")
			}

			specEnd := spec.End()
			if spec.Comment != nil {
				specEnd = spec.Comment.End()
			}
			text := utils.CutTextFromFile(fileBytes, currentFile, spec.Pos(), specEnd)
			utils.Report(pass, group.Lparen, group.Rparen+1, text, "incorrect single declaration style")
		}
	}

	return nil, nil
}
