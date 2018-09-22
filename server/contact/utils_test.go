package contact

import (
	"math/rand"
	"time"
)

var runeSet = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	r := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]rune, n)
	for i := range b {
		b[i] = runeSet[r.Intn(len(runeSet))]
	}

	return string(b)
}
