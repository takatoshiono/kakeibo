package domain

import (
	"time"
)

type AmountExpendedByMonth struct {
	Year   int
	Month  time.Month
	Amount int
}
