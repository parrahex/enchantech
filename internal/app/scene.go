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

var conditionNames = map[int]string{
	0:  "clear",
	1:  "mostly clear",
	2:  "partly cloudy",
	3:  "overcast",
	45: "fog",
	48: "fog",
	51: "light drizzle",
	53: "drizzle",
	55: "heavy drizzle",
	56: "freezing drizzle",
	57: "freezing drizzle",
	61: "light rain",
	63: "rain",
	65: "heavy rain",
	66: "freezing rain",
	67: "freezing rain",
	71: "light snow",
	73: "snow",
	75: "heavy snow",
	77: "snow grains",
	80: "light showers",
	81: "showers",
	82: "heavy showers",
	85: "snow showers",
	86: "heavy snow showers",
	95: "thunderstorm",
	96: "thunderstorm with hail",
	99: "thunderstorm with hail",
}

func conditionFor(current conditions) string {
	if !current.Known {
		return ""
	}

	return conditionNames[current.Code]
}

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
