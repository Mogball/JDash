package config

type Config struct {
	Word    map[string]string
	Number  map[string]int
	Decimal map[string]float32
}

func Make() *Config {
	config := Config{}
	config.Word = make(map[string]string)
	config.Number = make(map[string]int)
	config.Decimal = make(map[string]float32)

	config.Word[TRUMP_SITES] = "https://www.bloomberg.com," +
		"https://www.economist.com," +
		"https://www.nytimes.com," +
		"https://www.washingtonpost.com," +
		"https://www.usatoday.com," +
		"http://www.washingtonexaminer.com," +
		"https://www.huffingtonpost.com," +
		"https://news.vice.com/en_us," +
		"https://www.theglobeandmail.com"
	config.Word[TRUMP_FULL_MATCHER] = "([D|d]onald [T|t]rump)"
	config.Word[TRUMP_PART_MATCHER] = "([T|t]rump)|([D|d]onald)"

	config.Number[FIRESTORE_TRUMP_LOOKBACK] = 336

	return &config
}
