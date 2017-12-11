package strangetracker

import (
	"jdash/app"
	"jdash/config"
	"golang.org/x/net/context"
	"cloud.google.com/go/firestore"
	"time"
	"google.golang.org/api/iterator"
	"log"
)

type Point struct {
	TimeSeconds int64 `json:"x"`
	Value       int64 `json:"y"`
}

type Series struct {
	Points []Point `json:"values"`
	Name   string  `json:"key"`
}

func GetDataIteratorSince(lookbehindSeconds int64) *firestore.DocumentIterator {
	dailyData := app.FirestoreClient().Collection(config.FIRESTORE_STRANGE_TRACKER).Doc(config.FIRESTORE_DOM_DATA).Collection(config.DATA)
	return dailyData.Where("Time", ">=", lookbehindSeconds).Documents(context.Background())
}

func parseDataIterator(dataIt *firestore.DocumentIterator) Series {
	points := make([]Point, 0)
	for {
		doc, err := dataIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Println(err)
			break
		}
		results := doc.Data()
		point := Point{
			TimeSeconds: results["Time"].(int64),
			Value:       results["Count"].(int64),
		}
		points = append(points, point)
	}
	return Series{
		Points: points,
		Name:   "DOM Counts",
	}
}

func LookbehindFor(days int) int64 {
	return time.Now().Unix() - int64(days*config.SEC_IN_DAY)
}

func DefaultLookbehind() int64 {
	return time.Now().Unix() - int64(app.Config().Number[config.FIRESTORE_DOM_LOOKBACK]*config.SEC_IN_DAY)
}

func GetGraphData(lookbehindSeconds int64) map[string][]Series {
	dataSeries := parseDataIterator(GetDataIteratorSince(lookbehindSeconds))
	return map[string][]Series{"default": {dataSeries}}
}
