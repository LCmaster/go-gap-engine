package selection_test

import (
	"testing"
	"go-gap-engine/selection"
)

func TestTournament(t *testing.T) {
	pop := []string{"A", "B", "C", "D"}
	fits := []float64{1.0, 10.0, 2.0, 0.5}

	sel := selection.Tournament[string](2)
	selected := sel(pop, fits, 2)

	if len(selected) != 2 {
		t.Errorf("Expected 2 selected, got %v", len(selected))
	}
}

func TestRouletteWheel(t *testing.T) {
	pop := []string{"A", "B"}
	fits := []float64{1.0, 9.0}

	sel := selection.RouletteWheel[string]()
	selected := sel(pop, fits, 100)

	if len(selected) != 100 {
		t.Errorf("Expected 100 selected, got %v", len(selected))
	}

	countB := 0
	for _, s := range selected {
		if s == "B" {
			countB++
		}
	}

	// B should be selected roughly 90% of the time.
	// Allow some variance
	if countB < 60 || countB > 100 {
		t.Errorf("Expected B to be selected frequently, got %v times out of 100", countB)
	}
}
