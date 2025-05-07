package score

import (
	"strings"

	pairing "github.com/Yoaz/LavaPairingSystem/internal"
)

/* ***********************************************************************
 *                            STAKE SCORE                                *
 *********************************************************************** */

// Score calculates a normalized score based on the provider's stake relative to the maximum stake
// observed in the currently considered provider pool
// The maxStake value is in the PreScoreContext, which is passed to the Score method
// This allows the score to be calculated dynamically based on the current pool of providers
func (s *StakeScore) Score(p *pairing.Provider, _ *pairing.ConsumerPolicy, ctx *PreScoreContext) float64 {
	// Prevent division by zero if maxStake hasn't been set or is zero
	if ctx.MaxStake == 0 {
		return 0.0
	}
	// Normalize stake: provider's stake / maximum stake in the pool
	return float64(p.Stake) / float64(ctx.MaxStake)
}

func (s *StakeScore) Name() string { return "StakeScore" }

/* ***********************************************************************
 *                            FEATURE SCORE                              *
 *********************************************************************** */

// Score calculates a score based on the number of extra features the provider offers beyond
// those required by the policy, normalized by the total number of features the provider has
func (s *FeatureScore) Score(p *pairing.Provider, policy *pairing.ConsumerPolicy, ctx *PreScoreContext) float64 {
	// Prevent division by zero if the provider has no features
	if len(p.Features) == 0 {
		return 0.0
	}

	extra := 0
	required := make(map[string]bool)
	for _, feat := range policy.RequiredFeatures {
		required[feat] = true
	}
	for _, pf := range p.Features {
		if !required[pf] {
			extra++
		}
	}
	// Normalize score: number of extra features / total number of features
	// This gives a score between 0 and 1, rewarding providers with a higher proportion of extra features
	return float64(extra) / float64(len(p.Features))
}

func (s *FeatureScore) Name() string { return "FeatureScore" }

/* ***********************************************************************
 *                            LOCATION SCORE                             *
 *********************************************************************** */

// Score assigns a perfect score (1.0) if the provider's location matches the required location (case-insensitive),
// and a lower, fixed score (0.5) otherwise
func (s *LocationScore) Score(p *pairing.Provider, policy *pairing.ConsumerPolicy, ctx *PreScoreContext) float64 {
	if strings.EqualFold(p.Location, policy.RequiredLocation) {
		return 1.0
	}
	// Assign an arbitrary lower score for non-matching locations
	// NOTE: A more sophisticated approach might consider geographic proximity or other factors
	return 0.5
}

func (s *LocationScore) Name() string { return "LocationScore" }

/* ***********************************************************************
 *                            FEE SCORE                                  *
 *********************************************************************** */

// Score calculates a score based on the provider's fee, normalized to a range of 0 to 1
func (s *FeeScore) Score(provider *pairing.Provider, policy *pairing.ConsumerPolicy, ctx *PreScoreContext) float64 {
	fee, ok := ctx.NormalizedFees[provider.ID]
	if !ok {
		return 0
	}
	return 1.0 - fee // Lower fee is better
}

func (s *FeeScore) Name() string { return "FeeScore" }
