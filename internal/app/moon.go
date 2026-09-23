package app

import (
	"math"
	"time"
)

const synodicMonth = 29.530588853

var referenceNewMoon = time.Date(2000, time.January, 6, 18, 14, 0, 0, time.UTC)

func moonPhase(at time.Time, latitude float64) float64 {
	days := at.Sub(referenceNewMoon).Hours() / 24

	phase := math.Mod(days/synodicMonth, 1)
	if phase < 0 {
		phase++
	}

	if latitude < 0 {
		phase = 1 - phase
	}

	return math.Round(phase*1000) / 1000
}
