// Package xi implements the Ξ (Xi) Advisory Quality Function.
//
// The Ξ function measures advisory quality through three orthogonal dimensions:
// 1. Relevance (R): Do engineers apply the recommended patterns?
// 2. Actionability (A): Do generated artifacts require minimal modification?
// 3. Impact Realization (I): Do deployed services achieve their performance targets?
//
// The overall Ξ score is the geometric mean: Ξ = (R × A × I)^(1/3)
//
// This geometric mean ensures that excellent performance on one dimension
// cannot compensate for failure on another—all three must be high for
// high overall quality.
package xi

import (
	"fmt"
	"math"
)

// Score represents a Ξ quality score with its component dimensions.
type Score struct {
	Relevance     float64 // R: pattern application rate
	Actionability float64 // A: code modification rate
	Impact        float64 // I: performance target achievement
	Overall       float64 // Ξ: geometric mean of R, A, I
}

// Calculate computes the overall Ξ score from its components.
//
// Formula: Ξ = (R × A × I)^(1/3)
//
// Each component must be in the range [0.0, 1.0].
func Calculate(relevance, actionability, impact float64) (*Score, error) {
	// Validate inputs
	if err := validateComponent("Relevance", relevance); err != nil {
		return nil, err
	}
	if err := validateComponent("Actionability", actionability); err != nil {
		return nil, err
	}
	if err := validateComponent("Impact", impact); err != nil {
		return nil, err
	}

	// Compute geometric mean
	product := relevance * actionability * impact
	overall := math.Pow(product, 1.0/3.0)

	return &Score{
		Relevance:     relevance,
		Actionability: actionability,
		Impact:        impact,
		Overall:       overall,
	}, nil
}

// CalculateRelevance computes the Relevance (R) component.
//
// Formula: R = (patterns_applied / patterns_recommended) × confidence_weight
//
// Where:
//   patterns_applied = number of recommended patterns actually implemented
//   patterns_recommended = number of patterns suggested by OSE
//   confidence_weight = average confidence score of applied patterns
func CalculateRelevance(patternsApplied, patternsRecommended int, avgConfidence float64) (float64, error) {
	if patternsRecommended == 0 {
		// No patterns recommended—relevance is perfect by definition
		return 1.0, nil
	}

	if patternsApplied < 0 || patternsApplied > patternsRecommended {
		return 0, fmt.Errorf("invalid pattern counts: applied=%d, recommended=%d",
			patternsApplied, patternsRecommended)
	}

	if avgConfidence < 0.0 || avgConfidence > 1.0 {
		return 0, fmt.Errorf("invalid confidence: %f (must be in [0,1])", avgConfidence)
	}

	applicationRate := float64(patternsApplied) / float64(patternsRecommended)
	relevance := applicationRate * avgConfidence

	return relevance, nil
}

// CalculateActionability computes the Actionability (A) component.
//
// Formula: A = 1 - sqrt(lines_modified / lines_generated)
//
// The square root dampens the penalty—small modifications are expected
// (domain-specific business logic), but extensive rewrites indicate
// templates need improvement.
func CalculateActionability(linesModified, linesGenerated int) (float64, error) {
	if linesGenerated == 0 {
		// No code generated—actionability is perfect by definition
		return 1.0, nil
	}

	if linesModified < 0 {
		return 0, fmt.Errorf("invalid lines_modified: %d (must be >= 0)", linesModified)
	}

	if linesGenerated < 0 {
		return 0, fmt.Errorf("invalid lines_generated: %d (must be >= 0)", linesGenerated)
	}

	modificationRate := float64(linesModified) / float64(linesGenerated)

	// Dampen with square root
	dampedRate := math.Sqrt(modificationRate)

	// A = 1 - dampedRate, bounded to [0, 1]
	actionability := math.Max(0.0, 1.0-dampedRate)

	return actionability, nil
}

// CalculateImpact computes the Impact Realization (I) component.
//
// Formula: I = (met_targets / total_targets) × stability_factor
//
// Where:
//   met_targets = number of performance targets achieved
//   total_targets = number of performance targets specified
//   stability_factor = 1 - (incident_count × severity_weight)
//
// Severity weights:
//   SEV1 (Critical): 0.5
//   SEV2 (Major):    0.2
//   SEV3 (Minor):    0.1
func CalculateImpact(metTargets, totalTargets int, incidents []Incident) (float64, error) {
	if totalTargets == 0 {
		return 0, fmt.Errorf("total_targets must be > 0")
	}

	if metTargets < 0 || metTargets > totalTargets {
		return 0, fmt.Errorf("invalid target counts: met=%d, total=%d",
			metTargets, totalTargets)
	}

	// Calculate target achievement rate
	targetScore := float64(metTargets) / float64(totalTargets)

	// Calculate stability factor from incidents
	stabilityFactor := 1.0
	for _, incident := range incidents {
		weight := incidentSeverityWeight(incident.Severity)
		stabilityFactor -= weight
	}

	// Ensure stability factor doesn't go negative
	if stabilityFactor < 0 {
		stabilityFactor = 0
	}

	impact := targetScore * stabilityFactor

	return impact, nil
}

// Incident represents a production incident that affects Impact calculation.
type Incident struct {
	Severity         string // "SEV1", "SEV2", "SEV3"
	Description      string
	DurationMinutes  int
}

// incidentSeverityWeight returns the penalty for an incident based on severity.
func incidentSeverityWeight(severity string) float64 {
	switch severity {
	case "SEV1":
		return 0.5 // Critical outage
	case "SEV2":
		return 0.2 // Major issue
	case "SEV3":
		return 0.1 // Minor issue
	default:
		return 0.05 // Unknown severity - minimal penalty
	}
}

// validateComponent ensures a Ξ component is in valid range [0.0, 1.0].
func validateComponent(name string, value float64) error {
	if value < 0.0 || value > 1.0 {
		return fmt.Errorf("%s must be in [0.0, 1.0], got %f", name, value)
	}
	if math.IsNaN(value) {
		return fmt.Errorf("%s is NaN", name)
	}
	if math.IsInf(value, 0) {
		return fmt.Errorf("%s is infinite", name)
	}
	return nil
}

// QualityLevel returns a human-readable quality level for a Ξ score.
func QualityLevel(xi float64) string {
	switch {
	case xi >= 0.85:
		return "EXCELLENT"
	case xi >= 0.75:
		return "VERY GOOD"
	case xi >= 0.65:
		return "GOOD"
	case xi >= 0.50:
		return "ACCEPTABLE"
	case xi >= 0.35:
		return "POOR"
	default:
		return "CRITICAL"
	}
}

// String returns a formatted string representation of the Score.
func (s *Score) String() string {
	return fmt.Sprintf("Ξ=%.3f (%s) [R=%.3f, A=%.3f, I=%.3f]",
		s.Overall, QualityLevel(s.Overall),
		s.Relevance, s.Actionability, s.Impact)
}
