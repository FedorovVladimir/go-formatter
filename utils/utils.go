package utils

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
)

func Report(pass *analysis.Pass, pos token.Pos, end token.Pos, text []byte, msg string) {
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

func GetPosInFile(file *token.File, pos token.Pos) token.Pos {
	return token.Pos(int(pos) - file.Base())
}

func CutTextFromFile(fileBytes []byte, file *token.File, pos token.Pos, end token.Pos) []byte {
	posInFile := GetPosInFile(file, pos)
	endInFile := GetPosInFile(file, end)
	return fileBytes[posInFile:endInFile]
}

func GetSpecEnd(spec ast.Spec) token.Pos {
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

func GetSpecText(fileBytes []byte, currentFile *token.File, spec ast.Spec) []byte {
	return CutTextFromFile(fileBytes, currentFile, spec.Pos(), GetSpecEnd(spec))
}
