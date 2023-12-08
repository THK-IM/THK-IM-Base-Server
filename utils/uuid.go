package utils

import (
	"fmt"
	"math/rand"
)

func GetRandomString(length int) string {
	randBytes := make([]byte, length/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
