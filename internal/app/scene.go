package app

const defaultScene = "day"

var snowCodes = map[int]bool{
	71: true, 73: true, 75: true, 77: true, 85: true, 86: true,
}

var rainCodes = map[int]bool{
	51: true, 53: true, 55: true, 56: true, 57: true,
	61: true, 63: true, 65: true, 66: true, 67: true,
	80: true, 81: true, 82: true,
}

var fogCodes = map[int]bool{45: true, 48: true}

func sceneFor(current conditions, hour int) string {
	if !current.Known {
		return defaultScene
	}

	switch {
	case current.Code >= 95:
		return "thunder"
	case snowCodes[current.Code]:
		return "snow"
	case rainCodes[current.Code]:
		return "rain"
	}

	if current.Night {
		return "night"
	}

	switch {
	case hour < 12:
		return "morning"
	case hour < 18:
		return "day"
	default:
		return "evening"
	}
}
