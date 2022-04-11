package login

import (
	"math/rand"
	"time"
)

const (
	lLetters = "abcdefghijklmnopqrstuvwxyz"
)

func Generate() string {
	rand.Seed(time.Now().UnixNano())
	length := 6
	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = lLetters[rand.Intn(len(lLetters))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}
