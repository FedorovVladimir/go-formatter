package main

import (
	"log"

	"github.com/FedorovVladimir/go-formatter/config"
	"github.com/FedorovVladimir/go-formatter/decl_to_groups"
	"github.com/FedorovVladimir/go-formatter/formatter_order"
	"github.com/FedorovVladimir/go-formatter/single_decl_cleaner"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

var formatters = map[string]*analysis.Analyzer{
	"formatter_order":     formatter_order.Analyzer,
	"single_decl_cleaner": single_decl_cleaner.Analyzer,
	"decl_to_groups":      decl_to_groups.Analyzer,
}

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	var enabledFormatters []*analysis.Analyzer
	for _, formatter := range cfg.Formatters {
		if formatter.Enabled {
			enabledFormatters = append(enabledFormatters, formatters[formatter.Name])
		}
	}
	multichecker.Main(enabledFormatters...)
}
