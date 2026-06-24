package engine

import (
	"errors"
	"math"
	"math/rand/v2"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
)

// Engine is the core struct that drives the evolutionary process.
type Engine[T any] struct {
	cfg Config[T]
	rng *rand.Rand
}

// New returns a new Engine configured with the given Config.
func New[T any](cfg Config[T]) (*Engine[T], error) {
	if cfg.PopulationSize <= 0 {
		return nil, errors.New("PopulationSize must be > 0")
	}
	if cfg.Generations <= 0 {
		return nil, errors.New("Generations must be > 0")
	}
	if cfg.InitFunc == nil {
		return nil, errors.New("InitFunc is required")
	}
	if cfg.FitnessFunc == nil {
		return nil, errors.New("FitnessFunc is required")
	}
	if cfg.SelectionFunc == nil {
		return nil, errors.New("SelectionFunc is required")
	}
	if cfg.CrossoverFunc == nil {
		return nil, errors.New("CrossoverFunc is required")
	}
	if cfg.MutationFunc == nil {
		return nil, errors.New("MutationFunc is required")
	}
	if cfg.MutationRate < 0 || cfg.MutationRate > 1 {
		return nil, errors.New("MutationRate must be between 0 and 1")
	}
	if cfg.CrossoverRate < 0 || cfg.CrossoverRate > 1 {
		return nil, errors.New("CrossoverRate must be between 0 and 1")
	}
	if cfg.ElitismCount < 0 || cfg.ElitismCount > cfg.PopulationSize {
		return nil, errors.New("ElitismCount must be between 0 and PopulationSize")
	}

	if cfg.ConcurrencyLevel <= 0 {
		cfg.ConcurrencyLevel = runtime.NumCPU()
	}

	var rng *rand.Rand
	if cfg.Seed != nil {
		rng = rand.New(rand.NewChaCha8(*cfg.Seed))
	} else {
		rng = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}

	return &Engine[T]{cfg: cfg, rng: rng}, nil
}

// Evolve runs the genetic algorithm for the specified number of generations.
func (e *Engine[T]) Evolve() (best T, bestFitness float64) {
	// 1. Initialize population
	pop := make([]T, e.cfg.PopulationSize)
	for i := range pop {
		pop[i] = e.cfg.InitFunc(e.rng)
	}

	bestFitness = -math.MaxFloat64

	for gen := 0; gen < e.cfg.Generations; gen++ {
		// 2. Evaluate fitness concurrently
		fitnesses := e.EvaluatePopulation(pop)

		// 3. Track best and average
		var currentBest T
		currentBestFit := -math.MaxFloat64
		var sumFit float64

		type indFit struct {
			ind T
			fit float64
		}
		scoredPop := make([]indFit, e.cfg.PopulationSize)

		for i, f := range fitnesses {
			sumFit += f
			scoredPop[i] = indFit{ind: pop[i], fit: f}
			if f > currentBestFit {
				currentBestFit = f
				currentBest = pop[i]
			}
		}

		if currentBestFit > bestFitness {
			bestFitness = currentBestFit
			best = currentBest
		}

		avgFit := sumFit / float64(e.cfg.PopulationSize)

		if e.cfg.OnGeneration != nil {
			e.cfg.OnGeneration(gen, currentBest, currentBestFit, avgFit)
		}

		// If it's the last generation, we don't need to create a new population
		if gen == e.cfg.Generations-1 {
			break
		}

		// Sort population by fitness descending (for elitism)
		sort.SliceStable(scoredPop, func(i, j int) bool {
			return scoredPop[i].fit > scoredPop[j].fit
		})

		newPop := make([]T, 0, e.cfg.PopulationSize)

		// 4. Elitism
		for i := 0; i < e.cfg.ElitismCount && i < len(scoredPop); i++ {
			newPop = append(newPop, scoredPop[i].ind)
		}

		// 5. Selection, Crossover, Mutation
		for len(newPop) < e.cfg.PopulationSize {
			// Select parents
			parents := e.cfg.SelectionFunc(e.rng, pop, fitnesses, 2)
			p1, p2 := parents[0], parents[1]

			var o1, o2 T

			// Crossover
			if e.rng.Float64() < e.cfg.CrossoverRate {
				o1, o2 = e.cfg.CrossoverFunc(e.rng, p1, p2)
			} else {
				o1, o2 = p1, p2
			}

			// Mutation
			o1 = e.cfg.MutationFunc(e.rng, o1, e.cfg.MutationRate)
			o2 = e.cfg.MutationFunc(e.rng, o2, e.cfg.MutationRate)

			newPop = append(newPop, o1)
			if len(newPop) < e.cfg.PopulationSize {
				newPop = append(newPop, o2)
			}
		}

		pop = newPop
	}

	return best, bestFitness
}

// EvaluatePopulation evaluates all individuals in the population concurrently.
func (e *Engine[T]) EvaluatePopulation(pop []T) []float64 {
	fitnesses := make([]float64, len(pop))
	var wg sync.WaitGroup

	numWorkers := e.cfg.ConcurrencyLevel
	if numWorkers > len(pop) {
		numWorkers = len(pop)
	}

	var taskIdx atomic.Int64
	taskIdx.Store(-1)

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				idx := int(taskIdx.Add(1))
				if idx >= len(pop) {
					break
				}
				fitnesses[idx] = e.cfg.FitnessFunc(pop[idx])
			}
		}()
	}

	wg.Wait()
	return fitnesses
}
