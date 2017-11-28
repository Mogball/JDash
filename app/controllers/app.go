package controllers

import (
    "github.com/revel/revel"

    "jdash/app/cron"
    "jdash/app/trumptracker"
)

type App struct {
    *revel.Controller
}

func (c App) Index() revel.Result {
    return c.Render()
}

func (c App) TrumpSites() revel.Result {
    return c.RenderJSON(trumptracker.GetTrackedSites())
}

func (c App) TrumpTrack() revel.Result {
    trackerResultList, _ := trumptracker.TrumpTrackNow()
    return c.RenderJSON(trackerResultList)
}

func (c App) TrumpTrackPushResults() revel.Result {
    cron.TrumpTrackerTask()
    return c.RenderHTML("OK")
}
