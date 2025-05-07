package system

import (
	"fmt"
	"io"
	"log/slog"
	"sort"
	"sync"

	pairing "github.com/Yoaz/LavaPairingSystem/internal"
	"github.com/Yoaz/LavaPairingSystem/internal/filter"
	"github.com/Yoaz/LavaPairingSystem/internal/score"
	"github.com/Yoaz/LavaPairingSystem/internal/utils"
)

// NewPairingSystem creates a new PairingSystem instance with the provided filters, scorers, and logger
// StrictMode determines if the system should return an error when no providers match the filter criteria
func NewPairingSystem(filters []filter.Filter, scorers []score.Scorer, logger *slog.Logger, strictMode bool) PairingSystem {
	// Ensure logger is not nil, provide a default discard logger if it is
	if logger == nil {
		// If no logger is provided, default to discarding logs
		logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return &pairingSystem{
		filters:    filters,
		scorers:    scorers,
		logger:     logger,
		strictMode: strictMode, // NOTE: If true, returns error when no providers match; if false, returns empty list
	}
}

/* ***********************************************************************
 *                                   CORE                                *
 *********************************************************************** */

// FilterProviders filters the list of providers based on the consumer policy
// It applies each filter in the order they were added to the PairingSystem
func (ps *pairingSystem) FilterProviders(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) []*pairing.Provider {
	ps.logger.Debug("Starting provider filtering", "initial_count", len(providers))

	// Check if there are any providers to filter
	if len(providers) == 0 {
		return []*pairing.Provider{}
	}

	// Sequential filtering for small lists
	if len(providers) <= parallelFilterThreshold {
		filtered := providers
		for _, filter := range ps.filters {
			countBefore := len(filtered)
			filtered = filter.Apply(filtered, policy)
			countAfter := len(filtered)
			ps.logger.Debug("Filter applied", "filter_name", filter.Name(), "count_before", countBefore, "count_after", countAfter)
		}
		ps.logger.Debug("Finished sequential provider filtering", "final_count", len(filtered))
		return filtered
	}

	// Parallel filtering for large lists
	filtered := ps.parallelFilterProviders(providers, policy)
	ps.logger.Debug("Finished parallel provider filtering", "final_count", len(filtered))
	return filtered
}

// parallelFilterProviders filters providers in parallel using goroutines
// It creates a worker pool to process the providers concurrently
// Each worker applies the filters to a provider and sends the result to a results channel
func (ps *pairingSystem) parallelFilterProviders(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) []*pairing.Provider {
	tasks := make(chan *pairing.Provider, len(providers))
	results := make(chan *pairing.Provider, len(providers))

	var wg sync.WaitGroup

	// Start workers
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go ps.filterWorker(w, tasks, results, policy, &wg)
	}

	// Feed tasks
	for _, p := range providers {
		tasks <- p
	}
	close(tasks)

	// Close results channel once workers are done
	// Block until all workers finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	filtered := make([]*pairing.Provider, 0, len(providers))
	for p := range results {
		filtered = append(filtered, p)
	}

	return filtered
}

// RankProviders ranks the filtered providers based on the consumer policy
// It calculates scores using the provided scorers and returns a list of PairingScore
// Each PairingScore contains the provider, its score, and the individual components of the score
//
// NOTE: If weights are provided in the policy, they are used to calculate a weighted score
// If no weights are provided, the average score is used
func (ps *pairingSystem) RankProviders(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) []*pairing.PairingScore {
	ps.logger.Debug("Starting provider ranking", "provider_count", len(providers))

	if len(providers) == 0 {
		ps.logger.Debug("No providers to rank, returning empty list.")
		return []*pairing.PairingScore{}
	}

	// Compute max stake for normalization
	// This is done to ensure that the stake scores are relative to the maximum stake in the list
	currentMaxStake := utils.ComputeMaxStake(providers)
	if currentMaxStake == 0 {
		ps.logger.Debug("No providers with stake found, setting max stake to 1")
		currentMaxStake = 1
	} else {
		ps.logger.Debug("Computed max stake for normalization", "max_stake", currentMaxStake)
	}

	// Compute normalized fees for providers
	// This is done to ensure that the fee scores are relative to the maximum fee in the list
	normalizedFees := utils.ComputeNormalizedFees(providers)

	preScoreCtx := &score.PreScoreContext{
		MaxStake:       currentMaxStake,
		NormalizedFees: normalizedFees,
	}

	tasks := make(chan *pairing.Provider, len(providers))
	results := make(chan *pairing.PairingScore, len(providers))

	var wg sync.WaitGroup

	// Start worker goroutines
	for w := 0; w < workerCount; w++ {
		wg.Add(1)
		go ps.rankWorker(w, tasks, results, policy, preScoreCtx, &wg)
	}

	// Feed tasks
	for _, provider := range providers {
		tasks <- provider
	}
	close(tasks)

	// Wait for workers to finish and close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	scores := make([]*pairing.PairingScore, 0, len(providers))
	for score := range results {
		scores = append(scores, score)
	}

	ps.logger.Debug("Finished calculating all provider scores")
	return scores
}

// GetPairingList retrieves a list of top providers based on the consumer policy
// It filters, ranks, and sorts the providers, returning the top N providers
func (ps *pairingSystem) GetPairingList(providers []*pairing.Provider, policy *pairing.ConsumerPolicy) ([]*pairing.Provider, error) {
	ps.logger.Info("Starting GetPairingList", "initial_provider_count", len(providers))

	// Step 1: Filter providers based on policy requirements
	filtered := ps.FilterProviders(providers, policy)
	if len(filtered) == 0 {
		ps.logger.Warn("No providers matched the filter criteria.")

		if ps.strictMode {
			return nil, fmt.Errorf("strict mode: no providers matched the filter criteria")
		}

		return []*pairing.Provider{}, nil // Graceful: return empty list, no error
	}
	ps.logger.Debug("Filtering complete", "filtered_count", len(filtered))

	// Step 2: Rank the filtered providers based on scoring criteria
	scored := ps.RankProviders(filtered, policy)
	ps.logger.Debug("Ranking complete", "ranked_count", len(scored))

	// Step 3: Sort providers by their final score in descending order
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score // Higher score first
	})
	ps.logger.Debug("Sorting complete")

	// Step 4: Select the top N providers
	finalCount := utils.Min(topNProviders, len(scored)) // Handle fewer providers than topN
	topProviders := make([]*pairing.Provider, 0, finalCount)
	for i := 0; i < finalCount; i++ {
		topProviders = append(topProviders, scored[i].Provider)
		ps.logger.Debug("Selected provider",
			"rank", i+1,
			"address", scored[i].Provider.Address,
			"score", scored[i].Score,
			"components", scored[i].Components,
		)
	}

	ps.logger.Info("Finished GetPairingList", "selected_count", len(topProviders))
	return topProviders, nil
}

/* ***********************************************************************
 *                                   WORKERS                             *
 *********************************************************************** */

// rankWorker is a goroutine that processes providers and calculates their scores
// It takes a provider from the tasks channel, scores it using the provided scorers,
// and sends the result to the results channel
func (ps *pairingSystem) rankWorker(workerID int, tasks <-chan *pairing.Provider, results chan<- *pairing.PairingScore, policy *pairing.ConsumerPolicy, preScoreCtx *score.PreScoreContext, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range tasks {
		components := make(map[string]float64)
		var totalScore float64

		for _, scorer := range ps.scorers {
			s := scorer.Score(p, policy, preScoreCtx)
			components[scorer.Name()] = s
			totalScore += s
		}

		finalScore := 0.0
		// Check if weighted scoring should be applied
		// NOTE: Defined in struct as a map[string]float64 therefore no need to check for nil
		if len(policy.Weights) > 0 {
			ps.logger.Debug("Applying weighted scoring logic", "worker_id", workerID, "provider_id", p.ID)
			var weightedSum float64
			// The validation in main.go ensures that if policy.Weights is present, its values sum to 1.
			// Iterating through the components we calculated.
			// If a components's (scorer's) name is in policy.Weights, its score is weighted.
			// If not, its effective weight is 0 for this weighted sum.
			for name, scoreValue := range components {
				weight, ok := policy.Weights[name]
				if ok {
					weightedSum += scoreValue * weight
				} else {
					// If a scorer is not in the weights map, it contributes 0 to the weighted score.
					// This implies the user intentionally omitted it from the weighted scheme.
					ps.logger.Debug("Scorer not found in policy weights, applying 0 weight", "worker_id", workerID, "provider_id", p.ID, "scorer_name", name)
				}
			}
			finalScore = weightedSum
		} else {
			// Fallback to average scoring if weights are not provided
			ps.logger.Debug("Applying average (equal weight) scoring logic", "worker_id", workerID, "provider_id", p.ID)
			if len(ps.scorers) > 0 {
				finalScore = totalScore / float64(len(ps.scorers))
			}
		}

		results <- &pairing.PairingScore{
			Provider:   p,
			Score:      finalScore,
			Components: components,
		}

		ps.logger.Debug("Rank-Worker scored provider",
			"worker_id", workerID,
			"provider_id", p.ID,
			"score", finalScore,
			"components", components,
		)
	}
}

// filterWorker is a goroutine that processes providers and applies filters to them
func (ps *pairingSystem) filterWorker(workerID int, tasks <-chan *pairing.Provider, results chan<- *pairing.Provider, policy *pairing.ConsumerPolicy, wg *sync.WaitGroup) {
	defer wg.Done()

	for p := range tasks {
		pass := true
		for _, filter := range ps.filters {
			// Apply the filter to the provider
			if !filter.ApplySingle(p, policy) {
				ps.logger.Debug("Filter-Worker filter rejected provider",
					"worker_id", workerID,
					"provider_id", p.ID,
					"filter_name", filter.Name(),
				)
				// If the provider doesn't pass the filter, break out of the loop
				pass = false
				break
			}
		}
		// If the provider passes all filters, send it to the results channel
		if pass {
			results <- p
		}
	}
}
