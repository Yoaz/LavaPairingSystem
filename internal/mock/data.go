package mock

import (
	pairing "github.com/Yoaz/LavaPairingSystem/internal"
)

var (
	// Mocked Providers
	Providers = []*pairing.Provider{
		{ID: "1", Address: "provider1", Stake: 1000, Location: "US-West", Features: []string{"featA", "featB", "featC"}, Fee: 3.0},
		{ID: "2", Address: "provider2", Stake: 2000, Location: "US-East", Features: []string{"featA", "featB"}, Fee: 0.015},
		{ID: "3", Address: "provider3", Stake: 1500, Location: "EU-Central", Features: []string{"featA", "featC", "featD"}, Fee: 4.5},
		{ID: "4", Address: "provider4", Stake: 500, Location: "US-West", Features: []string{"featB"}, Fee: 0.005},
		{ID: "5", Address: "provider5", Stake: 2500, Location: "US-West", Features: []string{"featA", "featB", "featC", "featExtra"}, Fee: 0.8},
		{ID: "6", Address: "provider6", Stake: 1200, Location: "EU-Central", Features: []string{"featA", "featD", "featE"}, Fee: 1.7},
		{ID: "7", Address: "provider7", Stake: 800, Location: "US-East", Features: []string{"featA", "featB", "featC", "featX"}, Fee: 2.0},
		{ID: "8", Address: "provider8", Stake: 3000, Location: "US-West", Features: []string{"featA", "featB", "featC", "featY", "featZ"}, Fee: 2.5},
	}

	// Mocked Consumer Policy
	ConsumerPolicy = &pairing.ConsumerPolicy{
		RequiredLocation: "US-West",
		RequiredFeatures: []string{"featA", "featB"},
		MinStake:         1000,
		Weights: map[string]float64{
			"StakeScore":    0.5,
			"FeatureScore":  0.3,
			"LocationScore": 0.2,
			// NOTE: Commented in order to show ommiting of a score type therefore contributing 0 to the final score
			// "FeeScore":      0.1,
		},
	}
)
