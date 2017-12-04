package main

import (
	"fmt"
	"log"
	"google.golang.org/api/gmail/v1"
	"jdash/api"
)

func main() {
	srv, err := gmail.New(api.CreateGmailClient())
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	username := "me"
	r, err := srv.Users.Labels.List(username).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels. %v", err)
	}
	if len(r.Labels) > 0 {
		fmt.Print("Labels:\n")
		for _, l := range r.Labels {
			fmt.Printf("- %s\n",  l.Name)
		}
	} else {
		fmt.Print("No labels found.")
	}

}