package mac

import (
	"crypto/rand"
	"fmt"
)

func ProcessEvent() string {
	// 6e:0c:51:c6:c6:3a
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return "6e:0c:51:c6:c6:3a" // default
	}
	buf[0] |= 2
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}
