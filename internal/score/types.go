package score

import pairing "github.com/Yoaz/LavaPairingSystem/internal"

// Scorer is an interface for scoring providers based on a consumer policy
type Scorer interface {
	Score(provider *pairing.Provider, policy *pairing.ConsumerPolicy, ctx *PreScoreContext) float64
	Name() string
}

type (
	StakeScore    struct{}
	FeatureScore  struct{}
	LocationScore struct{}
	FeeScore      struct{}
)

// PreScoreContext holds the context for pre-scoring calculations
type PreScoreContext struct {
	MaxStake       int64
	AverageLatency float64
	NormalizedFees map[string]float64
}
