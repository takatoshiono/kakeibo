package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// Config is the configuration of this application.
type Config struct {
	DBDriverName         string `envconfig:"DB_DRIVER_NAME"`
	DBDSN                string `envconfig:"DB_DSN"`
	MoneyForwardEmail    string `envconfig:"MONEY_FORWARD_EMAIL"`
	MoneyForwardPassword string `envconfig:"MONEY_FORWARD_PASSWORD"`
	TestDBDriverName     string `envconfig:"TEST_DB_DRIVER_NAME"`
	TestDBDSN            string `envconfig:"TEST_DB_DSN"`
}

// Get returns the config.
func Get() (*Config, error) {
	c := &Config{}
	if err := envconfig.Process("", c); err != nil {
		return nil, fmt.Errorf("failed to process config: %w", err)
	}
	return c, nil
}
