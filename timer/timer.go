package timer

import (
	"fmt"
	"time"
)

// TimeGetter
type TimeGetter func() time.Time

type Timer struct {
	Now TimeGetter
}

const (
	dateNoTimezone = "2006-01-02T15:04:05"
	dateNoTime     = "2006-01-02"
	onlyTime       = "3:04PM"
)

func (timer Timer) ParseTime(inputTime string) (time.Time, error) {
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
				now := timer.Now()
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

	return time.Time{}, fmt.Errorf("timer error: unable to parse date/time %s", inputTime)
}
