package operators

import (
	"math/rand/v2"

	"github.com/LCmaster/go-gap-engine/engine"
	"github.com/LCmaster/go-gap-engine/gp/tree"
)

// SubtreeCrossover returns a crossover function for GP trees.
// It selects a random node in both parents and swaps the subtrees rooted at those nodes.
func SubtreeCrossover() engine.CrossoverFunc[tree.Tree] {
	return func(rng *rand.Rand, p1, p2 tree.Tree) (tree.Tree, tree.Tree) {
		o1 := p1.Clone()
		o2 := p2.Clone()

		if o1.Root == nil || o2.Root == nil {
			return o1, o2
		}

		// Collect all nodes in both trees with their parent pointers
		nodes1 := collectNodes(o1.Root)
		nodes2 := collectNodes(o2.Root)

		if len(nodes1) == 0 || len(nodes2) == 0 {
			return o1, o2
		}

		// Select random crossover points
		idx1 := rng.IntN(len(nodes1))
		idx2 := rng.IntN(len(nodes2))

		target1 := nodes1[idx1]
		target2 := nodes2[idx2]

		// Swap subtrees
		// We need to swap the actual node contents or replace them in their parents
		// An easy way in Go is to swap the struct values they point to
		tempType := target1.Type
		tempValue := target1.Value
		tempChildren := target1.Children

		target1.Type = target2.Type
		target1.Value = target2.Value
		target1.Children = target2.Children

		target2.Type = tempType
		target2.Value = tempValue
		target2.Children = tempChildren

		return o1, o2
	}
}

// collectNodes returns a flattened slice of pointers to all nodes in the tree
func collectNodes(root *tree.Node) []*tree.Node {
	if root == nil {
		return nil
	}
	nodes := []*tree.Node{root}
	for _, child := range root.Children {
		nodes = append(nodes, collectNodes(child)...)
	}
	return nodes
}
