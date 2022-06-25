package main

import (
	"fmt"

	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/config"
	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/context_first_parameter"
	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/empty_func_body"
	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/grouped_vars"
	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/many_arguments"
	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/methods_with_star_and_rename"
	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/new_line"
	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/rm_ignore_vars"
	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/start_enums_at_one"
	"git.user-penguin.space/vladimir/go-formatter/cmd/formatter/with"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

var m = map[string]*analysis.Analyzer{
	"context_first_parameter":      context_first_parameter.Analyzer,
	"empty_func_body":              empty_func_body.Analyzer,
	"grouped_vars":                 grouped_vars.Analyzer,
	"many_arguments":               many_arguments.Analyzer,
	"methods_with_star_and_rename": methods_with_star_and_rename.Analyzer,
	"new_line":                     new_line.Analyzer,
	"rm_ignore_vars":               rm_ignore_vars.Analyzer,
	"start_enums_at_one":           start_enums_at_one.Analyzer,
	"with":                         with.Analyzer,
}

func main() {
	c, err := config.ReadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	var formatters []*analysis.Analyzer
	for _, formatter := range c.Formatters {
		if formatter.On {
			formatters = append(formatters, m[formatter.Name])
		}
	}
	multichecker.Main(formatters...)
}
