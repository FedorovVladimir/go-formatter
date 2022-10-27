package shadowed_var

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Good(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "a")
}
