package cron

import (
	"strconv"
	"net/url"
	"cloud.google.com/go/firestore"

	"jdash/app"
	"jdash/config"
	"jdash/trumptracker"
	"log"
	"jdash/strangetracker"
)

func TrumpTrackerTask() {
	trackerResultList, timeSeconds := trumptracker.TrumpTrackNow()
	timeKey := strconv.FormatInt(timeSeconds, 10)
	hourlyData := app.FirestoreClient.Collection(config.FIRESTORE_TRUMP_DATA).Doc(config.HOURLY).Collection(config.DATA)
	trackerMap := make(map[string]*trumptracker.TrumpTrackResult)
	for i := 0; i < len(trackerResultList); i++ {
		trackedUrl, err := url.Parse(trackerResultList[i].Url)
		if err != nil {
			log.Printf("[WARNING] Invalid tracker URL [%s]", trackerResultList[i].Url)
			continue
		}
		trackerMap[trackedUrl.Host] = trackerResultList[i]
	}
	dataDocument := hourlyData.Doc(timeKey)
	_, err := dataDocument.Set(app.Context, trackerMap, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	_, err = dataDocument.Set(app.Context, map[string]int64 {config.TIME: timeSeconds}, firestore.MergeAll)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Pushed [%d] track results to Firestore", len(trackerResultList))
}

func StrangeTrackerDOMTask() {
	domResult := strangetracker.TrackDOMNow()
	timeKey := strconv.FormatInt(domResult.Time, 10)
	dailyData := app.FirestoreClient.Collection(config.FIRESTORE_STRANGE_TRACKER).Doc(config.FIRESTORE_DOM_DATA).Collection(config.DATA)
	dailyData.Doc(timeKey).Set(app.Context, domResult)
	log.Printf("Pushed DOM result to Firestore")
}
