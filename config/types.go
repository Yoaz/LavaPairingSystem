package config

import (
	"log/slog"

	"github.com/Yoaz/LavaPairingSystem/internal/filter"
	"github.com/Yoaz/LavaPairingSystem/internal/score"
	"github.com/Yoaz/LavaPairingSystem/internal/system"
)

// AppConfig holds the configuration for the application, including filters, scorers, and the pairing system
type AppConfig struct {
	Log           *slog.Logger
	Filters       []filter.Filter
	Scorers       []score.Scorer
	PairingSystem system.PairingSystem
}
