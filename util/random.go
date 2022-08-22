package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet string = "abcdefghijklmnopqrstuvwxyz"
	digit    string = "1234567"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomPhoneNumber() string {
	var sb strings.Builder
	k := len(digit)

	for i := 0; i < 10; i++ {
		c := digit[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomBoolean() bool {
	return rand.Float32() < 0.5
}

func RandomSubject() string {
	n := len(Subjects)
	subjects := make([]string, n)
	id := 0
	for _, value := range Subjects {
		subjects[id] = value
		id++
	}

	return subjects[rand.Intn(n)]
}
