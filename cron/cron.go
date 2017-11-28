package cron

import (
    "github.com/jasonlvhit/gocron"
    "log"
)

func ScheduleTasks() {
    ScheduleTrumpTracker()
    <- gocron.Start()
}

func ScheduleTrumpTracker() {
    log.Print("Scheduling TrumpTrackerTask")
    gocron.Every(1).Hour().Do(TrumpTrackerTask)
}