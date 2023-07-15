package util

const (
	USD = "USD"
	SGD = "SGD"
	EUR = "EUR"
	CAD = "CAD"
	MYR = "MYR"
)

// return true if currency is supported
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, SGD, EUR, CAD, MYR:
		return true
	}
	return false
}
