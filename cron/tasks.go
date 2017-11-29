package cron

import (
	"strconv"
	"net/url"
	"cloud.google.com/go/firestore"

	"jdash/app"
	"jdash/config"
	"jdash/trumptracker"
	"log"
)

func TrumpTrackerTask() {
	trackerResultList, timeSeconds := trumptracker.TrumpTrackNow()
	timeKey := strconv.FormatInt(timeSeconds, 10)
	hourlyData := app.FirestoreClient.Collection(config.FIRESTORE_TRUMP_DATA).Doc(config.HOURLY).Collection(config.DATA)
	trackerMap := make(map[string]*trumptracker.TrumpTrackResult)
	for i := 0; i < len(trackerResultList); i++ {
		trackedUrl, err := url.Parse(trackerResultList[i].Url)
		if err != nil {
			log.Fatalf("Invalid tracker URL [%s]", trackerResultList[i].Url)
			continue
		}
		trackerMap[trackedUrl.Host] = trackerResultList[i]
	}
	dataDocument := hourlyData.Doc(timeKey)
	_, err := dataDocument.Set(app.Context, trackerMap, firestore.MergeAll)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = dataDocument.Set(app.Context, map[string]int64 {"time": timeSeconds}, firestore.MergeAll)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Pushed [%d] track results to Firestore", len(trackerResultList))
}
