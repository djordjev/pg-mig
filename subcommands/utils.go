package subcommands

import (
	"errors"
	"time"
)

var UnableToParseTimeErr = errors.New("unable to parse date/time format")

const (
	dateNoTimezone = "2006-01-02T15:04:05"
	dateNoTime     = "2006-01-02"
	onlyTime       = "3:04PM"
)

func parseAnyTime(inputTime string, timeGetter TimeGetter) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		dateNoTimezone,
		dateNoTime,
		onlyTime,
	}

	for _, format := range formats {
		t, err := time.Parse(format, inputTime)
		if err == nil {
			if format == onlyTime {
				now := timeGetter()
				todayTime := time.Date(
					now.Year(),
					now.Month(),
					now.Day(),
					t.Hour(),
					t.Minute(),
					t.Second(),
					0,
					now.Location(),
				)

				return todayTime, nil
			}

			return t, nil
		}
	}

	return time.Time{}, UnableToParseTimeErr
}
