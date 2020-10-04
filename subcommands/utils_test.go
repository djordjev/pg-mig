package subcommands

import (
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

	table := []struct {
		time string
		res  time.Time
		err  error
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
			err:  UnableToParseTimeErr,
		},
	}

	for _, v := range table {
		t, err := parseAnyTime(v.time, now)
		r.Equal(t, v.res)
		r.Equal(err, v.err)
	}
}
