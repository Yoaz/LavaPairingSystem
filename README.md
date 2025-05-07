# Provider Pairing System

## Overview

This project implements the core filtering and scoring mechanisms for Lava Network’s provider pairing system. It processes a list of RPC service providers against a consumer's policy, filters the valid providers, scores them, ranks them, and returns the top 5 matches.

The system is designed for correctness, efficiency, and concurrency, adhering to best practices in Go development.

## Features

✅ **Filtering:**

- `LocationFilter`: Keeps providers matching the required location.
- `FeatureFilter`: Keeps providers supporting all required features.
- `StakeFilter`: Keeps providers meeting the minimum stake.

✅ **Scoring:**

- `StakeScore`: Higher score for higher stake (normalized).
- `FeatureScore`: Higher score for extra features beyond the minimum.
- `LocationScore`: Perfect score if matching location, lower otherwise.
- `FeeScore`: Adds an additional scoring strategy based on provider fees, normalized.

✅ **Weighted Scoring:**

- The system supports both **equal-weight averaging** (default when no weights are supplied) and **custom weighted scoring**.
- To use weighted scoring, the `ConsumerPolicy.Weights` map should include weights (as float64) for each scorer (e.g., `"StakeScore"`, `"FeatureScore"`).
- The sum of provided weights must equal **1.0**.
- If the `Weights` map is `nil` or empty, the system falls back to equal averaging: `finalScore = totalScore / numberOfScorers`.
- If weights are provided, only scorers present in the map contribute (missing scorers treated as zero).

## Project Structure

```
/           → Root
cmd/
  main.go                  → Entry point
config/
  config.go               → Configuration construction
internal/
  filter/                 → Filtering logic (e.g., by location, stake, features)
    filter.go
    types.go
  score/                  → Scoring logic (e.g., stake score, feature score, fee score)
    score.go
    types.go
  system/                 → Core system orchestration
    system.go
  models.go               → Shared models (Provider, ConsumerPolicy, PairingScore)
  logger/
    logger.go             → Custom slog-based logger
  utils/
    utils.go              → Utilities logic
```

## Usage

Make sure you have Go installed.

### Build and Run

```
make            # Builds and runs the app
make build      # Only builds the binary
make run        # Runs the app without rebuilding
make clean      # Cleans the binary directory
```

### Example One Liner Build & Run

```
make
```

This will execute the pairing system against a sample list of providers and a sample policy.

## Weighted Scoring Input Example

Example `ConsumerPolicy.Weights` map:

```go
Weights: map[string]float64{
    "StakeScore":    0.5,
    "FeatureScore":  0.3,
    "LocationScore": 0.2,
}
```

- If **no weights** are supplied (`nil` or empty map), the system averages all scorer results equally.
- If **partial weights** are supplied (e.g., only StakeScore), only those scorers contribute, and others are treated as zero.
- The `ValidateWeights` function ensures that provided weights sum to exactly 1.0 when present.

## Design Rationale

- **Separation of concerns:** Filters and scorers are separate for clarity and future extensibility.
- **Concurrency:** Ranking uses goroutines to parallelize score calculation for performance.
- **Normalization:** All scores are scaled between 0 and 1, allowing fair combination of diverse criteria.
- **Fallback:** If weights are missing or invalid, the system gracefully falls back to equal-weight averaging.
- **Error handling:** Clean and minimal, leaving room for expansion in production-ready systems.

## Architecture Diagram

```
+-----------------------+
|    PairingSystem      |
|-----------------------|
| - filters: []Filter   |
| - scorers: []Scorer   |
+-----------------------+
        |       |
        |       |
        v       v
+-------------+ +-------------+
|   Filter    | |   Scorer    |
|  Interface  | |  Interface  |
+-------------+ +-------------+
        |               |
        v               v
+----------------+ +----------------+
| SpecificFilter | | SpecificScorer |
+----------------+ +----------------+
        |               |
        v               v
+----------------+ +----------------+
| Provider       | | ConsumerPolicy |
+----------------+ +----------------+
        |               |
        v               v
+-----------------------------+
|     PairingScore            |
+-----------------------------+
```

## Flow Chart

```
[Start]
   |
   v
[Load Providers + Policy]
   |
   v
[Filter Providers] --(if none)--> [strictMode == true] --> [Return Error]
                      |                         |
                      |                         v
                      |                    [strictMode == false]
                      v                         |
             [Rank Providers (Concurrent Scoring)]
                      |
                      v
              [Sort by Final Score]
                      |
                      v
              [Select Top 5 Providers]
                      |
                      v
              [Return Pairing List]
                      |
                      v
                    [End]
```

## Requirements

- Go 1.20+

## Notes

This is a simplified implementation for the challenge and can be extended with:

- More sophisticated scoring mechanism.
- Realistic provider data sources.
- Integration with blockchain query interfaces.
- Integrating config.yaml for configuration setup

## Author

Yoaz Sh. | [YӨΛZ](https://yoaz.info)
