package data

import (
	"time"
)

const ISO = "2006-01-02T15:04Z05:00"

type Record struct {
	Time time.Time
	Radiation float64
	Humidity float64
	Temperature float64
	Wind float64
	Power float64
	empty bool
	null bool
}
