package money

import (
	"fmt"
)

type Money uint64

func (m Money) String() string {
	asFloat := float64(m)
	asFloat = asFloat / 100
	return fmt.Sprintf("£%.2f", asFloat)
}

func (m *Money) Debit(amount Money) Money {
	*m = *m - amount
	return *m
}

func (m *Money) Credit(amount Money) Money {
	*m = *m + amount
	return *m
}
