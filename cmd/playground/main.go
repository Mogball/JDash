package main

import (
	"jdash/render"
	"fmt"
	"jdash/ubercounter"
	"jdash/api"
	"log"
)

func main() {
	conf, err := api.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ubercounter.UberCountFor("jeffniu22@gmail.com", ))
}
