package integer

import (
	"strconv"

	"github.com/Pallinder/go-randomdata"
	"github.com/pkg/errors"
)

func ProcessEvent(replacement interface{}) string {
	conf, ok := replacement.([]int)
	if !ok || len(conf) != 2 {
		panic(errors.New("unknown replacement settings provided for 'integer'"))
	}
	return strconv.FormatInt(int64(randomdata.Number(conf[0], conf[1])), 10)
}
