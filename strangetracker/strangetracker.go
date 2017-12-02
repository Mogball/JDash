package strangetracker

import (
	"jdash/config"
	"jdash/app"
)

func buildCode(conf *config.Config, limit int) *CodeSet {
	return &CodeSet{
		A: int64(conf.Number[config.STRANGE_DOM_OFFSET_A]),
		B: int64(conf.Number[config.STRANGE_DOM_OFFSET_B]),
		C: int64(conf.Number[config.STRANGE_DOM_OFFSET_C]),
		Mod: int64(conf.Number[config.STRANGE_DOM_MOD]),
		Limit: int64(limit),
	}
}

func breakAndGetFromConfig() string {
	code := buildCode(app.Config, config.NUM_CHARS)
	sep := app.Config.Word[config.STRANGE_DOM_SEP]
	value := app.Config.Word[config.STRANGE_DOM_STRING]
	return CrackString(value, sep, code)
}


