package grouped_vars

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Good(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "a")
}

func Test_Bad_1(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "b")
}

func Test_Bad_2(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "c")
}
