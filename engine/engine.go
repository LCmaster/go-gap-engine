package engine

import (
	"context"
	"math"
	"math/rand/v2"
	"runtime"
	"sort"
	"sync"
	"sync/atomic"
)

// Engine is the core struct that drives the evolutionary process.
type Engine[T any] struct {
	cfg config[T]
	rng *rand.Rand
}

// New returns a new Engine configured with the given functional options.
func New[T any](opts ...Option[T]) (*Engine[T], error) {
	// Initialize with defaults
	cfg := config[T]{
		concurrencyLevel: runtime.NumCPU(),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.populationSize <= 0 {
		return nil, ErrInvalidPopulationSize
	}
	if cfg.generations <= 0 {
		return nil, ErrInvalidGenerations
	}
	if cfg.initFunc == nil {
		return nil, ErrMissingInitFunc
	}
	if cfg.fitnessFunc == nil {
		return nil, ErrMissingFitnessFunc
	}
	if cfg.selectionFunc == nil {
		return nil, ErrMissingSelectionFunc
	}
	if cfg.crossoverFunc == nil {
		return nil, ErrMissingCrossoverFunc
	}
	if cfg.mutationFunc == nil {
		return nil, ErrMissingMutationFunc
	}
	if cfg.mutationRate < 0 || cfg.mutationRate > 1 {
		return nil, ErrInvalidMutationRate
	}
	if cfg.crossoverRate < 0 || cfg.crossoverRate > 1 {
		return nil, ErrInvalidCrossoverRate
	}
	if cfg.elitismCount < 0 || cfg.elitismCount > cfg.populationSize {
		return nil, ErrInvalidElitismCount
	}

	var rng *rand.Rand
	if cfg.seed != nil {
		rng = rand.New(rand.NewChaCha8(*cfg.seed))
	} else {
		rng = rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))
	}

	return &Engine[T]{cfg: cfg, rng: rng}, nil
}

// Evolve runs the genetic algorithm for the specified number of generations.
func (e *Engine[T]) Evolve(ctx context.Context) (best T, bestFitness float64, err error) {
	// 1. Initialize population
	pop := make([]T, e.cfg.populationSize)
	for i := range pop {
		pop[i] = e.cfg.initFunc(e.rng)
	}

	bestFitness = -math.MaxFloat64

	for gen := 0; gen < e.cfg.generations; gen++ {
		select {
		case <-ctx.Done():
			return best, bestFitness, ctx.Err()
		default:
		}

		// 2. Evaluate fitness concurrently
		fitnesses, errEval := e.EvaluatePopulation(ctx, pop)
		if errEval != nil {
			return best, bestFitness, errEval
		}

		// 3. Track best and average
		var currentBest T
		currentBestFit := -math.MaxFloat64
		var sumFit float64

		type indFit struct {
			ind T
			fit float64
		}
		scoredPop := make([]indFit, e.cfg.populationSize)

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

		avgFit := sumFit / float64(e.cfg.populationSize)

		if e.cfg.onGeneration != nil {
			e.cfg.onGeneration(gen, currentBest, currentBestFit, avgFit)
		}

		// If it's the last generation, we don't need to create a new population
		if gen == e.cfg.generations-1 {
			break
		}

		// Sort population by fitness descending (for elitism)
		sort.SliceStable(scoredPop, func(i, j int) bool {
			return scoredPop[i].fit > scoredPop[j].fit
		})

		newPop := make([]T, 0, e.cfg.populationSize)

		// 4. Elitism
		for i := 0; i < e.cfg.elitismCount && i < len(scoredPop); i++ {
			newPop = append(newPop, scoredPop[i].ind)
		}

		// 5. Selection, Crossover, Mutation
		for len(newPop) < e.cfg.populationSize {
			// Select parents
			parents := e.cfg.selectionFunc(e.rng, pop, fitnesses, 2)
			p1, p2 := parents[0], parents[1]

			var o1, o2 T

			// Crossover
			if e.rng.Float64() < e.cfg.crossoverRate {
				o1, o2 = e.cfg.crossoverFunc(e.rng, p1, p2)
			} else {
				o1, o2 = p1, p2
			}

			// Mutation
			o1 = e.cfg.mutationFunc(e.rng, o1, e.cfg.mutationRate)
			o2 = e.cfg.mutationFunc(e.rng, o2, e.cfg.mutationRate)

			newPop = append(newPop, o1)
			if len(newPop) < e.cfg.populationSize {
				newPop = append(newPop, o2)
			}
		}

		pop = newPop
	}

	return best, bestFitness, nil
}

// EvaluatePopulation evaluates all individuals in the population concurrently.
func (e *Engine[T]) EvaluatePopulation(ctx context.Context, pop []T) ([]float64, error) {
	fitnesses := make([]float64, len(pop))
	var wg sync.WaitGroup

	numWorkers := e.cfg.concurrencyLevel
	if numWorkers > len(pop) {
		numWorkers = len(pop)
	}

	var taskIdx atomic.Int64
	taskIdx.Store(-1)

	var evalErr error
	var errOnce sync.Once

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					errOnce.Do(func() { evalErr = ctx.Err() })
					return
				default:
				}

				idx := int(taskIdx.Add(1))
				if idx >= len(pop) {
					break
				}
				fitnesses[idx] = e.cfg.fitnessFunc(pop[idx])
			}
		}()
	}

	wg.Wait()
	return fitnesses, evalErr
}
