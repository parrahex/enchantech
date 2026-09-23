package app

import "math"

type pageData struct {
	Profile     profile
	Scene       string
	Known       bool
	Night       bool
	Fog         bool
	Cover       int
	Wind        float64
	Condition   string
	Temperature int
	Moon        float64
	Season      string
	Clock       bool
	Offset      int
}

func newPageData(content profile, current conditions, now moment) pageData {
	return pageData{
		Profile:     content,
		Scene:       sceneFor(current, now.hour),
		Known:       current.Known,
		Night:       current.Known && current.Night,
		Fog:         current.Known && fogCodes[current.Code],
		Cover:       current.Cover,
		Wind:        current.Wind,
		Condition:   conditionFor(current),
		Temperature: int(math.Round(current.Temperature)),
		Moon:        now.moon,
		Season:      now.season,
		Clock:       now.local,
		Offset:      now.offset,
	}
}
