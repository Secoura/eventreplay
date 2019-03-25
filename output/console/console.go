package console

import (
	"encoding/json"
	"fmt"
	"log"
)

func ProcessEvent(ev string, codec string, identifier string) {
	switch codec {
	case "raw":
		fmt.Printf("%s\n", ev)
	case "json":
		j, _ := json.Marshal(map[string]interface{}{
			"identifier": identifier,
			"_raw": ev,
		})
		fmt.Printf("%s\n", j)
	default:
		log.Fatalf("unknown output codec: '%s'", codec)
	}
}