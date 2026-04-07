package models

type Currency struct {
	Code    string // "KZT"
	Numeric string // "001" — для номера счёта
}

var currencies = map[string]Currency{
	"KZT": {Code: "KZT", Numeric: "001"},
	"USD": {Code: "USD", Numeric: "002"},
	"EUR": {Code: "EUR", Numeric: "003"},
	"RUB": {Code: "RUB", Numeric: "004"},
	"CNY": {Code: "CNY", Numeric: "005"},
}

func GetCurrency(code string) (Currency, bool) {
	c, ok := currencies[code]
	return c, ok
}

func IsValidCurrency(code string) bool {
	_, ok := currencies[code]
	return ok
}

func SupportedCurrencies() []string {
	list := make([]string, 0, len(currencies))
	for code := range currencies {
		list = append(list, code)
	}
	return list
}
