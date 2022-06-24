package main

import (
	"go-formatter/cmd/formatter/empty_func_body"
	"go-formatter/cmd/formatter/with"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		with.Analyzer,
		empty_func_body.Analyzer,
	)
}
