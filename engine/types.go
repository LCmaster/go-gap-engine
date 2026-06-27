package engine

import (
	"errors"
	"math/rand/v2"
)

var (
	ErrInvalidPopulationSize = errors.New("population size must be > 0")
	ErrInvalidGenerations    = errors.New("generations must be > 0")
	ErrMissingInitFunc       = errors.New("initFunc is required")
	ErrMissingFitnessFunc    = errors.New("fitnessFunc is required")
	ErrMissingSelectionFunc  = errors.New("selectionFunc is required")
	ErrMissingCrossoverFunc  = errors.New("crossoverFunc is required")
	ErrMissingMutationFunc   = errors.New("mutationFunc is required")
	ErrInvalidMutationRate   = errors.New("mutation rate must be between 0 and 1")
	ErrInvalidCrossoverRate  = errors.New("crossover rate must be between 0 and 1")
	ErrInvalidElitismCount   = errors.New("elitism count must be >= 0 and <= population size")
)

// InitFunc generates a single random individual for the initial population.
type InitFunc[T any] func(rng *rand.Rand) T

// FitnessFunc evaluates the fitness of a single individual.
// Higher values should represent better fitness (maximization).
type FitnessFunc[T any] func(individual T) float64

// SelectionFunc selects 'num' individuals from the population based on their fitnesses.
type SelectionFunc[T any] func(rng *rand.Rand, population []T, fitnesses []float64, num int) []T

// CrossoverFunc performs crossover between two parents to produce two offspring.
type CrossoverFunc[T any] func(rng *rand.Rand, p1, p2 T) (T, T)

// MutationFunc mutates an individual based on the given mutation rate.
type MutationFunc[T any] func(rng *rand.Rand, individual T, rate float64) T

// Option configures the Engine via the functional options pattern.
type Option[T any] func(*config[T])

// config holds the internal configuration for the Engine.
type config[T any] struct {
	populationSize   int
	generations      int
	mutationRate     float64
	crossoverRate    float64
	elitismCount     int
	concurrencyLevel int
	seed             *[32]byte

	initFunc      InitFunc[T]
	fitnessFunc   FitnessFunc[T]
	selectionFunc SelectionFunc[T]
	crossoverFunc CrossoverFunc[T]
	mutationFunc  MutationFunc[T]

	onGeneration func(generation int, best T, bestFitness float64, avgFitness float64)
}

// WithPopulationSize sets the population size.
func WithPopulationSize[T any](size int) Option[T] {
	return func(c *config[T]) {
		c.populationSize = size
	}
}

// WithGenerations sets the number of generations to run.
func WithGenerations[T any](generations int) Option[T] {
	return func(c *config[T]) {
		c.generations = generations
	}
}

// WithMutationRate sets the mutation rate (0.0 to 1.0).
func WithMutationRate[T any](rate float64) Option[T] {
	return func(c *config[T]) {
		c.mutationRate = rate
	}
}

// WithCrossoverRate sets the crossover rate (0.0 to 1.0).
func WithCrossoverRate[T any](rate float64) Option[T] {
	return func(c *config[T]) {
		c.crossoverRate = rate
	}
}

// WithElitismCount sets how many of the best individuals carry over to the next generation.
func WithElitismCount[T any](count int) Option[T] {
	return func(c *config[T]) {
		c.elitismCount = count
	}
}

// WithConcurrencyLevel sets the number of goroutines for evaluating fitness.
func WithConcurrencyLevel[T any](level int) Option[T] {
	return func(c *config[T]) {
		c.concurrencyLevel = level
	}
}

// WithSeed sets a deterministic seed for the RNG.
func WithSeed[T any](seed [32]byte) Option[T] {
	return func(c *config[T]) {
		c.seed = &seed
	}
}

// WithInitFunc sets the initialization function.
func WithInitFunc[T any](f InitFunc[T]) Option[T] {
	return func(c *config[T]) {
		c.initFunc = f
	}
}

// WithFitnessFunc sets the fitness evaluation function.
func WithFitnessFunc[T any](f FitnessFunc[T]) Option[T] {
	return func(c *config[T]) {
		c.fitnessFunc = f
	}
}

// WithSelectionFunc sets the selection function.
func WithSelectionFunc[T any](f SelectionFunc[T]) Option[T] {
	return func(c *config[T]) {
		c.selectionFunc = f
	}
}

// WithCrossoverFunc sets the crossover function.
func WithCrossoverFunc[T any](f CrossoverFunc[T]) Option[T] {
	return func(c *config[T]) {
		c.crossoverFunc = f
	}
}

// WithMutationFunc sets the mutation function.
func WithMutationFunc[T any](f MutationFunc[T]) Option[T] {
	return func(c *config[T]) {
		c.mutationFunc = f
	}
}

// WithOnGeneration sets the callback for each generation's completion.
func WithOnGeneration[T any](f func(generation int, best T, bestFitness float64, avgFitness float64)) Option[T] {
	return func(c *config[T]) {
		c.onGeneration = f
	}
}
