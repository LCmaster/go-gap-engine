package operators

import (
	"math/rand/v2"

	"github.com/LCmaster/go-gap-engine/engine"
)

// SinglePointCrossover creates a CrossoverFunc for any slice-based genome.
func SinglePointCrossover[S ~[]E, E any]() engine.CrossoverFunc[S] {
	return func(rng *rand.Rand, p1, p2 S) (S, S) {
		if len(p1) == 0 || len(p1) != len(p2) {
			o1, o2 := make(S, len(p1)), make(S, len(p2))
			copy(o1, p1)
			copy(o2, p2)
			return o1, o2
		}

		o1, o2 := make(S, len(p1)), make(S, len(p2))
		point := rng.IntN(len(p1))

		copy(o1[:point], p1[:point])
		copy(o1[point:], p2[point:])

		copy(o2[:point], p2[:point])
		copy(o2[point:], p1[point:])

		return o1, o2
	}
}

// UniformCrossover creates a CrossoverFunc for any slice-based genome.
func UniformCrossover[S ~[]E, E any]() engine.CrossoverFunc[S] {
	return func(rng *rand.Rand, p1, p2 S) (S, S) {
		if len(p1) == 0 || len(p1) != len(p2) {
			o1, o2 := make(S, len(p1)), make(S, len(p2))
			copy(o1, p1)
			copy(o2, p2)
			return o1, o2
		}

		o1, o2 := make(S, len(p1)), make(S, len(p2))
		for i := 0; i < len(p1); i++ {
			if rng.Float64() < 0.5 {
				o1[i], o2[i] = p1[i], p2[i]
			} else {
				o1[i], o2[i] = p2[i], p1[i]
			}
		}

		return o1, o2
	}
}

// OrderCrossover is specifically for permutations (represented as slices of comparable elements).
// It preserves permutations without duplicates.
func OrderCrossover[S ~[]E, E comparable]() engine.CrossoverFunc[S] {
	return func(rng *rand.Rand, p1, p2 S) (S, S) {
		if len(p1) == 0 || len(p1) != len(p2) {
			o1, o2 := make(S, len(p1)), make(S, len(p2))
			copy(o1, p1)
			copy(o2, p2)
			return o1, o2
		}

		length := len(p1)
		o1, o2 := make(S, length), make(S, length)
		
		start := rng.IntN(length)
		end := rng.IntN(length)
		if start > end {
			start, end = end, start
		}

		// Copy substring
		set1 := make(map[E]bool)
		set2 := make(map[E]bool)
		for i := start; i <= end; i++ {
			o1[i] = p1[i]
			o2[i] = p2[i]
			set1[p1[i]] = true
			set2[p2[i]] = true
		}

		// Fill the rest
		fillIdx1 := (end + 1) % length
		fillIdx2 := (end + 1) % length
		
		for i := 0; i < length; i++ {
			idx := (end + 1 + i) % length
			
			if !set1[p2[idx]] {
				o1[fillIdx1] = p2[idx]
				fillIdx1 = (fillIdx1 + 1) % length
			}
			
			if !set2[p1[idx]] {
				o2[fillIdx2] = p1[idx]
				fillIdx2 = (fillIdx2 + 1) % length
			}
		}

		return o1, o2
	}
}
