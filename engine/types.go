package engine

// InitFunc generates a single random individual for the initial population.
type InitFunc[T any] func() T

// FitnessFunc evaluates the fitness of a single individual.
// Higher values should represent better fitness (maximization).
type FitnessFunc[T any] func(individual T) float64

// SelectionFunc selects 'num' individuals from the population based on their fitnesses.
type SelectionFunc[T any] func(population []T, fitnesses []float64, num int) []T

// CrossoverFunc performs crossover between two parents to produce two offspring.
type CrossoverFunc[T any] func(p1, p2 T) (T, T)

// MutationFunc mutates an individual based on the given mutation rate.
type MutationFunc[T any] func(individual T, rate float64) T

// Config holds the configuration for the Engine.
type Config[T any] struct {
	PopulationSize   int
	Generations      int
	MutationRate     float64
	CrossoverRate    float64
	ElitismCount     int
	ConcurrencyLevel int // Number of goroutines for fitness evaluation (default: runtime.NumCPU())

	InitFunc      InitFunc[T]
	FitnessFunc   FitnessFunc[T]
	SelectionFunc SelectionFunc[T]
	CrossoverFunc CrossoverFunc[T]
	MutationFunc  MutationFunc[T]

	// Optional callback executed at the end of each generation
	OnGeneration func(generation int, best T, bestFitness float64, avgFitness float64)
}
