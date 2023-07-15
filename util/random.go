package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano()) //because mustbe int64, so must convert using UnixNano()
}

// randomly generate a value between min, max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1) // rand.Int63n(max - min + 1) return a value in rage of [0, max-min]
	// so this will return  a value in range of [min, max]
}

// randomly generate a string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)] //any character of alphabet
		sb.WriteByte(c)             //write this char to the string builder
	}

	return sb.String()
}

// generate random Username using the above
func RandomUsername() string {
	return RandomString((6))
}

// generate random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// generate random currencies among some options
func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD", "SGD", "MYR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// generate random emails
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
