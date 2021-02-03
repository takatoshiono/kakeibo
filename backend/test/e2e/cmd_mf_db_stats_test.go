package e2e

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/takatoshiono/kakeibo/backend/internal/cmd/mf"
	"github.com/takatoshiono/kakeibo/backend/internal/testutil"
)

func TestCmdMFDBStats_ExpensesByMonth(t *testing.T) {
	c := testutil.MustGetConfig()
	opt := &mf.Option{
		DriverName: c.TestDBDriverName,
		DSN:        c.TestDBDSN,
	}
	cmd := mf.NewCmd(opt)

	// Setup
	cmd.SetArgs([]string{"db", "import", "--file", filenameMF})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		ctx := context.Background()
		db := testutil.MustGetDB()
		testutil.TruncateTable(ctx, t, db, "money_forward_records")
		testutil.TruncateTable(ctx, t, db, "categories")
		testutil.TruncateTable(ctx, t, db, "sources")
	}()

	// Run
	got := testutil.CaptureStdout(t, func() {
		cmd.SetArgs([]string{"db", "stats", "--query", "ExpensesByMonth", "--year", "2020"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})

	// Confirm
	want := `7,2801
`

	if d := cmp.Diff(want, got); d != "" {
		t.Error(d)
	}
}

func TestCmdMFDBStats_ExpensesByMonthAndCategory(t *testing.T) {
	c := testutil.MustGetConfig()
	opt := &mf.Option{
		DriverName: c.TestDBDriverName,
		DSN:        c.TestDBDSN,
	}
	cmd := mf.NewCmd(opt)

	// Setup
	cmd.SetArgs([]string{"db", "import", "--file", filenameMF})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}

	defer func() {
		ctx := context.Background()
		db := testutil.MustGetDB()
		testutil.TruncateTable(ctx, t, db, "money_forward_records")
		testutil.TruncateTable(ctx, t, db, "categories")
		testutil.TruncateTable(ctx, t, db, "sources")
	}()

	// Run
	got := testutil.CaptureStdout(t, func() {
		cmd.SetArgs([]string{"db", "stats", "--query", "ExpensesByMonthAndCategory", "--year", "2020"})
		if err := cmd.Execute(); err != nil {
			t.Fatal(err)
		}
	})

	// Confirm
	want := `7,健康・医療,2370
7,食費,431
`

	if d := cmp.Diff(want, got); d != "" {
		t.Error(d)
	}
}
