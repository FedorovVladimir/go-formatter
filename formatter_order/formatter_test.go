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

func Test_Bad5(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "h")
}

func Test_Bad6(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "i")
}

func Test_Bad7(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "j")
}

func Test_Bad8(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "k")
}

func Test_Bad9(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "l")
}

func Test_Bad10(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "m")
}

func Test_Good4(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "n")
}
