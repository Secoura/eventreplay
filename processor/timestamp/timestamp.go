package timestamp

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/secoura/jodaTime"
)

func randomTime(start, end time.Time) time.Time {
	startUnix := start.Unix()
	endUnix := end.Unix()
	newTime := startUnix + int64(rand.Intn(int(endUnix) - int(startUnix)))
	return time.Unix(newTime, 0)
}

func ProcessEvent(replacement interface{}, earliestTime, latestTime string) string {
	now := time.Now()
	earliest, err := getRelativeTime(earliestTime, now)
	if err != nil {
		log.Fatalf("Could not parse earliest time: %v", err)
	}
	latest, err := getRelativeTime(latestTime, now)
	if err != nil {
		log.Fatalf("Could not parse latest time: %v", err)
	}

	t := randomTime(earliest, latest)
	var replacementVal string

	switch replacement {
	case "UNIX":
		replacementVal = fmt.Sprintf("%d", t.Unix())
	case "RFC3339", "ISO8601":
		replacementVal = t.Format(time.RFC3339)
	default:
		replacementVal = jodaTime.Format(replacement.(string), t)
	}

	return replacementVal
}