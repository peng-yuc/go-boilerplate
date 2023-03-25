package util

// Constants for all supported currencies.
const (
	CAD = "CAD"
	EUR = "EUR"
	USD = "USD"
)

// IsSupportedCurrency returns true if the currency is supported.
func IsSupportedCurrency(currency string) bool {
	switch currency {
	case CAD, EUR, USD:
		return true
	}
	return false
}
