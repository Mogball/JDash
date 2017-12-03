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
	"strconv"
	"encoding/json"
	"jdash/strangetracker"
)

func main() {
	port := os.Getenv("PORT")
	mode := os.Getenv("MODE")

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
	router.GET("strangetracker/dom/now/get", func(c *gin.Context) {
		domResult := strangetracker.TrackDOMNow()
		c.JSON(http.StatusOK, domResult)
	})
	router.GET("strangetracker/dom/now/push", func(c *gin.Context) {
		cron.StrangeTrackerDOMTask()
		c.String(http.StatusOK, "OK")
	})
	router.GET("trumptracker/view/data/:lookbehind", trumpTrackerViewData)
	router.GET("trumptracker/view", trumpTrackerView)

	router.GET("render/view/:file", func(c *gin.Context) {
		c.HTML(http.StatusOK, "render.tmpl.html", gin.H{
			"file": c.Param("file"),
		})
	})

	app.Init()
	if mode == "LOCAL" {
		cron.ScheduleTasks()
	}

	router.Run(":" + port)
}

func trumpTrackerViewData(c *gin.Context) {
	lookbehindStr := c.Param("lookbehind")
	var lookbehindTime int64
	if lookbehindStr == "" {
		lookbehindTime = trumptracker.DefaultLookbehind()
	} else {
		hours, err := strconv.Atoi(lookbehindStr)
		if err != nil {
			c.Error(err)
			return
		}
		lookbehindTime = trumptracker.LookbehindFor(hours)
	}
	c.JSON(http.StatusOK, trumptracker.GetGraphData(lookbehindTime))
}

func trumpTrackerView(c *gin.Context) {
	bootstrapData, err := json.Marshal(trumptracker.GetGraphData(trumptracker.DefaultLookbehind()))
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "graph.tmpl.html", gin.H{
		"data": string(bootstrapData),
	})
}
