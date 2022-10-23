package utils

import (
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
