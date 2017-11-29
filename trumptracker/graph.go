package trumptracker

import (
	"jdash/app"
	"jdash/config"
	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"log"
	"time"
)

var metrics = [2]string{"MajorMatches", "MinorMatches"}

type TrackerValue struct {
	MajorMatches int64 `firestore:"majorMatches"`
	MinorMatches int64 `firestore:"minorMatches"`
	Time         int64 `firestore:"time"`
	Hostname     string
}

type Point struct {
	TimeSeconds int64 `json:"x"`
	Value       int64 `json:"y"`
}

type Series struct {
	Points []Point `json:"values"`
	Name   string  `json:"key"`
}

func getDataIteratorSince(lookbehindSeconds int64) *firestore.DocumentIterator {
	hourlyData := app.FirestoreClient.Collection(config.FIRESTORE_TRUMP_DATA).Doc(config.HOURLY).Collection(config.DATA)
	return hourlyData.Where("time", ">=", lookbehindSeconds).Documents(app.Context)
}

func resultMapToValueSlice(resultMap map[string]interface{}) []TrackerValue {
	values := make([]TrackerValue, 0, len(resultMap))
	for hostname, value := range resultMap {
		valMap := value.(map[string]interface{})
		values = append(values, TrackerValue{
			MajorMatches: valMap["MajorMatches"].(int64),
			MinorMatches: valMap["MinorMatches"].(int64),
			Time:         valMap["Time"].(int64),
			Hostname:     hostname,
		})
	}
	return values
}

func getFieldforMetric(metric string, value *TrackerValue) int64 {
	if metric == "MajorMatches" {
		return value.MajorMatches
	} else {
		return value.MinorMatches
	}
}

func parseDataIterator(dataIter *firestore.DocumentIterator) map[string][]Series {
	pointMap := make(map[string]map[string][]Point)
	for _, metric := range metrics {
		pointMap[metric] = make(map[string][]Point)

	}
	for {
		doc, err := dataIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
			break
		}
		resultMap := doc.Data()
		delete(resultMap, config.TIME)
		values := resultMapToValueSlice(resultMap)
		for _, value := range values {
			for _, metric := range metrics {
				if _, has := pointMap[metric][value.Hostname]; !has {
					pointMap[metric][value.Hostname] = make([]Point, 0)
				}
				pointMap[metric][value.Hostname] = append(pointMap[metric][value.Hostname], Point{
					TimeSeconds: value.Time,
					Value:       getFieldforMetric(metric, &value),
				})
			}
		}
	}
	seriesMap := make(map[string][]Series)
	for _, metric := range metrics {
		for hostname, points := range pointMap[metric] {
			seriesMap[metric] = append(seriesMap[metric], Series{
				Points: points,
				Name:   hostname,
			})
		}
	}
	return seriesMap
}

func DefaultLookbehind() int64 {
	return time.Now().Unix() - int64(app.Config.Number[config.FIRESTORE_TRUMP_LOOKBACK] * config.SEC_IN_HRS)
}

func LookbehindFor(hours int) int64 {
	return time.Now().Unix() - int64(hours * config.SEC_IN_HRS)
}

func GetGraphData(lookbehind int64) map[string][]Series {
	return parseDataIterator(getDataIteratorSince(lookbehind))
}