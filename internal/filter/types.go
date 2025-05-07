package filter

import pairing "github.com/Yoaz/LavaPairingSystem/internal"

// Filter is an interface for filtering providers based on a consumer policy
type Filter interface {
	Apply(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) []*pairing.Provider
	ApplySingle(provider *pairing.Provider, policy *pairing.ConsumerPolicy) bool
	Name() string // for tracking filter name
}

// Filter implementations for different criteria
type (
	LocationFilter struct{} // Filters providers based on location
	FeatureFilter  struct{} // Filters providers based on features
	StakeFilter    struct{} // Filters providers based on stake
)
