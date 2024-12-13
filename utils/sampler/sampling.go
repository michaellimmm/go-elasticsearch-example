package sampler

import (
	"math"
)

func CalculateTotalSampleSize(totalDatasetSize int, confidenceLevel float64, marginOfError float64) int {
	// Z-scores for common confidence levels
	zScores := map[float64]float64{
		0.90: 1.645,
		0.95: 1.96,
		0.99: 2.576,
	}

	// Find the closest confidence level if exact match not found
	var zScore float64
	minDiff := math.Inf(1)
	for level, score := range zScores {
		diff := math.Abs(level - confidenceLevel)
		if diff < minDiff {
			minDiff = diff
			zScore = score
		}
	}

	// Simplified sample size calculation formula
	// n = (Z^2 * p * (1-p)) / E^2
	// Z = z-score, p = expected proportion (0.5 for max variability), E = margin of error
	p := 0.5 // Use 0.5 for maximum variability
	numerator := math.Pow(zScore, 2) * p * (1 - p)
	denominator := math.Pow(marginOfError, 2)

	sampleSize := int(math.Ceil(numerator / denominator))

	// Adjust if sample size is larger than population
	if sampleSize > totalDatasetSize {
		return totalDatasetSize
	}

	return sampleSize
}
