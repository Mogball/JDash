package main

import (
	"jdash/app"
	"jdash/trumptracker"
	"time"
	"jdash/config"
	"fmt"
)

func main() {
	app.Init()
	lookbehindTime := time.Now().Unix() - int64(app.Config.Number[config.FIRESTORE_TRUMP_LOOKBACK] * config.SEC_IN_HRS)
	it := trumptracker.GetDataIteratorSince(lookbehindTime)
	resultMap := trumptracker.ParseDataIterator(it)
	fmt.Println(resultMap)
}
