package main

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/LCmaster/go-gap-engine/engine"
	"github.com/LCmaster/go-gap-engine/ga/operators"
	"github.com/LCmaster/go-gap-engine/ga/types"
	"github.com/LCmaster/go-gap-engine/selection"
)

type Point struct {
	x, y float64
}

func dist(p1, p2 Point) float64 {
	dx := p1.x - p2.x
	dy := p1.y - p2.y
	return math.Sqrt(dx*dx + dy*dy)
}

func main() {
	numCities := 20
	cities := make([]Point, numCities)
	
	// Use global random for initialization of the problem
	for i := range cities {
		cities[i] = Point{rand.Float64() * 100, rand.Float64() * 100}
	}

	seed := [32]byte{1} // fixed seed for reproducible evolution

	cfg := engine.Config[types.Permutation]{
		PopulationSize:   100,
		Generations:      100,
		MutationRate:     0.1,
		CrossoverRate:    0.8,
		ElitismCount:     2,
		ConcurrencyLevel: 4,
		Seed:             &seed,
		InitFunc: func(rng *rand.Rand) types.Permutation {
			p := make(types.Permutation, numCities)
			for i := range p {
				p[i] = i
			}
			rng.Shuffle(numCities, func(i, j int) { p[i], p[j] = p[j], p[i] })
			return p
		},
		FitnessFunc: func(p types.Permutation) float64 {
			totalDist := 0.0
			for i := 0; i < len(p)-1; i++ {
				totalDist += dist(cities[p[i]], cities[p[i+1]])
			}
			totalDist += dist(cities[p[len(p)-1]], cities[p[0]])
			// Maximize fitness -> use negative distance
			return -totalDist
		},
		SelectionFunc: selection.Tournament[types.Permutation](5),
		CrossoverFunc: operators.OrderCrossover[types.Permutation, int](),
		MutationFunc:  operators.Swap[types.Permutation, int](),
		OnGeneration: func(gen int, best types.Permutation, bestFit float64, avgFit float64) {
			if gen%10 == 0 || gen == 99 {
				fmt.Printf("Gen %3d: Best Fitness (Neg Dist) = %8.2f, Avg = %8.2f\n", gen, bestFit, avgFit)
			}
		},
	}

	eng, err := engine.New(cfg)
	if err != nil {
		panic(err)
	}
	best, fit := eng.Evolve()
	fmt.Printf("Final Best Distance: %.2f\nPath: %v\n", -fit, best)
}
