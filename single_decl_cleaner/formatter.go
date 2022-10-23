package single_decl_cleaner

import (
	"go/ast"
	"go/token"
	"io/ioutil"

	"github.com/FedorovVladimir/go-formatter/utils"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

const Doc = `formatter_order`

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
				pos := utils.GetPosInFile(currentFile, spec.Doc.Pos())
				end := utils.GetPosInFile(currentFile, spec.Doc.End())
				text := append(fileBytes[pos:end], []byte("\n")...)
				utils.Report(pass, group.Pos(), group.Pos(), text, "incorrect single declaration style for doc")
			}

			specEnd := spec.End()
			if spec.Comment != nil {
				specEnd = spec.Comment.End()
			}
			pos := utils.GetPosInFile(currentFile, spec.Pos())
			end := utils.GetPosInFile(currentFile, specEnd)
			text := fileBytes[pos:end]
			utils.Report(pass, group.Lparen, group.Rparen+1, text, "incorrect single declaration style")
		}
	}

	return nil, nil
}
