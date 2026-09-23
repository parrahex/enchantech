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
}

func newPageData(content profile, current conditions, hour int) pageData {
	return pageData{
		Profile:     content,
		Scene:       sceneFor(current, hour),
		Known:       current.Known,
		Night:       current.Known && current.Night,
		Fog:         current.Known && fogCodes[current.Code],
		Cover:       current.Cover,
		Wind:        current.Wind,
		Condition:   conditionFor(current),
		Temperature: int(math.Round(current.Temperature)),
	}
}
