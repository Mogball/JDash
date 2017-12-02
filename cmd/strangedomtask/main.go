package main

import (
	"jdash/app"
	"jdash/cron"
)

func main() {
	app.Init()
	cron.StrangeTrackerDOMTask()
}
