package engine

import (
	"math/rand"
	"runtime"
	"sort"
	"sync"
)

// Engine is the core struct that drives the evolutionary process.
type Engine[T any] struct {
	cfg Config[T]
}

// New returns a new Engine configured with the given Config.
func New[T any](cfg Config[T]) *Engine[T] {
	if cfg.ConcurrencyLevel <= 0 {
		cfg.ConcurrencyLevel = runtime.NumCPU()
	}
	return &Engine[T]{cfg: cfg}
}

// Evolve runs the genetic algorithm for the specified number of generations.
func (e *Engine[T]) Evolve() (best T, bestFitness float64) {
	// 1. Initialize population
	pop := make([]T, e.cfg.PopulationSize)
	for i := range pop {
		pop[i] = e.cfg.InitFunc()
	}

	bestFitness = -1.7976931348623157e+308 // math.MaxFloat64 * -1 (lowest possible)

	for gen := 0; gen < e.cfg.Generations; gen++ {
		// 2. Evaluate fitness concurrently
		fitnesses := e.EvaluatePopulation(pop)

		// 3. Track best and average
		var currentBest T
		currentBestFit := -1.7976931348623157e+308
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
			parents := e.cfg.SelectionFunc(pop, fitnesses, 2)
			p1, p2 := parents[0], parents[1]

			var o1, o2 T

			// Crossover
			if rand.Float64() < e.cfg.CrossoverRate {
				o1, o2 = e.cfg.CrossoverFunc(p1, p2)
			} else {
				o1, o2 = p1, p2
			}

			// Mutation
			o1 = e.cfg.MutationFunc(o1, e.cfg.MutationRate)
			o2 = e.cfg.MutationFunc(o2, e.cfg.MutationRate)

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
	
	// Create a channel for task indices
	tasks := make(chan int, len(pop))
	for i := range pop {
		tasks <- i
	}
	close(tasks)

	numWorkers := e.cfg.ConcurrencyLevel
	if numWorkers > len(pop) {
		numWorkers = len(pop)
	}

	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range tasks {
				fitnesses[i] = e.cfg.FitnessFunc(pop[i])
			}
		}()
	}

	wg.Wait()
	return fitnesses
}
