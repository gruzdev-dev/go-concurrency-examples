package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	GoroutinesCount int    `json:"goroutines_count"`
	Iterations      int    `json:"iterations"`
	Mode            string `json:"mode"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.GoroutinesCount <= 0 {
		return nil, fmt.Errorf("goroutines_count must be positive")
	}
	if cfg.Iterations <= 0 {
		return nil, fmt.Errorf("iterations must be positive")
	}

	switch cfg.Mode {
	case "unsafe", "mutex", "atomic", "channel", "mutex_copy", "interface_tearing":
	default:
		return nil, fmt.Errorf("invalid mode: %s", cfg.Mode)
	}

	return &cfg, nil
}
