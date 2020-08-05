package domain

import (
	"cloud.google.com/go/civil"
)

type MoneyForwardRecord struct {
	// ID
	ID string

	// 計算対象
	IsCalculationTarget bool

	// 日付
	RecordedOn civil.Date

	// 内容
	Title string

	// 金額（円）
	Amount int

	// 保有金融機関
	Source string

	// 大項目
	Category1 string

	// 中項目
	Category2 string

	// メモ
	Memo string

	// 振替
	IsTransfer bool
}
