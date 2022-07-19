package passutil

import (
	"math/rand"
	"time"
)

func Generate() string {
	const passwordLen = 32
	const passwordChar = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	rand.Seed(time.Now().UnixNano())

	secureBytes := make([]byte, passwordLen)
	rand.Read(secureBytes)

	for k, v := range secureBytes {
		secureBytes[k] = passwordChar[v%byte(len(passwordChar))]
	}

	return string(secureBytes)
}
