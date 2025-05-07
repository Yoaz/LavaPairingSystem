package pairing

// Provider represents a provider in the pairing system.
type Provider struct {
	ID       string  // Unique identifier for the provider (--> NOTE: ADDED TO GIVE AN EXAMPLE FOR ANOTHER SCORE TYPE)
	Fee      float64 // Fee charged by the provider (--> NOTE: ADDED TO GIVE AN EXAMPLE FOR ANOTHER SCORE TYPE)
	Address  string
	Stake    int64
	Location string
	Features []string
}

// ConsumerPolicy represents the policy requirements for a consumer
type ConsumerPolicy struct {
	RequiredLocation string
	RequiredFeatures []string
	MinStake         int64
	// Weights for different scoring components (e.g., {"Stake": 0.5, "Location": 0.3, "Feature": 0.2})
	// This allows for flexible scoring based on the consumer's preferences.
	// NOTE: Th weights should sum to 1.0
	Weights map[string]float64 // (--> NOTE: ADDED TO GIVE AN EXAMPLE FOR WEIGHTED SCORING MECHANISM)
}

// PairingScore represents the score of a provider based on the consumer policy
type PairingScore struct {
	Provider   *Provider
	Score      float64
	Components map[string]float64 // (e.g., {"StakeScore": 0.8, "FeatureScore": 1.0}
}
