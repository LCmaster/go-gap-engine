package tree

import (
	"math/rand/v2"
)

type NodeType int

const (
	FunctionNode NodeType = iota
	TerminalNode
)

// Node represents a single node in the abstract syntax tree.
type Node struct {
	Type     NodeType
	Value    string // The name of the function or the terminal value (e.g. "ADD", "X", "5")
	Children []*Node
}

// Clone creates a deep copy of the node and its children.
func (n *Node) Clone() *Node {
	if n == nil {
		return nil
	}
	clone := &Node{
		Type:     n.Type,
		Value:    n.Value,
		Children: make([]*Node, len(n.Children)),
	}
	for i, child := range n.Children {
		clone.Children[i] = child.Clone()
	}
	return clone
}

// Tree represents the genome for Genetic Programming.
type Tree struct {
	Root *Node
}

// Clone returns a deep copy of the Tree.
func (t Tree) Clone() Tree {
	return Tree{Root: t.Root.Clone()}
}

// PrimitiveSet defines the available functions and terminals.
type PrimitiveSet struct {
	Functions []string
	Terminals []string
	// Arity defines how many children each function takes
	Arity map[string]int
}

func NewPrimitiveSet() *PrimitiveSet {
	return &PrimitiveSet{
		Arity: make(map[string]int),
	}
}

func (p *PrimitiveSet) AddFunc(name string, arity int) {
	p.Functions = append(p.Functions, name)
	p.Arity[name] = arity
}

func (p *PrimitiveSet) AddTerm(name string) {
	p.Terminals = append(p.Terminals, name)
}

// GenerateFull creates a tree where all branches reach the exact maxDepth.
func GenerateFull(rng *rand.Rand, maxDepth int, pset PrimitiveSet) *Node {
	if maxDepth == 0 || len(pset.Functions) == 0 {
		return &Node{
			Type:  TerminalNode,
			Value: pset.Terminals[rng.IntN(len(pset.Terminals))],
		}
	}

	funcName := pset.Functions[rng.IntN(len(pset.Functions))]
	arity := pset.Arity[funcName]

	node := &Node{
		Type:     FunctionNode,
		Value:    funcName,
		Children: make([]*Node, arity),
	}

	for i := 0; i < arity; i++ {
		node.Children[i] = GenerateFull(rng, maxDepth-1, pset)
	}

	return node
}

// GenerateGrow creates a tree where branches may end before maxDepth.
func GenerateGrow(rng *rand.Rand, maxDepth int, pset PrimitiveSet) *Node {
	if maxDepth == 0 || len(pset.Functions) == 0 {
		return &Node{
			Type:  TerminalNode,
			Value: pset.Terminals[rng.IntN(len(pset.Terminals))],
		}
	}

	// Choose randomly between a function and a terminal
	isTerminal := rng.IntN(len(pset.Functions)+len(pset.Terminals)) >= len(pset.Functions)

	if isTerminal {
		return &Node{
			Type:  TerminalNode,
			Value: pset.Terminals[rng.IntN(len(pset.Terminals))],
		}
	}

	funcName := pset.Functions[rng.IntN(len(pset.Functions))]
	arity := pset.Arity[funcName]

	node := &Node{
		Type:     FunctionNode,
		Value:    funcName,
		Children: make([]*Node, arity),
	}

	for i := 0; i < arity; i++ {
		node.Children[i] = GenerateGrow(rng, maxDepth-1, pset)
	}

	return node
}
