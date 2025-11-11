package xi

import (
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name          string
		relevance     float64
		actionability float64
		impact        float64
		wantOverall   float64
		wantErr       bool
	}{
		{
			name:          "perfect score",
			relevance:     1.0,
			actionability: 1.0,
			impact:        1.0,
			wantOverall:   1.0,
			wantErr:       false,
		},
		{
			name:          "excellent scores (Sarah's example)",
			relevance:     0.563,
			actionability: 0.621,
			impact:        0.90,
			wantOverall:   0.680, // (0.563 × 0.621 × 0.90)^(1/3)
			wantErr:       false,
		},
		{
			name:          "zero score",
			relevance:     0.0,
			actionability: 0.0,
			impact:        0.0,
			wantOverall:   0.0,
			wantErr:       false,
		},
		{
			name:          "invalid relevance (too high)",
			relevance:     1.5,
			actionability: 1.0,
			impact:        1.0,
			wantErr:       true,
		},
		{
			name:          "invalid impact (negative)",
			relevance:     1.0,
			actionability: 1.0,
			impact:        -0.1,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score, err := Calculate(tt.relevance, tt.actionability, tt.impact)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(score.Overall-tt.wantOverall) > 0.01 {
				t.Errorf("overall score = %.3f, want %.3f", score.Overall, tt.wantOverall)
			}

			// Verify geometric mean property
			expectedOverall := math.Pow(tt.relevance*tt.actionability*tt.impact, 1.0/3.0)
			if math.Abs(score.Overall-expectedOverall) > 0.001 {
				t.Errorf("geometric mean calculation incorrect: got %.3f, want %.3f",
					score.Overall, expectedOverall)
			}
		})
	}
}

func TestCalculateRelevance(t *testing.T) {
	tests := []struct {
		name               string
		patternsApplied    int
		patternsRecommended int
		avgConfidence      float64
		want               float64
		wantErr            bool
	}{
		{
			name:                "all patterns applied (Sarah's example)",
			patternsApplied:     2,
			patternsRecommended: 3,
			avgConfidence:       0.845,
			want:                0.563, // (2/3) × 0.845
			wantErr:             false,
		},
		{
			name:                "perfect application",
			patternsApplied:     5,
			patternsRecommended: 5,
			avgConfidence:       1.0,
			want:                1.0,
			wantErr:             false,
		},
		{
			name:                "no patterns applied",
			patternsApplied:     0,
			patternsRecommended: 3,
			avgConfidence:       0.9,
			want:                0.0,
			wantErr:             false,
		},
		{
			name:                "no patterns recommended",
			patternsApplied:     0,
			patternsRecommended: 0,
			avgConfidence:       1.0,
			want:                1.0, // Perfect by definition
			wantErr:             false,
		},
		{
			name:                "invalid: more applied than recommended",
			patternsApplied:     5,
			patternsRecommended: 3,
			avgConfidence:       0.9,
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateRelevance(tt.patternsApplied, tt.patternsRecommended, tt.avgConfidence)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("relevance = %.3f, want %.3f", got, tt.want)
			}
		})
	}
}

func TestCalculateActionability(t *testing.T) {
	tests := []struct {
		name           string
		linesModified  int
		linesGenerated int
		want           float64
		wantErr        bool
	}{
		{
			name:           "Sarah's example",
			linesModified:  180,
			linesGenerated: 1247,
			want:           0.621, // 1 - sqrt(180/1247) = 1 - sqrt(0.144) = 1 - 0.379
			wantErr:        false,
		},
		{
			name:           "no modifications",
			linesModified:  0,
			linesGenerated: 1000,
			want:           1.0,
			wantErr:        false,
		},
		{
			name:           "complete rewrite",
			linesModified:  1000,
			linesGenerated: 1000,
			want:           0.0, // 1 - sqrt(1.0) = 0
			wantErr:        false,
		},
		{
			name:           "small modifications (10%)",
			linesModified:  100,
			linesGenerated: 1000,
			want:           0.684, // 1 - sqrt(0.1) = 1 - 0.316
			wantErr:        false,
		},
		{
			name:           "no code generated",
			linesModified:  0,
			linesGenerated: 0,
			want:           1.0, // Perfect by definition
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateActionability(tt.linesModified, tt.linesGenerated)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("actionability = %.3f, want %.3f", got, tt.want)
			}
		})
	}
}

func TestCalculateImpact(t *testing.T) {
	tests := []struct {
		name         string
		metTargets   int
		totalTargets int
		incidents    []Incident
		want         float64
		wantErr      bool
	}{
		{
			name:         "Sarah's example (perfect targets, one SEV3)",
			metTargets:   4,
			totalTargets: 4,
			incidents: []Incident{
				{Severity: "SEV3", Description: "Cache invalidation bug"},
			},
			want:    0.90, // (4/4) × (1 - 0.1) = 1.0 × 0.9
			wantErr: false,
		},
		{
			name:         "perfect (no incidents)",
			metTargets:   4,
			totalTargets: 4,
			incidents:    []Incident{},
			want:         1.0,
			wantErr:      false,
		},
		{
			name:         "partial targets met",
			metTargets:   3,
			totalTargets: 4,
			incidents:    []Incident{},
			want:         0.75, // (3/4) × 1.0
			wantErr:      false,
		},
		{
			name:         "SEV1 incident impact",
			metTargets:   4,
			totalTargets: 4,
			incidents: []Incident{
				{Severity: "SEV1", Description: "Complete outage"},
			},
			want:    0.50, // (4/4) × (1 - 0.5)
			wantErr: false,
		},
		{
			name:         "multiple incidents",
			metTargets:   4,
			totalTargets: 4,
			incidents: []Incident{
				{Severity: "SEV2"},
				{Severity: "SEV3"},
				{Severity: "SEV3"},
			},
			want:    0.60, // (4/4) × (1 - 0.2 - 0.1 - 0.1)
			wantErr: false,
		},
		{
			name:         "incidents exceed stability (floor at 0)",
			metTargets:   4,
			totalTargets: 4,
			incidents: []Incident{
				{Severity: "SEV1"},
				{Severity: "SEV1"},
				{Severity: "SEV2"},
			},
			want:    0.0, // (4/4) × max(0, 1 - 0.5 - 0.5 - 0.2) = 0
			wantErr: false,
		},
		{
			name:         "invalid: zero targets",
			metTargets:   0,
			totalTargets: 0,
			incidents:    []Incident{},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateImpact(tt.metTargets, tt.totalTargets, tt.incidents)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(got-tt.want) > 0.01 {
				t.Errorf("impact = %.3f, want %.3f", got, tt.want)
			}
		})
	}
}

func TestQualityLevel(t *testing.T) {
	tests := []struct {
		xi   float64
		want string
	}{
		{0.95, "EXCELLENT"},
		{0.85, "EXCELLENT"},
		{0.80, "VERY GOOD"},
		{0.75, "VERY GOOD"},
		{0.70, "GOOD"},
		{0.65, "GOOD"},
		{0.60, "ACCEPTABLE"},
		{0.50, "ACCEPTABLE"},
		{0.40, "POOR"},
		{0.35, "POOR"},
		{0.20, "CRITICAL"},
		{0.00, "CRITICAL"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := QualityLevel(tt.xi)
			if got != tt.want {
				t.Errorf("QualityLevel(%.2f) = %s, want %s", tt.xi, got, tt.want)
			}
		})
	}
}

func TestScoreString(t *testing.T) {
	score := &Score{
		Relevance:     0.563,
		Actionability: 0.621,
		Impact:        0.90,
		Overall:       0.680,
	}

	got := score.String()
	want := "Ξ=0.680 (GOOD) [R=0.563, A=0.621, I=0.900]"

	if got != want {
		t.Errorf("String() = %s, want %s", got, want)
	}
}

// Benchmark the geometric mean calculation
func BenchmarkCalculate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Calculate(0.8, 0.7, 0.9)
	}
}
