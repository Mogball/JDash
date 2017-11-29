package main

import (
    "jdash/cron"
    "jdash/app"
)

func main() {
    app.Init()
    cron.TrumpTrackerTask()
}
