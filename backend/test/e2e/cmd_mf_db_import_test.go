package e2e

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	"github.com/takatoshiono/kakeibo/backend/internal/cmd/mf"
	"github.com/takatoshiono/kakeibo/backend/internal/domain"
	"github.com/takatoshiono/kakeibo/backend/internal/repository/database"
	"github.com/takatoshiono/kakeibo/backend/internal/testutil"
)

func TestCmdMFDBImport(t *testing.T) {
	// TODO: CreateのケースとUpdateのケースを明示的にテストする
	c := testutil.MustGetConfig()
	opt := &mf.Option{
		DriverName: c.TestDBDriverName,
		DSN:        c.TestDBDSN,
	}
	cmd := mf.NewCmd(opt)
	cmd.SetArgs([]string{"db", "import", "--file", filenameMF})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	expected := []*domain.MoneyForwardRecord{
		{
			ID:                  "5nNlbKBTtKW2DJOK8vIN8g",
			IsCalculationTarget: true,
			RecordedOn:          time.Date(2020, 7, 10, 0, 0, 0, 0, time.UTC),
			Title:               "ジヨイナス",
			Amount:              -431,
			SourceName:          "ジャックスカード",
			SourceID:            "",
			Category1Name:       "食費",
			Category1ID:         "",
			Category2Name:       "食費",
			Category2ID:         "",
			Memo:                "ミスド",
			IsTransfer:          false,
		},
		{
			ID:                  "X_UMG0Ztge8nKZl01uiVZw",
			IsCalculationTarget: true,
			RecordedOn:          time.Date(2020, 7, 11, 0, 0, 0, 0, time.UTC),
			Title:               "歯医者",
			Amount:              -2370,
			SourceName:          "家族の財布",
			SourceID:            "",
			Category1Name:       "健康・医療",
			Category1ID:         "",
			Category2Name:       "医療費",
			Category2ID:         "",
			Memo:                "",
			IsTransfer:          false,
		},
	}
	notExpected := []*domain.MoneyForwardRecord{
		{
			ID:                  "CETj9SskFWNAoj_d6GkWhQ",
			IsCalculationTarget: false,
			RecordedOn:          time.Date(2020, 7, 12, 0, 0, 0, 0, time.UTC),
			Title:               "定期預金",
			Amount:              -10000,
			SourceName:          "三井住友銀行",
			SourceID:            "",
			Category1Name:       "未分類",
			Category1ID:         "",
			Category2Name:       "未分類",
			Category2ID:         "",
			Memo:                "",
			IsTransfer:          true,
		},
	}

	db := testutil.MustGetDB()
	transaction := database.NewTransaction(db)
	masterRepo := database.NewMasterRepository(transaction)
	mfRepo := database.NewMoneyForwardRepository(transaction)
	ctx := context.Background()

	defer func() {
		testutil.TruncateTable(ctx, t, db, "money_forward_records")
		testutil.TruncateTable(ctx, t, db, "categories")
		testutil.TruncateTable(ctx, t, db, "sources")
	}()

	for _, want := range expected {
		s, err := masterRepo.FindSourceByName(ctx, want.SourceName)
		if err != nil {
			t.Fatal(err)
		}
		c1, err := masterRepo.FindCategoryByNameAndLevel(ctx, want.Category1Name, domain.CategoryLevel1)
		if err != nil {
			t.Fatal(err)
		}
		c2, err := masterRepo.FindCategoryByNameAndLevel(ctx, want.Category2Name, domain.CategoryLevel2)
		if err != nil {
			t.Fatal(err)
		}

		want.SourceID = s.ID
		want.Category1ID = c1.ID
		want.Category2ID = c2.ID

		got, err := mfRepo.FindRecord(ctx, want.ID)
		if err != nil {
			t.Fatal(err)
		}

		opts := cmp.Options{
			// 今のところなくていいが、将来的に必要になったら取れるようにする
			cmpopts.IgnoreFields(domain.MoneyForwardRecord{},
				"IsCalculationTarget", // not saved
				"SourceName",          // not get
				"Category1Name",       // not get
				"Category1ID",         // not get
				"Category2Name",       // not get
				"IsTransfer",          // not saved
			),
		}
		if d := cmp.Diff(want, got, opts); d != "" {
			t.Error(d)
		}
	}

	for _, want := range notExpected {
		_, err := mfRepo.FindRecord(ctx, want.ID)
		if !errors.Is(err, sql.ErrNoRows) {
			t.Errorf("want err no rows, but got err %v", err)
		}
	}
}
