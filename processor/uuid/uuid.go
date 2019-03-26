package uuid

import "github.com/google/uuid"

func ProcessEvent() string {
	u := uuid.New()
	return u.String()
}
