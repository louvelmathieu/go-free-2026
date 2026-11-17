// Package config lit la configuration du service depuis l'environnement.
package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"
)

// Config est la configuration complète de provisio. Une seule struct,
// construite une seule fois, injectée partout où il faut.
type Config struct {
	Env            string
	Addr           string
	LogLevel       slog.Level
	StepTimeout    time.Duration
	Stock          int
	LegacyFailRate float64
}

// Load lit l'environnement et applique les défauts.
func Load() (Config, error) {
	cfg := Config{
		Env:            getenv("APP_ENV", "dev"),
		Addr:           getenv("APP_ADDR", ":8080"),
		LogLevel:       slog.LevelDebug,
		StepTimeout:    2 * time.Second,
		Stock:          100,
		LegacyFailRate: 0.1,
	}
	if cfg.Env == "prod" {
		cfg.LogLevel = slog.LevelInfo
	}

	if v := os.Getenv("APP_STEP_TIMEOUT"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return Config{}, fmt.Errorf("APP_STEP_TIMEOUT: %w", err)
		}
		cfg.StepTimeout = d
	}
	if v := os.Getenv("APP_LEGACY_FAIL_RATE"); v != "" {
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return Config{}, fmt.Errorf("APP_LEGACY_FAIL_RATE: %w", err)
		}
		cfg.LegacyFailRate = f
	}
	return cfg, nil
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
