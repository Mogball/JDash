package main

import (
	"jdash/api/gmail"
	"log"
	"jdash/app"
	"fmt"
)

func main() {
	app.Init()
	srv, err := gmail.CreateClient()
	log.Println(err)
	r, err := srv.Users.Labels.List("me").Do()
	if err != nil {
		log.Fatalln(err)
	}
	if len(r.Labels) > 0 {
		fmt.Print("Labels:\n")
		for _, l := range r.Labels {
			fmt.Printf("- %s\n", l.Name)
		}
	} else {
		fmt.Println("No labels found")
	}
}
