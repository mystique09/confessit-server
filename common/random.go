package common

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"
)

const characters = "abcdefghijklmnopqrstuvwxyz0123456789"

func RandomString(l int) string {
	var sb strings.Builder

	k := len(characters)

	for i := 0; i < l; i++ {
		randn, err := rand.Int(rand.Reader, big.NewInt(int64(k)))
		if err != nil {
			log.Fatal(err)
		}

		char := characters[randn.Int64()]
		sb.WriteByte(char)
	}
	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", RandomString(12), RandomString(6))
}
