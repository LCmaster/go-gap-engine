package engine_test

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/LCmaster/go-gap-engine/engine"
)

func mockInitFunc(rng *rand.Rand) int { return 1 }
func mockFitnessFunc(ind int) float64 { return float64(ind) }
func mockSelectionFunc(rng *rand.Rand, pop []int, fits []float64, num int) []int { return []int{pop[0], pop[0]} }
func mockCrossoverFunc(rng *rand.Rand, p1, p2 int) (int, int) { return p1, p2 }
func mockMutationFunc(rng *rand.Rand, ind int, rate float64) int { return ind + 1 }

func TestEngineEvolve(t *testing.T) {
	tests := []struct {
		name        string
		options     []engine.Option[int]
		wantErr     bool
		checkResult func(t *testing.T, best int, bestFit float64, err error)
	}{
		{
			name: "Successful Evolution",
			options: []engine.Option[int]{
				engine.WithPopulationSize[int](10),
				engine.WithGenerations[int](5),
				engine.WithMutationRate[int](0.1),
				engine.WithCrossoverRate[int](0.9),
				engine.WithElitismCount[int](1),
				engine.WithConcurrencyLevel[int](2),
				engine.WithInitFunc(mockInitFunc),
				engine.WithFitnessFunc(mockFitnessFunc),
				engine.WithSelectionFunc(mockSelectionFunc),
				engine.WithCrossoverFunc(mockCrossoverFunc),
				engine.WithMutationFunc(mockMutationFunc),
			},
			wantErr: false,
			checkResult: func(t *testing.T, best int, bestFit float64, err error) {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if bestFit < 1 {
					t.Errorf("Expected fitness >= 1, got %v", bestFit)
				}
				if best < 1 {
					t.Errorf("Expected best individual >= 1, got %v", best)
				}
			},
		},
		{
			name: "Missing InitFunc",
			options: []engine.Option[int]{
				engine.WithPopulationSize[int](10),
				engine.WithGenerations[int](5),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eng, err := engine.New(tt.options...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}
			
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			
			best, bestFit, err := eng.Evolve(ctx)
			if tt.checkResult != nil {
				tt.checkResult(t, best, bestFit, err)
			}
		})
	}
}

func BenchmarkEngineEvolve(b *testing.B) {
	eng, _ := engine.New(
		engine.WithPopulationSize[int](100),
		engine.WithGenerations[int](50),
		engine.WithMutationRate[int](0.1),
		engine.WithCrossoverRate[int](0.9),
		engine.WithElitismCount[int](5),
		engine.WithInitFunc(mockInitFunc),
		engine.WithFitnessFunc(mockFitnessFunc),
		engine.WithSelectionFunc(mockSelectionFunc),
		engine.WithCrossoverFunc(mockCrossoverFunc),
		engine.WithMutationFunc(mockMutationFunc),
	)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		eng.Evolve(context.Background())
	}
}
