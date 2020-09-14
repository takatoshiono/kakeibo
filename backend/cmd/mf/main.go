package main

import (
	"fmt"
	"os"

	"github.com/takatoshiono/kakeibo/backend/internal/cmd/mf"
)

func main() {
	if err := mf.NewCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
