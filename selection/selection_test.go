package selection_test

import (
	"math/rand/v2"
	"testing"

	"github.com/LCmaster/go-gap-engine/selection"
)

func TestTournament(t *testing.T) {
	tests := []struct {
		name           string
		tournamentSize int
		numSelect      int
		wantLen        int
	}{
		{"Select 2, Size 2", 2, 2, 2},
		{"Select 5, Size 1", 1, 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := rand.New(rand.NewPCG(1, 2))
			pop := []string{"A", "B", "C", "D"}
			fits := []float64{1.0, 10.0, 2.0, 0.5}

			sel := selection.Tournament[string](tt.tournamentSize)
			selected := sel(rng, pop, fits, tt.numSelect)

			if len(selected) != tt.wantLen {
				t.Errorf("Expected %v selected, got %v", tt.wantLen, len(selected))
			}
		})
	}
}

func TestRouletteWheel(t *testing.T) {
	tests := []struct {
		name      string
		numSelect int
		wantLen   int
		check     func(t *testing.T, selected []string)
	}{
		{
			name:      "Select 100",
			numSelect: 100,
			wantLen:   100,
			check: func(t *testing.T, selected []string) {
				countB := 0
				for _, s := range selected {
					if s == "B" {
						countB++
					}
				}
				if countB < 60 || countB > 100 {
					t.Errorf("Expected B to be selected frequently, got %v times out of 100", countB)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rng := rand.New(rand.NewPCG(1, 2))
			pop := []string{"A", "B"}
			fits := []float64{1.0, 9.0}

			sel := selection.RouletteWheel[string]()
			selected := sel(rng, pop, fits, tt.numSelect)

			if len(selected) != tt.wantLen {
				t.Errorf("Expected %v selected, got %v", tt.wantLen, len(selected))
			}

			if tt.check != nil {
				tt.check(t, selected)
			}
		})
	}
}
