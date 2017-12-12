package ubercounter

import (
	"golang.org/x/net/context"
	"google.golang.org/api/gmail/v1"
	"github.com/jfcote87/google-api-go-client/batch"
	"jdash/config"
	"jdash/app"
	"github.com/jfcote87/google-api-go-client/batch/credentials"
	"regexp"
	"strings"
	"strconv"
	"math"
	"golang.org/x/oauth2"
	"fmt"
)

type UberCount struct {
	TotalSpent     int
	UberSpent      int
	UberEatsSpent  int
	UberCount      int
	UberEatsCount  int
	CancelledCount int
}

func PrettyPrintCount(count *UberCount) {
	fmt.Printf("Total Spent: $%.2f\n", float32(count.TotalSpent)/100.0)
	fmt.Printf("Uber Count: %d\n", count.UberCount)
	fmt.Printf("Uber Spent: $%.2f\n", float32(count.UberSpent)/100.0)
	fmt.Printf("UberEATS Count: %d\n", count.UberEatsCount)
	fmt.Printf("UberEATS Spent: $%.2f\n", float32(count.UberEatsSpent)/100.0)
}

func UberCountFor(username string, conf *oauth2.Config, token *oauth2.Token) (*UberCount, error) {
	ctx := context.Background()
	client := conf.Client(ctx, token)
	gsv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}
	batchSrv := &batch.Service{}
	batchGsv, err := gmail.New(batch.BatchClient)
	if err != nil {
		return nil, err
	}
	receipts, err := fetchUberReceipts(username, gsv, app.Config())
	if err != nil {
		return nil, err
	}
	creds := &credentials.Oauth2Credentials{TokenSource: conf.TokenSource(ctx, token)}
	for _, receipt := range receipts {
		if receipt == nil {
			continue
		}
		getCall := batchGsv.Users.Messages.Get(username, receipt.Id)
		getCall.Fields("snippet")
		res, err := getCall.Do()
		if err = batchSrv.AddRequest(err,
			batch.SetResult(&res),
			batch.SetCredentials(creds)); err != nil {
			return nil, err
		}
	}
	receiptResults, err := batchSrv.DoCtx(ctx)
	if err != nil {
		return nil, err
	}
	uberCount := &UberCount{}
	dollarRegex := regexp.MustCompile(app.Config().Word[config.UBER_COUNT_DOLLAR])
	for _, res := range receiptResults {
		receipt := res.Result.(*gmail.Message)
		match := dollarRegex.FindString(receipt.Snippet)
		if match != "" {
			match = strings.Replace(match[1:], ",", "", -1)
			money, err := strconv.ParseFloat(match, 32)
			if err != nil {
				continue
			}
			cents := int(math.Floor(money*100 + 0.5))
			uberCount.TotalSpent += cents
			if strings.Contains(receipt.Snippet, "UberEATS") {
				uberCount.UberEatsCount++
				uberCount.UberEatsSpent += cents
			} else {
				uberCount.UberCount++
				uberCount.UberSpent += cents
			}
		} else {
			uberCount.CancelledCount++
		}
	}
	return uberCount, nil
}

func fetchUberReceipts(username string, service *gmail.Service, conf *config.Config) ([]*gmail.Message, error) {
	messagesCall := service.Users.Messages.List(username)
	messagesCall.Q(conf.Word[config.UBER_COUNT_QUERY])
	messagesCall.Fields("nextPageToken,messages(id)")
	response, err := messagesCall.Do()
	if err != nil {
		return nil, err
	}
	messages := make([]*gmail.Message, 0, len(response.Messages))
	for ; response.Messages != nil; {
		messages = append(messages, response.Messages...)
		if len(response.NextPageToken) > 0 {
			messagesCall.PageToken(response.NextPageToken)
			response, err = messagesCall.Do()
			if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	return messages, nil
}
