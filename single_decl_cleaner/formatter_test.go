package single_decl_cleaner

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Good(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "a")
}
