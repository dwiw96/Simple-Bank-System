package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

/*
 * As rand.Seed() expect an int64 as input,
 * we should convert the time to unix nano before passing it to the function.
 */
func init() {
	rand.Seed(time.Now().UnixNano())
}

/*
 * rand.Int63n(n) function returns a random integer between 0 and n-1.
 * So rand.Int63n(max - min + 1) will return a random integer between 0 and max - min.
 */
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomCurrency() string {
	currencies := []string{"IDR", "USD", "EUR"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomByte(n int) ([]byte, error) {
	res := make([]byte, n)
	len, err := rand.Read(res)
	if len != n {
		return nil, err
	}
	return res, err
}

func RandomPassword() string {
	return RandomString(5) + fmt.Sprint(rand.Intn(10))
}
