package selection

import (
	"math/rand/v2"

	"github.com/LCmaster/go-gap-engine/engine"
)

// Tournament returns a SelectionFunc that performs tournament selection.
func Tournament[T any](tournamentSize int) engine.SelectionFunc[T] {
	return func(rng *rand.Rand, pop []T, fitnesses []float64, num int) []T {
		selected := make([]T, num)
		for i := 0; i < num; i++ {
			bestIdx := -1
			bestFit := -1.7976931348623157e+308 // math.MaxFloat64 * -1

			for j := 0; j < tournamentSize; j++ {
				idx := rng.IntN(len(pop))
				if bestIdx == -1 || fitnesses[idx] > bestFit {
					bestIdx = idx
					bestFit = fitnesses[idx]
				}
			}
			selected[i] = pop[bestIdx]
		}
		return selected
	}
}

// RouletteWheel returns a SelectionFunc that performs roulette wheel (fitness proportionate) selection.
// Note: This implementation requires all fitness values to be positive.
func RouletteWheel[T any]() engine.SelectionFunc[T] {
	return func(rng *rand.Rand, pop []T, fitnesses []float64, num int) []T {
		selected := make([]T, num)
		
		var sumFit float64
		for _, f := range fitnesses {
			if f > 0 {
				sumFit += f
			}
		}

		for i := 0; i < num; i++ {
			r := rng.Float64() * sumFit
			var currentSum float64
			selectedIndex := len(pop) - 1 // fallback

			for j, f := range fitnesses {
				if f > 0 {
					currentSum += f
					if currentSum >= r {
						selectedIndex = j
						break
					}
				}
			}
			selected[i] = pop[selectedIndex]
		}
		
		return selected
	}
}
