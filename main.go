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
	"jdash/render"
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
	router.GET("strangetracker/dom/view", strangeDomView)

	router.GET("render/view", func(c *gin.Context) {
		files, _ := json.Marshal(render.GetRenderFiles())
		c.HTML(http.StatusOK, "render.tmpl.html", gin.H{"files": string(files)})
	})

	router.POST("code/encode", encodeStringAndSend)
	router.POST("code/decode", decodeStringAndSend)
	router.GET("code", func(c *gin.Context) {
		c.HTML(http.StatusOK, "code.tmpl.html", gin.H{
			"encodeUrl": "code/encode",
			"decodeUrl": "code/decode",
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
		"axis": "Trump Mentions",
		"metric": "MinorMatches",
	})
}

func encodeStringAndSend(c *gin.Context) {
	message := c.DefaultPostForm("message", "DEFAULT")
	encoded := strangetracker.EncodeString(message, " ", strangetracker.AppDefaultCode())
	c.String(http.StatusOK, encoded)
}

func decodeStringAndSend(c *gin.Context) {
	encoded := c.PostForm("encoded")
	message := strangetracker.CrackString(encoded, " ", strangetracker.AppDefaultCode())
	c.String(http.StatusOK, message)
}

func strangeDomView(c *gin.Context) {
	bootstrapData, err := json.Marshal(strangetracker.GetGraphData(strangetracker.DefaultLookbehind()))
	if err != nil {
		c.Error(err)
		return
	}
	c.HTML(http.StatusOK, "graph.tmpl.html", gin.H{
		"data": string(bootstrapData),
		"axis": "DOM Count",
		"metric": "default",
	})
}
