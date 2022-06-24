package with_test

import (
	"learn_go/cmd/formatter/with"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Good(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, with.Analyzer, "a")
}

func Test_Bad(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, with.Analyzer, "b")
}
