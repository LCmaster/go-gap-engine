package main

import (
	"fmt"
	"math"
	"strconv"

	"github.com/LCmaster/go-gap-engine/engine"
	"github.com/LCmaster/go-gap-engine/gp/operators"
	"github.com/LCmaster/go-gap-engine/gp/tree"
	"github.com/LCmaster/go-gap-engine/selection"
)

// evaluate evaluates the tree given a value for 'X'.
func evaluate(node *tree.Node, x float64) float64 {
	if node.Type == tree.TerminalNode {
		if node.Value == "X" {
			return x
		}
		val, _ := strconv.ParseFloat(node.Value, 64)
		return val
	}

	left := evaluate(node.Children[0], x)
	right := evaluate(node.Children[1], x)

	switch node.Value {
	case "ADD":
		return left + right
	case "SUB":
		return left - right
	case "MUL":
		return left * right
	case "DIV":
		if right == 0 {
			return 1 // Protected division
		}
		return left / right
	}
	return 0
}

func printTree(node *tree.Node) string {
	if node.Type == tree.TerminalNode {
		return node.Value
	}
	return fmt.Sprintf("(%s %s %s)", printTree(node.Children[0]), node.Value, printTree(node.Children[1]))
}

func main() {
	// Target function: f(x) = x^2 + 2x + 1
	// Data points
	type point struct{ x, y float64 }
	var data []point
	for i := -10; i <= 10; i++ {
		x := float64(i)
		data = append(data, point{x, x*x + 2*x + 1})
	}

	pset := tree.PrimitiveSet{
		Functions: []string{"ADD", "SUB", "MUL", "DIV"},
		Terminals: []string{"X", "1", "2"},
		Arity: map[string]int{
			"ADD": 2, "SUB": 2, "MUL": 2, "DIV": 2,
		},
	}

	cfg := engine.Config[tree.Tree]{
		PopulationSize:   200,
		Generations:      50,
		MutationRate:     0.2,
		CrossoverRate:    0.7,
		ElitismCount:     1,
		ConcurrencyLevel: 4,
		InitFunc: func() tree.Tree {
			return tree.Tree{Root: tree.GenerateGrow(4, pset)}
		},
		FitnessFunc: func(t tree.Tree) float64 {
			var errSum float64
			for _, p := range data {
				pred := evaluate(t.Root, p.x)
				diff := pred - p.y
				errSum += diff * diff
			}
			mse := errSum / float64(len(data))
			if math.IsNaN(mse) || math.IsInf(mse, 0) {
				return -1e9
			}
			// Maximize fitness -> use negative MSE
			return -mse
		},
		SelectionFunc: selection.Tournament[tree.Tree](3),
		CrossoverFunc: operators.SubtreeCrossover(),
		MutationFunc:  operators.SubtreeMutation(4, pset),
		OnGeneration: func(gen int, best tree.Tree, bestFit float64, avgFit float64) {
			if gen%5 == 0 || gen == 49 {
				fmt.Printf("Gen %3d: Best Fitness (Neg MSE) = %8.4f, Expr: %s\n", gen, bestFit, printTree(best.Root))
			}
		},
	}

	eng := engine.New(cfg)
	best, fit := eng.Evolve()
	fmt.Printf("\nFinal Best Fitness: %.4f\nExpression: %s\n", fit, printTree(best.Root))
}
