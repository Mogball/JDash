package main

import (
	"jdash/app"
	"jdash/config"
	"time"
	"google.golang.org/api/iterator"
	"log"
	"fmt"
)

func main() {
	app.Init()
	lookbehindSeconds := time.Now().Unix() - int64(app.Config.Number[config.FIRESTORE_TRUMP_LOOKBACK] * config.SEC_IN_HRS)
	hourlyData := app.FirestoreClient.Collection(config.FIRESTORE_TRUMP_DATA).Doc(config.HOURLY).Collection(config.DATA)
	dataIter := hourlyData.Where("time", ">=", lookbehindSeconds).Documents(app.Context)
	for {
		doc, err := dataIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalln(err)
			break
		}
		fmt.Println(doc.Data())
	}
}
