package main

import (
	"jdash/app"
	"jdash/trumptracker"
	"google.golang.org/api/iterator"
	"log"
	"jdash/config"
	"cloud.google.com/go/firestore"
	"sort"
	"strconv"
	"golang.org/x/net/context"
)

func main() {
	app.Init()
	runTask()
}

type TimeList []int64

func (l TimeList) Len() int {
	return len(l)
}

func (l TimeList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l TimeList) Less(i, j int) bool {
	return l[i] < l[j]
}

func runTask() {
	// Scan last 36 hours of data and check for anomalies
	log.Println("Maintain task checking 54 hours of data")
	dataIt := trumptracker.GetDataIteratorSince(trumptracker.LookbehindFor(54))
	dataMap, times := parseDataIterator(dataIt)
	ctx := context.Background()
	sort.Sort(times)
	for i, time := range times[1:] {
		for name, value := range dataMap[time] {
			valMap := value.(map[string]interface{})
			major := valMap["MajorMatches"].(int64)
			minor := valMap["MinorMatches"].(int64)
			if major <= 0 && minor <= 0 {
				log.Printf("Found error at time [%d] for [%s]", time, name)
				prevVal, prevOk := dataMap[times[i]][name]
				nextVal, nextOk := dataMap[times[i+2]][name]
				if !prevOk || !nextOk {
					continue
				}
				curRes := parseTrumpTrackResult(value)
				prevRes := parseTrumpTrackResult(prevVal)
				nextRes := parseTrumpTrackResult(nextVal)
				if (prevRes.MinorMatches <= 0 && prevRes.MajorMatches <= 0) || (nextRes.MinorMatches <= 0 && prevRes.MajorMatches <= 0) {
					continue
				}
				curRes.MajorMatches = (prevRes.MajorMatches + nextRes.MajorMatches) / 2
				curRes.MinorMatches = (prevRes.MinorMatches + nextRes.MinorMatches) / 2
				trackerData := app.FirestoreClient().Collection(config.FIRESTORE_TRUMP_DATA)
				hourlyData := trackerData.Doc(config.HOURLY).Collection(config.DATA)
				document := hourlyData.Doc(strconv.FormatInt(time, 10))
				_, err := document.Set(ctx, map[string]interface{}{name: curRes}, firestore.MergeAll)
				if err != nil {
					log.Printf("Error while fixing anomaly: %v", err)
				}
			}
		}
	}
}

func parseTrumpTrackResult(i interface{}) *trumptracker.TrumpTrackResult {
	valMap := i.(map[string]interface{})
	return &trumptracker.TrumpTrackResult{
		MajorMatches: int(valMap["MajorMatches"].(int64)),
		MinorMatches: int(valMap["MinorMatches"].(int64)),
		Url:          valMap["Url"].(string),
		Time:         valMap["Time"].(int64),
	}
}

func parseDataIterator(dataIt *firestore.DocumentIterator) (map[int64]map[string]interface{}, TimeList) {
	dataMap := make(map[int64]map[string]interface{}, 72)
	times := make(TimeList, 0, 72)
	for {
		doc, err := dataIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		resultMap := doc.Data()
		time := resultMap[config.TIME].(int64)
		delete(resultMap, config.TIME)
		dataMap[time] = resultMap
		times = append(times, time)
	}
	log.Println("Finished parsing data")
	return dataMap, times
}
