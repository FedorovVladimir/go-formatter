package formatter_order

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Good1(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "a")
}

func Test_Bad1(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "b")
}

func Test_Good2(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "c")
}

func Test_Bad2(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "d")
}

func Test_Bad3(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "e")
}

func Test_Good3(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "f")
}

func Test_Bad4(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "g")
}
