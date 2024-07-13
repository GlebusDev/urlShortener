package random

import (
	"math/rand"
	"time"
)

func NewRandomString(length int) string {

	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	var symbols = make([]byte, length)

	for i := range symbols {
		symbols[i] = byte(rnd.Intn(26) + 65)
	}
	res := string(symbols)
	return res
}
