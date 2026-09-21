package app

type pageData struct {
	Profile profile
	Scene   string
	Known   bool
	Fog     bool
	Cover   int
	Wind    float64
}

func newPageData(content profile, current conditions, hour int) pageData {
	return pageData{
		Profile: content,
		Scene:   sceneFor(current, hour),
		Known:   current.Known,
		Fog:     current.Known && fogCodes[current.Code],
		Cover:   current.Cover,
		Wind:    current.Wind,
	}
}
