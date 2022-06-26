package main

import (
	"fmt"

	"github.com/FedorovVladimir/go-formatter/arguments_form"
	"github.com/FedorovVladimir/go-formatter/config"
	"github.com/FedorovVladimir/go-formatter/context_first_parameter"
	"github.com/FedorovVladimir/go-formatter/empty_func_body"
	"github.com/FedorovVladimir/go-formatter/grouped_vars"
	"github.com/FedorovVladimir/go-formatter/many_arguments"
	"github.com/FedorovVladimir/go-formatter/methods_with_star_and_rename"
	"github.com/FedorovVladimir/go-formatter/new_line"
	"github.com/FedorovVladimir/go-formatter/order"
	"github.com/FedorovVladimir/go-formatter/return_value"
	"github.com/FedorovVladimir/go-formatter/rm_ignore_vars"
	"github.com/FedorovVladimir/go-formatter/start_enums_at_one"
	"github.com/FedorovVladimir/go-formatter/with"
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
	"order":                        order.Analyzer,
	"rm_ignore_vars":               rm_ignore_vars.Analyzer,
	"start_enums_at_one":           start_enums_at_one.Analyzer,
	"with":                         with.Analyzer,
	"arguments_form":               arguments_form.Analyzer,
	"return_value":                 return_value.Analyzer,
}

func main() {
	c, err := config.ReadConfig()
	if err != nil {
		fmt.Println(err)
		return
	}

	var formatters1 []*analysis.Analyzer
	for _, formatter := range c.Formatters {
		if formatter.On {
			formatters1 = append(formatters1, m[formatter.Name])
		}
	}
	multichecker.Main(formatters1...)
}
