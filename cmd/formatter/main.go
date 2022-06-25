package main

import (
	"go-formatter/cmd/formatter/many_arguments"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		//with.Analyzer,
		//empty_func_body.Analyzer,
		//grouped_vars.Analyzer,
		many_arguments.Analyzer,
	)
}
