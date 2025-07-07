package util

import (
	"time"

	"github.com/google/uuid"
)

func GenerateOrderID() uint {
	// Generate a new UUID using the standard library
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)

	uniqueID := uuid.New().ID()

	return uint(currentTimestamp) + uint(uniqueID)

}
