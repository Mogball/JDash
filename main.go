package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    _ "github.com/heroku/x/hmetrics/onload"
    "jdash/app"
    "jdash/trumptracker"
    "jdash/cron"
)

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        log.Fatal("$PORT must be set")
    }

    router := gin.New()
    router.Use(gin.Logger())
    router.LoadHTMLGlob("templates/*.tmpl.html")
    router.Static("/static", "static")
    router.StaticFile("/favicon.ico", "resources/favicon.ico")

    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.tmpl.html", nil)
    })
    router.GET("/trumptracker/sites", func(c *gin.Context) {
        c.JSON(http.StatusOK, trumptracker.GetTrackedSites())
    })
    router.GET("/trumptracker/now/get", func(c *gin.Context) {
        trackerResultList, _ := trumptracker.TrumpTrackNow()
        c.JSON(http.StatusOK, trackerResultList)
    })
    router.GET("trumptracker/now/push", func(c *gin.Context) {
        cron.TrumpTrackerTask()
        c.String(http.StatusOK, "OK")
    })

    app.Init()
    cron.ScheduleTasks()

    router.Run(":" + port)
}
