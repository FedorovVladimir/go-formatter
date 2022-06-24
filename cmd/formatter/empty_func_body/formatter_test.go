package empty_func_body_test

import (
	"go-formatter/cmd/formatter/empty_func_body"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test_Good(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, empty_func_body.Analyzer, "a")
}

func Test_Bad(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, empty_func_body.Analyzer, "b")
}
