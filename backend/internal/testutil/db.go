package testutil

import (
	"database/sql"
	"sync"

	"github.com/takatoshiono/kakeibo/backend/internal/config"
)

var (
	cfg     *config.Config
	cfgOnce sync.Once
)

// MustGetConfig gets the config only once and returns it or panic if error occured.
func MustGetConfig() *config.Config {
	cfgOnce.Do(func() {
		var err error
		if cfg, err = config.Get(); err != nil {
			panic("failed to get config: " + err.Error())
		}
	})
	return cfg
}

// MustGetDB returns *sql.DB or panic if error occured.
func MustGetDB() *sql.DB {
	c := MustGetConfig()
	db, err := sql.Open(c.TestDBDriverName, c.TestDBDSN)
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	return db
}
