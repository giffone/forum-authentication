package password

import (
	"math/rand"
	"time"
)

const (
	lLetters = "abcdefghijklmnopqrstuvwxyz"
	uLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits   = "0123456789"
)

func Generate() string {
	rand.Seed(time.Now().UnixNano())
	length := 8
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = uLetters[rand.Intn(len(uLetters))]
	for i := 2; i < length; i++ {
		buf[i] = lLetters[rand.Intn(len(lLetters))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	return string(buf)
}
