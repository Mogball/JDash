package cron

import (
    "strconv"
    "net/url"
    "cloud.google.com/go/firestore"

    "jdash/app"
    "jdash/app/config"
    "jdash/app/trumptracker"
    "log"
)

func TrumpTrackerTask() {
    trackerResultList, timeSeconds := trumptracker.TrumpTrackNow()
    timeKey := strconv.FormatInt(timeSeconds, 10)
    hourlyData := app.FirestoreClient.Collection(config.FIRESTORE_TRUMP_DATA).Doc(config.HOURLY)
    trackerData := hourlyData.Collection(config.DATA)
    trackerMap := make(map[string]*trumptracker.TrumpTrackResult)
    for i := 0; i < len(trackerResultList); i++ {
        trackedUrl, err := url.Parse(trackerResultList[i].Url)
        if err != nil {
            log.Fatalf("Invalid tracker URL [%s]", trackerResultList[i].Url)
            continue
        }
        trackerMap[trackedUrl.Host] = trackerResultList[i]
    }
    _, err := trackerData.Doc(timeKey).Set(app.Context, trackerMap, firestore.MergeAll)
    if err != nil {
        log.Fatalln(err)
    } else {
        log.Printf("Pushed [%d] track results to Firestore", len(trackerResultList))
    }
}
