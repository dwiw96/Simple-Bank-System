package util

import (
	"time"
)

// constants for all supported currencies.
const (
	IDR = "IDR"
	USD = "USD"
	EUR = "EUR"
	YEN = "YEN"
)

// func to check if the input currency is supported, renturn "true" if supported.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case IDR, USD, EUR, YEN:
		return true
	}
	return false
}

func CurrencyExchangeToIDR(currency string, amount int64) int64 {
	switch currency {
	case "USD":
		return amount * 14500
	case "EUR":
		return amount * 16500
	case "YEN":
		return amount * 100
	}

	return 0
}

func GetDOB(input string) (time.Time, error) {
	//YYYYMMDD := "2022-01-20"
	t, err := time.Parse(time.DateOnly, input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
