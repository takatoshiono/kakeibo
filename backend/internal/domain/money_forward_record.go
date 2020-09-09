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
	SourceName string

	// 保有金融機関ID
	SourceID string

	// 大項目
	Category1Name string

	// 大項目ID
	Category1ID string

	// 中項目
	Category2Name string

	// 中項目ID
	Category2ID string

	// メモ
	Memo string

	// 振替
	IsTransfer bool
}
