package main

import (
	"jdash/ubercounter"
	"jdash/api"
	"log"
)

func main() {
	conf, err := api.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	tok := api.GetCacheToken(conf)
	if err != nil {
		log.Fatal(err)
	}
	count, _ := ubercounter.UberCountFor("jeffniu22@gmail.com", conf, tok)
	ubercounter.PrettyPrintCount(count)
}
