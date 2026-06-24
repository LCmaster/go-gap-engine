# Go GAP Engine Examples

This directory contains complete, runnable examples demonstrating how to solve various optimization problems using the Go GAP Engine.

## Running the Examples

You can run any example directly using the `go run` command from the root of the project or inside the specific example's directory.

### 1. Traveling Salesperson Problem (GA)

Demonstrates solving the Traveling Salesperson Problem (TSP) using a classic Genetic Algorithm with permutation genomes, Order Crossover, and Swap Mutation.

```bash
cd ga_tsp
go run main.go
```

### 2. Symbolic Regression (GP)

Demonstrates solving symbolic regression using Genetic Programming. The algorithm evolves an Abstract Syntax Tree (AST) to discover the underlying mathematical function `f(x) = x^2 + 2x + 1` given a set of data points.

```bash
cd gp_symbolic_regression
go run main.go
```
