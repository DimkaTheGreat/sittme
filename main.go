package main

import (
	"flag"

	"github.com/DimkaTheGreat/sittme/models"
	"github.com/DimkaTheGreat/sittme/routing"
)

var (
	timeout = flag.Int("timeout", 20, "timeout between interrupted and finished state")
	port    = flag.String("port", "8086", "server port")
)

func main() {
	flag.Parse()

	translations := models.Translations{}

	translations.LoadTestData()

	routing.Run(translations, *timeout, *port)

}
