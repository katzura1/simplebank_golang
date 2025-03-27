package util

func IsSupportedCurrency(currency string) bool {
	supportedCurrencies := []string{"USD", "EUR", "CAD", "IDR"}

	for _, c := range supportedCurrencies {
		if c == currency {
			return true
		}
	}
	return false
}
