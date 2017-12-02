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

	config.Word[STRANGE_DOM_RAW_COUNT] = "^(\\([0-9,]*\\))$"
	config.Word[STRANGE_DOM_STRING] = "225|249|183|229|31|31|78|144|69|30|102|54|222|124|59|181|111|153|222|144|183|79|19|237|196|233|200|9|19|194|155|242|194|236|111|89|115|124|116|51|48"
	config.Word[STRANGE_DOM_SEP] = "|"
	config.Number[STRANGE_DOM_OFFSET_A] = 11
	config.Number[STRANGE_DOM_OFFSET_B] = 5
	config.Number[STRANGE_DOM_OFFSET_C] = 41
	config.Number[STRANGE_DOM_MOD] = 130968


	return &config
}
