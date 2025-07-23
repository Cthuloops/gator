package config

import (
	"gator/internal/database"
)

type State struct {
	Config *Config
	DB     *database.Queries
}
