package cron

import (
    "github.com/jasonlvhit/gocron"
    "log"
)

func ScheduleTrumpTracker() {
    log.Print("Scheduling TrumpTrackerTask")
    gocron.Every(5).Minutes().Do(TrumpTrackerTask)
}
