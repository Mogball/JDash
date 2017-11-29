package cron

import (
	"log"
	"github.com/robfig/cron"
)

var Scheduler *cron.Cron

func ScheduleTasks() {
	Scheduler = cron.New()
	ScheduleTrumpTracker()
	log.Println("Starting Cron Scheduler")
	Scheduler.Start()
}

func ScheduleTrumpTracker() {
	log.Println("Scheduling TrumpTrackerTask")
	Scheduler.AddFunc("@hourly", TrumpTrackerTask)
}
