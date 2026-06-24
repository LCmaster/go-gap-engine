# Go GAP Engine (Genetic Algorithm & Programming)

A modular, type-safe, and extensible Genetic Algorithm (GA) and Genetic Programming (GP) engine written in Go.

By leveraging Go generics, the Go GAP Engine decouples the core evolutionary loop from specific genome representations and genetic operators. This facilitates reuse across diverse optimization problems like symbolic regression, the Traveling Salesperson Problem (TSP), and more.

## Features

- **Type-Safe Evolution**: Built entirely using Go generics for compile-time safety and flexibility.
- **Genetic Algorithms (GA)**: Supports classic array/slice representations, custom permutations, etc.
- **Genetic Programming (GP)**: Built-in support for Abstract Syntax Trees (ASTs) for solving symbolic regression and other GP problems.
- **Concurrency**: Fitness evaluation is highly parallelized to speed up computations, customizable via the `ConcurrencyLevel` setting.
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
	"fmt"
	"github.com/LCmaster/go-gap-engine/engine"
	// Import your preferred selection and operators
)

func main() {
    // 1. Define your Configuration
	cfg := engine.Config[YourType]{
		PopulationSize:   100,
		Generations:      50,
		MutationRate:     0.1,
		CrossoverRate:    0.8,
		ElitismCount:     2,
		ConcurrencyLevel: 4,
		
		InitFunc:      yourInitFunc,
		FitnessFunc:   yourFitnessFunc,
		SelectionFunc: yourSelectionFunc,
		CrossoverFunc: yourCrossoverFunc,
		MutationFunc:  yourMutationFunc,
		
		OnGeneration: func(gen int, best YourType, bestFit float64, avgFit float64) {
			fmt.Printf("Generation %d: Best Fitness = %f\n", gen, bestFit)
		},
	}

    // 2. Initialize Engine
	eng := engine.New(cfg)
	
	// 3. Evolve
	best, bestFitness := eng.Evolve()
	
	fmt.Printf("Evolution complete! Best fitness: %f\n", bestFitness)
}
```

## Structure

- `engine/`: Core evolutionary loop and concurrency management.
- `ga/`: Genetic Algorithm specific operators (e.g., Order Crossover, Swap Mutation).
- `gp/`: Genetic Programming specific features (e.g., Tree representations, Primitive sets, Subtree Mutation).
- `selection/`: Shared selection algorithms (e.g., Tournament selection).
- `examples/`: Practical examples demonstrating how to use the engine.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
