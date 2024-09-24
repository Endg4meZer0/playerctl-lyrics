package util

import "strconv"

func TimecodeToFloat(timecode string) float64 {
	// [00:00.00]
	if len(timecode) != 10 {
		return -1
	}
	minutes, err := strconv.ParseFloat(timecode[1:3], 64)
	if err != nil {
		return -1
	}
	seconds, err := strconv.ParseFloat(timecode[4:9], 64)
	if err != nil {
		return -1
	}
	return minutes*60.0 + seconds
}
