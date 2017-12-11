package main

import (
	"log"
	"google.golang.org/api/gmail/v1"
	"jdash/api"
	"regexp"
	"github.com/jfcote87/google-api-go-client/batch"
	"github.com/jfcote87/google-api-go-client/batch/credentials"
	"golang.org/x/net/context"
	"fmt"
	"strings"
	"strconv"
	"math"
)

func main() {
	clientConfig, err := api.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	token := api.GetToken(clientConfig)
	ctx := context.Background()
	client := clientConfig.Client(ctx, token)
	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}
	batchsrv := batch.Service{}
	batchgsv, _ := gmail.New(batch.BatchClient)

	username := "me"
	messages := make([]*gmail.Message, 0)
	response, err := srv.Users.Messages.List(username).Q("from:'Uber Receipts'").Fields("nextPageToken,messages(id)").Do()
	if err != nil {
		log.Fatal(err)
	}
	for ; response.Messages != nil; {
		messages = append(messages, response.Messages...)
		if len(response.NextPageToken) > 0 {
			response, err = srv.Users.Messages.List(username).Q("from:'Uber Receipts'").PageToken(response.NextPageToken).Do()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			break
		}
	}
	dollarRegex := regexp.MustCompile("(\\$[0-9,]+(\\.[0-9]{2})?)")
	dollarRegex.FindString("sd")
	for _, message := range messages {
		res, err := batchgsv.Users.Messages.Get(username, message.Id).Fields("snippet").Do()
		cred := &credentials.Oauth2Credentials{
			TokenSource: clientConfig.TokenSource(ctx, token),
		}
		if err = batchsrv.AddRequest(err,
			batch.SetResult(&res),
			batch.SetCredentials(cred)); err != nil {
			log.Fatal(err)
		}
	}
	responses, err := batchsrv.DoCtx(ctx)
	if err != nil {
		log.Fatal(err)
	}
	/*for _, r := range responses {
		if r.Err != nil {
			log.Fatal(r.Err)
			continue
		}
		res := r.Result.(*gmail.Message)
		fmt.Println(res.Snippet)
	}*/
	totalSpent := 0
	uberCount := 0
	uberSpent := 0
	uberEatsCount := 0
	uberEats := 0
	uberEatsFee := 0
	cancelled := 0
	cancelledFee := 0
	for _, r := range responses {
		res := r.Result.(*gmail.Message)
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
	fmt.Printf("Total: $%.2f\n", float64(totalSpent)/100.0)
	fmt.Printf("UberEATS Count: %d\n", uberEatsCount)
	fmt.Printf("UberEATS: $%.2f\n", float64(uberEats)/100.0)
	fmt.Printf("UberEATS Fee: $%.2f\n", float64(uberEatsFee)/100.0)
	fmt.Printf("Uber Count: %d\n", uberCount)
	fmt.Printf("Uber: $%.2f", float64(uberSpent)/100.0)
}
