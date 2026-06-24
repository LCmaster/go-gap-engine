package operators

import (
	"math/rand/v2"

	"github.com/LCmaster/go-gap-engine/engine"
)

// BitFlip creates a MutationFunc for boolean slices (e.g. BitString).
// The rate is the probability of flipping each individual bit.
func BitFlip[S ~[]bool]() engine.MutationFunc[S] {
	return func(rng *rand.Rand, ind S, rate float64) S {
		o := make(S, len(ind))
		copy(o, ind)
		for i := 0; i < len(o); i++ {
			if rng.Float64() < rate {
				o[i] = !o[i]
			}
		}
		return o
	}
}

// Gaussian creates a MutationFunc for float64 slices (e.g. FloatVector).
// The rate is the probability of mutating each individual gene.
// The standardDeviation controls the magnitude of the Gaussian noise added.
func Gaussian[S ~[]float64](standardDeviation float64) engine.MutationFunc[S] {
	return func(rng *rand.Rand, ind S, rate float64) S {
		o := make(S, len(ind))
		copy(o, ind)
		for i := 0; i < len(o); i++ {
			if rng.Float64() < rate {
				o[i] += rng.NormFloat64() * standardDeviation
			}
		}
		return o
	}
}

// Swap creates a MutationFunc for any slice-based genome.
// The rate is the probability that a swap occurs. If it occurs, two random positions are swapped.
func Swap[S ~[]E, E any]() engine.MutationFunc[S] {
	return func(rng *rand.Rand, ind S, rate float64) S {
		o := make(S, len(ind))
		copy(o, ind)
		if rng.Float64() < rate && len(o) > 1 {
			i := rng.IntN(len(o))
			j := rng.IntN(len(o))
			o[i], o[j] = o[j], o[i]
		}
		return o
	}
}
