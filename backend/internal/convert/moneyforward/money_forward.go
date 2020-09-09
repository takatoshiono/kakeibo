package moneyforward

import (
	"fmt"
	"strconv"
	"time"

	"cloud.google.com/go/civil"

	"github.com/takatoshiono/kakeibo/backend/internal/domain"
)

const (
	csvFieldIsCalculationTarget = iota
	csvFieldRecordedOn
	csvFieldTitle
	csvFieldAmount
	csvFieldSource
	csvFieldCategory1
	csvFieldCategory2
	csvFieldMemo
	csvFieldIsTransfer
	csvFieldID
)

const recordedOnFormat = "2006/01/02"

func ConvCSVToDomain(fields []string) (*domain.MoneyForwardRecord, error) {
	recordedOn, err := convRecordedOnToDomain(fields[csvFieldRecordedOn])
	if err != nil {
		return nil, err
	}

	amount, err := strconv.ParseInt(fields[csvFieldAmount], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse int: %w", err)
	}

	return &domain.MoneyForwardRecord{
		ID:                  fields[csvFieldID],
		IsCalculationTarget: convIsCalculationTargetToDomain(fields[csvFieldIsCalculationTarget]),
		RecordedOn:          recordedOn,
		Title:               fields[csvFieldTitle],
		Amount:              int(amount),
		SourceName:          fields[csvFieldSource],
		Category1Name:       fields[csvFieldCategory1],
		Category2Name:       fields[csvFieldCategory2],
		Memo:                fields[csvFieldMemo],
		IsTransfer:          convIsTransferToDomain(fields[csvFieldIsTransfer]),
	}, nil
}

func convIsCalculationTargetToDomain(v string) bool {
	return v == "1"
}

func convRecordedOnToDomain(v string) (civil.Date, error) {
	t, err := time.Parse(recordedOnFormat, v)
	if err != nil {
		return civil.Date{}, fmt.Errorf("failed to parse: %w", err)
	}
	return civil.DateOf(t), nil
}

func convIsTransferToDomain(v string) bool {
	return v == "1"
}
