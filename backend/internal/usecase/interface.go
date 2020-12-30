package usecase

import (
	"context"
	"io"
	"time"

	"github.com/takatoshiono/kakeibo/backend/internal/domain"
)

// GoogleDrive is a interface that interact with google drive.
type GoogleDrive interface {
	CreateFile(ctx context.Context, r io.Reader, filename string, ctype string, parentID string) (string, error)
}

// MoneyForward is a interface tha interact with MoneyForward ME.
type MoneyForward interface {
	Login(ctx context.Context) error
	DownloadCSV(ctx context.Context, year int, month time.Month) (io.ReadCloser, error)
}

// MoneyForwardCSVReader is a interface that reads CSV data.
type MoneyForwardCSVReader interface {
	Read() ([]string, error)
	ReadAll() ([][]string, error)
}

// Transaction is a interface that executes a database transaction.
type Transaction interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// MasterRepository is a interface of MasterRepository.
type MasterRepository interface {
	FindOrCreateSource(ctx context.Context, name string) (*domain.Source, error)
	FindOrCreateCategory(ctx context.Context, name string, level domain.CategoryLevel, parentID string) (*domain.Category, error)
}

// MoneyForwardRepository is a interface of MoneyForwardRepository.
type MoneyForwardRepository interface {
	CreateOrUpdateRecord(ctx context.Context, record *domain.MoneyForwardRecord) error
}

// StatsRepository is a interface of StatsRepository.
type StatsRepository interface {
	FindAmountExpendedByMonth(ctx context.Context, year int) ([]*domain.AmountExpendedByMonth, error)
}
