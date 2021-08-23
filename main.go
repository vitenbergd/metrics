package main

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"vitenberg/metrics"
)

type Specification struct {
	HttpAddr    string `default:":8080" split_words:"true" desc:"HTTP <addr>:<port> to bind "`
	HandlerPath string `default:"/metrics" split_words:"true" desc:"HTTP GET handler path"`
	TokenFile   string `split_words:"true" required:"true" desc:"Path to CSV token file"`
}

func main() {
	var s Specification
	err := envconfig.Process("metrics", &s)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Fatalf("Couldn't start metrics exporter: %s", metrics.Metrics(s.HttpAddr, s.HandlerPath, s.TokenFile))
}
