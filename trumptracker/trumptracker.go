package trumptracker

import (
	"strings"
	"time"
	"regexp"

	"net/http"
	"io/ioutil"

	"jdash/app"
	"jdash/config"
	"log"
)

type TrumpTrackResult struct {
	MajorMatches int    `json:"majorMatches",firestore:"majorMatches"`
	MinorMatches int    `json:"minorMatches",firestore:"minorMatches"`
	Url          string `json:"url",firestore:"url"`
	Time         int64  `json:"time",firestore:"time"`
}

func GetTrackedSites() []string {
	return strings.Split(app.Config.Word[config.TRUMP_SITES], ",")
}

func CountTrumps(url string) (int, int) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return -1, -1
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return -1, -1
	}
	fullMatcher := regexp.MustCompile(app.Config.Word[config.TRUMP_FULL_MATCHER])
	partMatcher := regexp.MustCompile(app.Config.Word[config.TRUMP_PART_MATCHER])
	fullMatches := fullMatcher.FindAllIndex(body, -1)
	partMatches := partMatcher.FindAllIndex(body, -1)
	return len(fullMatches), len(partMatches)
}

func GetTrumpTrackResult(url string, timeSeconds int64) *TrumpTrackResult {
	numMajor, numMinor := CountTrumps(url)
	return &TrumpTrackResult{
		MajorMatches: numMajor,
		MinorMatches: numMinor,
		Url:          url,
		Time:         timeSeconds,
	}
}

func TrumpTrackNow() ([]*TrumpTrackResult, int64) {
	trumpSites := GetTrackedSites()
	resultList := make([]*TrumpTrackResult, len(trumpSites))
	timeSeconds := time.Now().Unix()
	log.Printf("Producing TrumpTrack result at time [%d]", timeSeconds)
	for i := 0; i < len(trumpSites); i++ {
		url := trumpSites[i]
		result := GetTrumpTrackResult(url, timeSeconds)
		resultList[i] = result
	}
	return resultList, timeSeconds
}
