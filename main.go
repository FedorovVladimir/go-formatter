package main

import (
	"log"

	"github.com/FedorovVladimir/go-formatter/config"
	"github.com/FedorovVladimir/go-formatter/formatter_order"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

var m = map[string]*analysis.Analyzer{
	"formatter_order": formatter_order.Analyzer,
}

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	var enabledFormatters []*analysis.Analyzer
	for _, formatter := range cfg.Formatters {
		if formatter.Enabled {
			enabledFormatters = append(enabledFormatters, m[formatter.Name])
		}
	}
	multichecker.Main(enabledFormatters...)
}
