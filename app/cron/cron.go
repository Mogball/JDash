package cron

import (
    "github.com/jasonlvhit/gocron"
)

func ScheduleTrumpTracker() {
    gocron.Every(5).Minutes().Do(TrumpTrackerTask)
}
