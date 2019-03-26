package hex

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/pkg/errors"
)

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func ProcessEvent(replacement interface{}) string {
	l, ok := replacement.(int)
	if !ok {
		panic(errors.New("unknown replacement settings provided for 'string'"))
	}
	h, err := randomHex(l)
	if err != nil {
		panic(err)
	}
	return h
}
