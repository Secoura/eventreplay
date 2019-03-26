package float

import (
	"strconv"

	"github.com/brianvoe/gofakeit"
	"github.com/pkg/errors"
)

func ProcessEvent(replacement interface{}) string {
	conf, ok := replacement.([]float64)
	if !ok || len(conf) != 3 {
		panic(errors.New("unknown replacement settings provided for 'float'"))
	}
	precision := conf[2]
	randFloat := gofakeit.Float64Range(conf[0], conf[1])
	return strconv.FormatFloat(randFloat, 'f', int(precision), 64)
}
