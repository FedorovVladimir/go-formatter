package single_decl_cleaner

import (
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
	return nil, nil
}
