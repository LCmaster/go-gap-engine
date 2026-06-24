package operators

import (
	"math/rand"

	"go-gap-engine/engine"
	"go-gap-engine/gp/tree"
)

// SubtreeMutation returns a mutation function that replaces a random node with a newly generated subtree.
func SubtreeMutation(maxDepth int, pset tree.PrimitiveSet) engine.MutationFunc[tree.Tree] {
	return func(ind tree.Tree, rate float64) tree.Tree {
		if rand.Float64() >= rate {
			return ind
		}

		o := ind.Clone()
		if o.Root == nil {
			o.Root = tree.GenerateGrow(maxDepth, pset)
			return o
		}

		nodes := collectNodes(o.Root)
		target := nodes[rand.Intn(len(nodes))]

		// Generate a new subtree
		newSubtree := tree.GenerateGrow(maxDepth, pset)

		// Replace the target node with the new subtree
		target.Type = newSubtree.Type
		target.Value = newSubtree.Value
		target.Children = newSubtree.Children

		return o
	}
}

// PointMutation returns a mutation function that randomly changes the function or terminal value of a node.
func PointMutation(pset tree.PrimitiveSet) engine.MutationFunc[tree.Tree] {
	return func(ind tree.Tree, rate float64) tree.Tree {
		if rand.Float64() >= rate {
			return ind
		}

		o := ind.Clone()
		if o.Root == nil {
			return o
		}

		nodes := collectNodes(o.Root)
		target := nodes[rand.Intn(len(nodes))]

		if target.Type == tree.FunctionNode {
			// Find another function with the same arity
			arity := len(target.Children)
			var validFuncs []string
			for _, f := range pset.Functions {
				if pset.Arity[f] == arity {
					validFuncs = append(validFuncs, f)
				}
			}
			if len(validFuncs) > 0 {
				target.Value = validFuncs[rand.Intn(len(validFuncs))]
			}
		} else {
			// Mutate terminal
			target.Value = pset.Terminals[rand.Intn(len(pset.Terminals))]
		}

		return o
	}
}
