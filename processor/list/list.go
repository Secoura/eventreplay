package list

import (
	"math/rand"

	"github.com/pkg/errors"
)

func ProcessEvent(replacement interface{}) string {
	list, ok := replacement.([]string)
	if !ok {
		panic(errors.New("unknown replacement settings provided for 'list'"))
	}
	return list[rand.Intn(len(list))]
}
