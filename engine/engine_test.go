package engine_test

import (
	"testing"
	"github.com/LCmaster/go-gap-engine/engine"
)

func TestEngineEvolve(t *testing.T) {
	cfg := engine.Config[int]{
		PopulationSize:   10,
		Generations:      5,
		MutationRate:     0.1,
		CrossoverRate:    0.9,
		ElitismCount:     1,
		ConcurrencyLevel: 2,
		InitFunc: func() int {
			return 1
		},
		FitnessFunc: func(ind int) float64 {
			return float64(ind)
		},
		SelectionFunc: func(pop []int, fits []float64, num int) []int {
			return []int{pop[0], pop[0]}
		},
		CrossoverFunc: func(p1, p2 int) (int, int) {
			return p1, p2
		},
		MutationFunc: func(ind int, rate float64) int {
			return ind + 1
		},
	}

	eng := engine.New(cfg)
	best, bestFit := eng.Evolve()

	if bestFit < 1 {
		t.Errorf("Expected fitness to improve or stay same, got %v", bestFit)
	}
	if best < 1 {
		t.Errorf("Expected best individual to be >= 1, got %v", best)
	}
}
