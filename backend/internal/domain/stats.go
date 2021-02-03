package domain

import (
	"time"
)

type AmountByYearMonth struct {
	Year   int
	Month  time.Month
	Amount int
}

type AmountByYearMonthCategory struct {
	Year       int
	Month      time.Month
	CategoryID string
	Amount     int
}
