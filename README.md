# Go GAP Engine (Genetic Algorithm & Programming)

A modular, type-safe, and extensible Genetic Algorithm (GA) and Genetic Programming (GP) engine written in Go.

By leveraging Go generics, the Go GAP Engine decouples the core evolutionary loop from specific genome representations and genetic operators. This facilitates reuse across diverse optimization problems like symbolic regression, the Traveling Salesperson Problem (TSP), and more.

## Features

- **Type-Safe Evolution**: Built entirely using Go generics for compile-time safety and flexibility.
- **Genetic Algorithms (GA)**: Supports classic array/slice representations, custom permutations, etc.
- **Genetic Programming (GP)**: Built-in support for Abstract Syntax Trees (ASTs) for solving symbolic regression and other GP problems.
- **Concurrency & Context-Awareness**: Fitness evaluation is highly parallelized and gracefully handles `context.Context` cancellation to prevent goroutine leaks.
- **Functional Options**: Clean and extensible initialization via the `engine.Option` pattern.
- **Reproducible Runs**: Pass `engine.WithSeed(seed)` to get fully deterministic evolution across runs.
- **Extensible**: Easily provide your own initialization, selection, crossover, mutation, and fitness functions.

## Installation

```bash
go get github.com/LCmaster/go-gap-engine
```


## Quick Start

See the `examples/` directory for full implementations:

- **[Traveling Salesperson Problem (TSP)](examples/ga_tsp)**: Uses Genetic Algorithms with permutation arrays.
- **[Symbolic Regression](examples/gp_symbolic_regression)**: Uses Genetic Programming to evolve an AST that fits the function `f(x) = x^2 + 2x + 1`.

### Basic Usage Pattern

```go
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"

	"github.com/LCmaster/go-gap-engine/engine"
	// Import your preferred selection and operators
)

func main() {
    // 1. Define your Options
	opts := []engine.Option[YourType]{
		engine.WithPopulationSize[YourType](100),
		engine.WithGenerations[YourType](50),
		engine.WithMutationRate[YourType](0.1),
		engine.WithCrossoverRate[YourType](0.8),
		engine.WithElitismCount[YourType](2),
		engine.WithConcurrencyLevel[YourType](4),

		// Optional: provide a fixed seed for reproducible results
		// engine.WithSeed[YourType]([32]byte{1, 2, 3}),

		// All operator functions now receive a *rand.Rand for safe, reproducible randomness
		engine.WithInitFunc(func(rng *rand.Rand) YourType {
			return yourInitLogic(rng)
		}),
		engine.WithFitnessFunc(func(ind YourType) float64 {
			return yourFitnessLogic(ind)
		}),
		engine.WithSelectionFunc(func(rng *rand.Rand, pop []YourType, fits []float64, num int) []YourType {
			return yourSelectionLogic(rng, pop, fits, num)
		}),
		engine.WithCrossoverFunc(func(rng *rand.Rand, p1, p2 YourType) (YourType, YourType) {
			return yourCrossoverLogic(rng, p1, p2)
		}),
		engine.WithMutationFunc(func(rng *rand.Rand, ind YourType, rate float64) YourType {
			return yourMutationLogic(rng, ind, rate)
		}),

		engine.WithOnGeneration(func(gen int, best YourType, bestFit float64, avgFit float64) {
			fmt.Printf("Generation %d: Best Fitness = %f\n", gen, bestFit)
		}),
	}

    // 2. Initialize Engine (returns an error if options are invalid)
	eng, err := engine.New(opts...)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Evolve with Context
	best, bestFitness, err := eng.Evolve(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Evolution complete! Best fitness: %f\n", bestFitness)
}
```

## Structure

- `engine/`: Core evolutionary loop and concurrency management.
- `ga/`: Genetic Algorithm specific operators (e.g., Order Crossover, Swap Mutation).
- `gp/`: Genetic Programming specific features (e.g., Tree representations, Primitive sets, Subtree Mutation).
- `selection/`: Shared selection algorithms (e.g., Tournament selection, Roulette Wheel).
- `examples/`: Practical examples demonstrating how to use the engine.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
