package util

import (
	"crudly/util/result"
	"strings"
	"time"
)

const IncomingTimeFormat = "2006-01-02T15:04:05"

func ValidateIncomingTime(timeString string) result.R[time.Time] {
	timeString = strings.TrimSuffix(timeString, "Z")

	parsedTime, err := time.Parse(IncomingTimeFormat, timeString)

	if err != nil {
		return result.Err[time.Time](err)
	}

	return result.Ok(parsedTime)
}
