package filedriller

import "github.com/google/uuid"

// GetUUID returns a UUID v4 as a string
func CreateUUID() string {
	newuuid, err := uuid.NewRandom()
	e(err)

	return newuuid.String()
}
