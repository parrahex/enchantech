package app

import "time"

func seasonFor(at time.Time, latitude float64) string {
	month := int(at.Month())
	if latitude < 0 {
		month = (month+5)%12 + 1
	}

	switch {
	case month == 12 || month <= 2:
		return "winter"
	case month <= 5:
		return "spring"
	case month <= 9:
		return "summer"
	default:
		return "autumn"
	}
}
