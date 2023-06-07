package util

// constants for all supported currencies.
const (
	IDR = "IDR"
	USD = "USD"
	EUR = "EUR"
)

// func to check if the input currency is supported, renturn "true" if supported.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case IDR, USD, EUR:
		return true
	}
	return false
}
