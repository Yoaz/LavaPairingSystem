package utils

import (
	"fmt"

	pairing "github.com/Yoaz/LavaPairingSystem/internal"
)

// Min returns the minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Compute max stake from a list of providers
func ComputeMaxStake(providers []*pairing.Provider) int64 {
	var maxStake int64
	for _, p := range providers {
		if p.Stake > maxStake {
			maxStake = p.Stake
		}
	}
	return maxStake
}

// ComputeNormalizedFees computes the normalized fees for a list of providers
// This function normalizes the fee of each provider in the list by scaling it
// relative to the maximum fee in the list. The normalized fee is calculated as
// the provider's fee divided by the maximum fee, ensuring that the highest fee
// becomes 1 and all other fees are scaled accordingly
func ComputeNormalizedFees(providers []*pairing.Provider) map[string]float64 {
	// Step 1: Find the maximum fee in the list of providers
	var maxFee float64
	for _, p := range providers {
		// Update maxFee if the current provider's fee is greater
		if p.Fee > maxFee {
			maxFee = p.Fee
		}
	}

	// Step 2: If the maximum fee is 0, set it to 1 to avoid division by zero
	if maxFee == 0 {
		maxFee = 1 // Default to 1 to avoid division by zero
	}

	// Step 3: Create a map to store the normalized fee for each provider
	normalized := make(map[string]float64)
	for _, p := range providers {
		// Normalize each provider's fee by dividing their fee by the max fee
		normalized[p.ID] = float64(p.Fee) / float64(maxFee)
	}

	// Step 4: Return the map containing each provider's normalized fee
	return normalized
}

// ValidateWeights checks if the sum of weights in the given map equals 1.0
// The presence of all specific keys is NOT mandetory, allowing users to provide
// weights only for the components they care about. Unspecified components will effectively
// have a weight of 0 in the weighted scoring logic
func ValidateWeights(weights map[string]float64) error {
	// If weights map is nil or empty, it's considered valid (will fallback to average scoring)
	// Or, if non-empty, the sum must be 1.0
	// NOTE: Defined in struct as a map[string]float64 therefore no need to check for nil
	if len(weights) == 0 {
		return nil // No weights provided, valid for average scoring fallback
	}

	// Only check the sum if weights are provided
	if err := checkWeightSum(weights); err != nil {
		return err
	}
	return nil
}

// CheckWeightSum checks if the sum of weights in the given map equals 1.0
func checkWeightSum(weights map[string]float64) error {
	var total float64
	for _, weight := range weights {
		total += weight
	}
	if total != 1.0 {
		return fmt.Errorf("weights must sum to 1, got %.2f", total)
	}
	return nil
}
