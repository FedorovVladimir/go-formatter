package main

import (
	"go-formatter/cmd/formatter/context_first_parameter"
	"go-formatter/cmd/formatter/empty_func_body"
	"go-formatter/cmd/formatter/grouped_vars"
	"go-formatter/cmd/formatter/many_arguments"
	"go-formatter/cmd/formatter/rm_ignore_vars"
	"go-formatter/cmd/formatter/start_enums_at_one"
	"go-formatter/cmd/formatter/with"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		context_first_parameter.Analyzer,
		empty_func_body.Analyzer,
		grouped_vars.Analyzer,
		many_arguments.Analyzer,
		with.Analyzer,
		rm_ignore_vars.Analyzer,
		start_enums_at_one.Analyzer,
	)
}
