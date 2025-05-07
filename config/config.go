package config

import (
	"log/slog"

	"github.com/Yoaz/LavaPairingSystem/internal/filter"
	"github.com/Yoaz/LavaPairingSystem/internal/logger"
	"github.com/Yoaz/LavaPairingSystem/internal/score"
	"github.com/Yoaz/LavaPairingSystem/internal/system"
)

// Init initializes the application configuration, including filters, scorers, and the pairing system
// It takes a strictMode boolean to determine if strict mode is enabld and a logLevel for logging
func Init(strictMode bool, logLevel slog.Level) *AppConfig {
	log := logger.NewWithLevel(logLevel)
	log.Info("Initializing LavaPairingSystem...")

	filters := []filter.Filter{
		filter.LocationFilter{},
		filter.FeatureFilter{},
		filter.StakeFilter{},
	}
	log.Debug("Initialized filters", "count", len(filters))

	scorers := []score.Scorer{
		&score.StakeScore{},
		&score.FeatureScore{},
		&score.LocationScore{},
		&score.FeeScore{},
	}
	log.Debug("Initialized scorers", "count", len(scorers))

	pairingSystem := system.NewPairingSystem(filters, scorers, log, strictMode)
	log.Info("Pairing system initialized successfully.")

	return &AppConfig{
		Log:           log,
		Filters:       filters,
		Scorers:       scorers,
		PairingSystem: pairingSystem,
	}
}
