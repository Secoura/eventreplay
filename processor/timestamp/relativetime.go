package timestamp

import (
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"github.com/pkg/errors"
)

type timeUnit int

const (
	Second timeUnit = iota
	Minute
	Hour
	Day
	Week
	Month
	Year
)

func getRelativeTime(relTime string, refTime time.Time) (time.Time, error) {
	if len(relTime) == 0 {
		return time.Time{}, errors.New("failed to parse relative time")
	}

	nowTime := now.New(refTime)

	// Special cases
	switch relTime {
	case "now":
		return nowTime.Time, nil
	case "today":
		return nowTime.BeginningOfDay(), nil
	}

	// s,m,h,d,w,M,y
	// We can cheat by using the ParseDuration function from Go.
	// However, as it only supports up to "hours", so for days, months and years, we will use another implementation.

	unit := relTime[len(relTime)-1]

	if unit == 's' || unit == 'm' || unit == 'h' {
		dur, err := time.ParseDuration(relTime)
		if err != nil {
			return time.Time{}, errors.Wrap(err, "failed to parse relative time")
		}
		return nowTime.Add(dur), nil
	}

	valStr := relTime[0 : len(relTime)-1]
	val, err := strconv.ParseInt(valStr, 10, 64)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to parse relative time")
	}

	if unit == 'd' {
		return nowTime.AddDate(0, 0, int(val)), nil
	}
	if unit == 'w' {
		return nowTime.AddDate(0, 0, int(val*7)), nil
	}
	if unit == 'M' {
		return nowTime.AddDate(0, int(val), 0), nil
	}
	if unit == 'y' {
		return nowTime.AddDate(int(val), 0, 0), nil
	}

	return time.Time{}, errors.Errorf("failed to parse relative time, unknown unit: '%c'", unit)
}