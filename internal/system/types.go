package system

import (
	"log/slog"

	pairing "github.com/Yoaz/LavaPairingSystem/internal"
	"github.com/Yoaz/LavaPairingSystem/internal/filter"
	"github.com/Yoaz/LavaPairingSystem/internal/score"
)

// topN is the number of top providers to return
const (
	topNProviders           = 5
	parallelFilterThreshold = 50
	workerCount             = 10
)

// NewPairingSystem creates a new PairingSystem instance with the provided filters, scorers, and logger
type PairingSystem interface {
	// FilterProviders returns a list of providers that match the policy requirements
	FilterProviders(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) []*pairing.Provider
	// RankProviders assigns scores to providers based on the policy requirements
	RankProviders(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) []*pairing.PairingScore
	// GetPairingList returns the top-5 best provider for the given consumer policy
	GetPairingList(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) ([]*pairing.Provider, error)
}

// pairingSystem is the implementation of the PairingSystem interface
type pairingSystem struct {
	filters    []filter.Filter
	scorers    []score.Scorer
	logger     *slog.Logger
	strictMode bool // If true, returns error when no providers match; if false, returns empty list
}
