package string

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/pkg/errors"
)

func ProcessEvent(replacement interface{}) string {
	l, ok := replacement.(int)
	if !ok {
		panic(errors.New("unknown replacement settings provided for 'string'"))
	}
	return randomdata.RandStringRunes(l)
}
