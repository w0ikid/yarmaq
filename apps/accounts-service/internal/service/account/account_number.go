package account

import (
	"math/big"
	"github.com/w0ikid/yarmaq/pkg/models"
	"fmt"
)

func calcCheckDigits(bban string) string {
	// BBAN + "KZ00" → переводим буквы в цифры (A=10, B=11, ... Z=35)
	raw := bban + "2035" + "00" // K=20, Z=35

	num := new(big.Int)
	num.SetString(raw, 10)

	mod := new(big.Int)
	mod.Mod(num, big.NewInt(97))

	check := 98 - mod.Int64()
	return fmt.Sprintf("%02d", check)
}

func generateAccountNumber(currency string, seq int64) (string, error) {
    c, ok := models.GetCurrency(currency)
    if !ok {
        return "", fmt.Errorf("unsupported currency: %s", currency)
    }
    bban := fmt.Sprintf("02%s%08d", c.Numeric, seq)
    check := calcCheckDigits(bban)
    return fmt.Sprintf("KZ%s%s", check, bban), nil
}
