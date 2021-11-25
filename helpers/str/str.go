package str

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

// Random Random Unique String
func Random(length int) string {
	buf := make([]byte, length/2)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%x", buf)
}

// UUID Generate UUID version 4
func UUID() string {
	return uuid.Must(uuid.NewRandom()).String()
}
