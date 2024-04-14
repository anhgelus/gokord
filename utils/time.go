package utils

import "math"

// HoursOfUnix returns the hours of a unix timestamp
func HoursOfUnix(unix int64) uint {
	return uint(math.Floor(float64(unix) / (60 * 60)))
}
