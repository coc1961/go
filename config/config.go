package config

import (
	"os"
	"strings"
)

//Config config
type Config struct {
	FindLimit      int
	DatabaseConfig []string
}

var config *Config

// Get get
func Get() *Config {
	if config == nil {
		//databaseConfig := strings.Split("127.0.0.1|crud", "|")
		databaseConfig := strings.Split("/tmp|crudDB", "|")
		if os.Getenv("DATABASE_CONFIG") != "" {
			databaseConfig = strings.Split(os.Getenv("DATABASE_CONFIG"), "|")
		}

		config = &Config{FindLimit: 1000, DatabaseConfig: databaseConfig}
	}
	return config
}
