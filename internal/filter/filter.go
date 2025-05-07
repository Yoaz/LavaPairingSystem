package filter

import pairing "github.com/Yoaz/LavaPairingSystem/internal"

/* ***********************************************************************
 *                            LOCATION FILTER                            *
 *********************************************************************** */

// Apply filters providers based on exact match with the required location in the policy
// It retains only those providers whose Location field matches the policy's RequiredLocation
func (f LocationFilter) Apply(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) []*pairing.Provider {
	var result []*pairing.Provider
	for _, p := range providers {
		if p.Location == policy.RequiredLocation {
			result = append(result, p)
		}
	}
	return result
}

// ApplySingle checks if a single provider matches the required location in the policy
// It returns true if the provider's Location field matches the policy's RequiredLocation
func (f LocationFilter) ApplySingle(provider *pairing.Provider, policy *pairing.ConsumerPolicy) bool {
	return provider.Location == policy.RequiredLocation
}

func (f LocationFilter) Name() string { return "LocationFilter" }

/* ***********************************************************************
 *                            FEATURE FILTER                             *
 *********************************************************************** */

// Apply filters providers ensuring they support all features specified in the policy's RequiredFeatures
// It retains only those providers whose Features list contains every feature listed in RequiredFeatures
func (f FeatureFilter) Apply(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) []*pairing.Provider {
	// Use a map for efficient lookup of required features.
	required := make(map[string]bool)
	for _, feature := range policy.RequiredFeatures {
		required[feature] = true
	}

	var result []*pairing.Provider
outer: // Label for breaking out of nested loops
	for _, p := range providers {
		// Check if the provider has all required features
		for req := range required {
			found := false
			for _, pf := range p.Features {
				if pf == req {
					found = true
					break // Found this required feature, check the next one
				}
			}
			// If a required feature wasn't found in the provider's list, skip this provider
			if !found {
				continue outer // Go to the next provider
			}
		}
		// If all required features were found, add the provider to the result list
		result = append(result, p)
	}
	return result
}

// ApplySingle checks if a single provider matches the required location in the policy
// It returns true if the provider's Location field matches the policy's RequiredLocation
func (f FeatureFilter) ApplySingle(provider *pairing.Provider, policy *pairing.ConsumerPolicy) bool {
	// Use a map for efficient lookup of required features
	required := make(map[string]bool)
	for _, feature := range policy.RequiredFeatures {
		required[feature] = true
	}

	// Check if the provider has all required features
	for req := range required {
		found := false
		for _, pf := range provider.Features {
			if pf == req {
				found = true
				break // Found this required feature, check the next one
			}
		}
		if !found {
			return false // A required feature wasn't found in the provider's list
		}
	}
	return true // All required features were found
}

func (f FeatureFilter) Name() string { return "FeatureFilter" }

/* ***********************************************************************
 *                            STAKE FILTER                               *
 *********************************************************************** */

// Apply filters providers based on the minimum stake requirement in the policy
// It retains only those providers whose Stake field is greater than or equal to the policy's MinStake
func (f StakeFilter) Apply(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) []*pairing.Provider {
	var result []*pairing.Provider
	for _, p := range providers {
		if p.Stake >= policy.MinStake {
			result = append(result, p)
		}
	}
	return result
}

// ApplySingle checks if a single provider meets the minimum stake requirement in the policy
// It returns true if the provider's Stake field is greater than or equal to the policy's MinStake
func (f StakeFilter) ApplySingle(provider *pairing.Provider, policy *pairing.ConsumerPolicy) bool {
	return provider.Stake >= policy.MinStake
}

func (f StakeFilter) Name() string { return "StakeFilter" }
