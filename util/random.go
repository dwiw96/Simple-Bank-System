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
	return RandomInt(10, 1000)
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomCurrency() string {
	currencies := []string{"IDR", "USD", "EUR", "YEN"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

func RandomEmailWithUsername(username string) string {
	return fmt.Sprintf("%s@email.com", username)
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

const idNumb = "0123456789"

func RandomAccountID() string {
	var sb strings.Builder
	sb.WriteString("1010")
	k := len(alphabet)

	for i := 0; i < 8; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomWalletID(n string) string {
	var sb strings.Builder
	sb.WriteString(n)
	k := len(alphabet)

	for i := 0; i < 8; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomAdress() (string, string, int64) {
	provinces := []string{"Banten", "DKI Jakarta", "Jogjakarta", "Kalimantan Timur", "Papua Barat"}
	city := []string{"Pandeglang", "Serang", "Jakarta Barat", "Gunung Kidul", "Balikpapan", "Tenggarong", "Manokwari", "Fakfak"}

	provLen := len(provinces)
	cityLen := len(city)

	return provinces[rand.Intn(provLen)], city[rand.Intn(cityLen)], RandomInt(1000, 9999)
}

func RandomDate() string {
	startDate := time.Date(1960, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)

	duration := endDate.Sub(startDate)
	randomDuration := time.Duration(rand.Int63n(int64(duration)))
	randomDate := startDate.Add(randomDuration)

	return randomDate.Format("2006-01-02")
}
