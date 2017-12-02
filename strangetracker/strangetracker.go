package strangetracker

import (
	"jdash/config"
	"jdash/app"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"errors"
	"log"
	"regexp"
	"strconv"
	"time"
)

type DOMResult struct {
	Time  int64 `json:"time",firebase:"time"`
	Count int   `json:"count",firebase:"count"`
}

func buildCode(conf *config.Config, limit int) *CodeSet {
	return &CodeSet{
		A:     int64(conf.Number[config.STRANGE_DOM_OFFSET_A]),
		B:     int64(conf.Number[config.STRANGE_DOM_OFFSET_B]),
		C:     int64(conf.Number[config.STRANGE_DOM_OFFSET_C]),
		Mod:   int64(conf.Number[config.STRANGE_DOM_MOD]),
		Limit: int64(limit),
	}
}

func breakAndGetDOMFromConfig() string {
	code := buildCode(app.Config, config.NUM_CHARS)
	sep := app.Config.Word[config.STRANGE_DOM_SEP]
	value := app.Config.Word[config.STRANGE_DOM_STRING]
	return CrackString(value, sep, code)
}

func stripRawCountFromTarget(url string, countRegex string) (int, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return 0, err
	}
	selection := doc.Find("span.count")
	if selection.Length() == 0 {
		return 0, errors.New("failed to find count element")
	}
	if selection.Length() > 1 {
		log.Println("[WARNING] Found more than one count element")
	}
	rawCount := selection.First().Text()
	countMatcher := regexp.MustCompile(countRegex)
	if countMatcher.FindStringIndex(rawCount) == nil {
		return 0, errors.New("count does not match regex")
	}
	log.Printf("Raw count value is [%s]", rawCount)
	commaStripped := strings.Replace(rawCount, ",", "", -1)
	count, err := strconv.Atoi(commaStripped[1:len(commaStripped)-1])
	if err != nil {
		return 0, errors.New("count was not a number")
	}
	return count, nil
}

func getCountFromDOMConfig() int {
	regex := app.Config.Word[config.STRANGE_DOM_RAW_COUNT]
	url := breakAndGetDOMFromConfig()
	count, err := stripRawCountFromTarget(url, regex)
	if err != nil {
		log.Println(err)
		count = 0
	}
	return count
}

func TrackDOMNow() DOMResult {
	timeSeconds := time.Now().Unix()
	log.Printf("Generating DOM result for time [%d]", timeSeconds)
	count := getCountFromDOMConfig()
	return DOMResult{
		Time:  timeSeconds,
		Count: count,
	}
}
