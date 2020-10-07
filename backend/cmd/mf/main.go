package main

import (
	"fmt"
	"os"

	"github.com/takatoshiono/kakeibo/backend/internal/cmd/mf"
	"github.com/takatoshiono/kakeibo/backend/internal/config"
)

func main() {
	if err := realMain(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func realMain() error {
	c, err := config.Get()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	if err := mf.NewCmd(&mf.Option{
		DriverName:           c.DBDriverName,
		DSN:                  c.DBDSN,
		MoneyForwardEmail:    c.MoneyForwardEmail,
		MoneyForwardPassword: c.MoneyForwardPassword,
	}).Execute(); err != nil {
		return fmt.Errorf("failed to execute: %w", err)
	}

	return nil
}
