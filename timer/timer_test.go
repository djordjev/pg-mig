package timer

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestParseAnyTime(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2020-09-20T15:04:05Z")
	t2, _ := time.Parse("2006-01-02T15:04:05", "2020-09-20T15:04:05")
	t3, _ := time.Parse("2006-01-02", "2020-09-20")

	r := require.New(t)

	now := func() time.Time {
		t1, _ := time.Parse(time.RFC3339, "2020-11-01T00:00:00Z")
		return t1
	}

	tNowWithTime, _ := time.Parse(time.RFC3339, "2020-11-01T11:02:00Z")
	timerError := "timer error: unable to parse date/time invalid time"

	table := []struct {
		time string
		res  time.Time
		err  *string
	}{
		{
			time: "2020-09-20T15:04:05Z",
			res:  t1,
			err:  nil,
		},
		{
			time: "2020-09-20T15:04:05",
			res:  t2,
			err:  nil,
		},
		{
			time: "2020-09-20",
			res:  t3,
			err:  nil,
		},
		{
			time: "11:02AM",
			res:  tNowWithTime,
			err:  nil,
		},
		{
			time: "invalid time",
			res:  time.Time{},
			err:  &timerError,
		},
		{
			time: fmt.Sprintf("%d", t1.Unix()),
			res: t1,
			err: nil,
		},
	}

	for _, v := range table {
		t.Run(v.time, func (t *testing.T) {
			currentTimer := Timer{Now: now}
			time, err := currentTimer.ParseTime(v.time)
			isEqual := v.res.Equal(time)
			r.True(isEqual)
			if v.err != nil {
				r.EqualError(err, *v.err)
			} else {
				r.NoError(err)
			}
		})
	}
}
