package domain

import (
	"time"
)

type AmountByYearMonth struct {
	Year   int
	Month  time.Month
	Amount int
}
