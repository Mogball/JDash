package cron

import (
	"log"
	"github.com/robfig/cron"
)

var Scheduler *cron.Cron

func ScheduleTasks() {
	Scheduler = cron.New()
	ScheduleTrumpTracker()
	ScheduleDOMTracker()
	log.Println("Starting Cron Scheduler")
	Scheduler.Start()
}

func ScheduleTrumpTracker() {
	log.Println("Scheduling TrumpTrackerTask for [HOURLY]")
	Scheduler.AddFunc("@hourly", TrumpTrackerTask)
}

func ScheduleDOMTracker() {
	log.Println("Scheduling StrangeTrackDOMTask for [DAILY]")
	Scheduler.AddFunc("@daily", StrangeTrackerDOMTask)
}
