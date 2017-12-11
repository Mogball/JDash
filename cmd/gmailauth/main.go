package main

import (
	"fmt"
	"log"
	"google.golang.org/api/gmail/v1"
	"jdash/api"
	"strings"
	"regexp"
	"strconv"
	"math"
)

func main() {
	srv, err := gmail.New(api.CreateGmailClient())
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	username := "me"
	response, err := srv.Users.Messages.List(username).Q("from:'Uber Receipts'").Do()
	if err != nil {
		log.Fatal(err)
	}
	dollarRegex := regexp.MustCompile("(\\$[0-9,]+(\\.[0-9]{2})?)")
	totalSpent := 0
	uberCount := 0
	uberSpent := 0
	uberEatsCount := 0
	uberEats := 0
	uberEatsFee := 0
	cancelled := 0
	cancelledFee := 0
	for _, message := range response.Messages {
		res, err := srv.Users.Messages.Get(username, message.Id).Do()
		if err != nil {
			log.Println(err)
		}
		match := dollarRegex.FindString(res.Snippet)
		if match != "" {
			match = strings.Replace(match[1:], ",", "", -1)
			money, err := strconv.ParseFloat(match, 32)
			if err != nil {
				continue
			}
			cents := int(math.Floor(money*100 + 0.5))
			totalSpent += cents
			if strings.Contains(res.Snippet, "UberEATS") {
				uberEatsCount++
				uberEats += cents
				uberEatsFee += 500
			} else {
				uberCount++
				uberSpent += cents
			}
		} else {
			cancelled++
			cancelledFee += 500
		}
	}
	fmt.Printf("Total: $%.2f\n", float64(totalSpent) / 100.0)
	fmt.Printf("UberEATS Count: %d\n", uberEatsCount)
	fmt.Printf("UberEATS: $%.2f\n", float64(uberEats) / 100.0)
	fmt.Printf("UberEATS Fee: $%.2f\n", float64(uberEatsFee) / 100.0)
	fmt.Printf("Uber Count: %d\n", uberCount)
	fmt.Printf("Uber: $%.2f", float64(uberSpent) / 100.0)
}
