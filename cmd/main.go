package main

import (
	// Added for Provider and ConsumerPolicy types
	"log/slog"

	"github.com/Yoaz/LavaPairingSystem/config"
	"github.com/Yoaz/LavaPairingSystem/internal/mock"
	"github.com/Yoaz/LavaPairingSystem/internal/utils"
)

func main() {
	// Initialize with logger `debug` level && strict mode enabled
	app := config.Init(true, slog.LevelDebug)
	log := app.Log

	// --- Example Usage  ---
	providers := mock.Providers   // Mock data for providers
	policy := mock.ConsumerPolicy // Mock data for consumer policy

	// Making sure consumer policy assigned weights is valid
	err := utils.ValidateWeights(policy.Weights)
	if err != nil {
		log.With("error", err).Error("Invalid weights in consumer policy")
		return
	}

	log.Info("Attempting to get pairing list with mock data", "policy_location", policy.RequiredLocation, "policy_min_stake", policy.MinStake, "policy_features_count", len(policy.RequiredFeatures))

	topProviders, err := app.PairingSystem.GetPairingList(providers, policy)
	if err != nil { // If strict mode is enabled, expect an error if no providers match the policy
		log.With("error", err).Error("Failed to get pairing list")
	} else {
		log.Info("Successfully retrieved pairing list", "count", len(topProviders))
		log.Info("----------------------------- TOP PROVIDERS -----------------------------")
		if len(topProviders) == 0 {
			log.Info("No providers matched the policy and were selected.")
		}
		for i, p := range topProviders {
			if p != nil { // nil check for safety, though GetPairingList should not return nil providers in the list
				log.Info("Selected Provider", "rank", i+1, "ID", p.ID, "address", p.Address, "stake", p.Stake, "location", p.Location, "fee", p.Fee, "features", p.Features)
			}
		}
	}
	// --- End Example Usage ---

}
