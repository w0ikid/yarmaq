package exchange

import (
	"fmt"
	"math"

	"github.com/w0ikid/yarmaq/pkg/errs"
)

// Service handles currency conversions with KZT as the anchor currency.
type Service interface {
	Convert(amount int64, from, to string) (int64, float64, error)
}

type implementation struct {
	// rates maps currency code to its value in KZT (1 unit = X KZT)
	rates map[string]float64
}

func NewService() Service {
	return &implementation{
		rates: map[string]float64{
			"KZT": 1.0,
			"USD": 480.0,
			"EUR": 520.0,
			"RUB": 5.2,
			"CNY": 65.0,
		},
	}
}

func (s *implementation) Convert(amount int64, from, to string) (int64, float64, error) {
	if from == to {
		return amount, 1.0, nil
	}

	fromRate, ok := s.rates[from]
	if !ok {
		return 0, 0, fmt.Errorf("%w: unsupported source currency: %s", errs.ErrValidation, from)
	}

	toRate, ok := s.rates[to]
	if !ok {
		return 0, 0, fmt.Errorf("%w: unsupported target currency: %s", errs.ErrValidation, to)
	}

	// Conversion logic:
	// AmountInKZT = amount * fromRate
	// TargetAmount = AmountInKZT / toRate
	// EffectiveRate = fromRate / toRate

	effectiveRate := fromRate / toRate
	targetAmount := float64(amount) * effectiveRate

	return int64(math.Round(targetAmount)), effectiveRate, nil
}
