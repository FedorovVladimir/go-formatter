package methods_with_star_and_rename

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const Doc = `methods_with_star_and_rename`

var Analyzer = &analysis.Analyzer{
	Name:     "methods_with_star_and_rename",
	Doc:      Doc,
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

func run(pass *analysis.Pass) (interface{}, error) {
	m := map[string]string{}
	nodeFilter := []ast.Node{(*ast.FuncDecl)(nil)}
	pass.ResultOf[inspect.Analyzer].(*inspector.Inspector).Preorder(nodeFilter, func(n ast.Node) {
		e := n.(*ast.FuncDecl)
		if e.Recv == nil {
			return
		}
		oldNameNode := e.Recv.List[0].Names[0]
		t := e.Recv.List[0].Type
		recvName := ""
		needStar := false
		switch name := t.(type) {
		case *ast.StarExpr:
			ident, ok := name.X.(*ast.Ident)
			if !ok {
				return
			}
			recvName = ident.Name
		case *ast.Ident:
			needStar = true
			recvName = name.Name
		}
		_, ok := m[recvName]
		if !ok {
			m[recvName] = oldNameNode.Name
		}
		name := m[recvName]
		if needStar {
			name = name + " *"
		}
		if oldNameNode.Name != name {
			pass.Report(analysis.Diagnostic{
				Pos:      oldNameNode.Pos(),
				End:      oldNameNode.End(),
				Category: "func",
				Message:  "methods_with_star_and_rename",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "many arguments",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     oldNameNode.Pos(),
								End:     oldNameNode.End(),
								NewText: []byte(name),
							},
						},
					},
				},
				Related: nil,
			})
		}
	})
	return nil, nil
}
