package types

// BitString represents a genome as a slice of booleans.
type BitString []bool

// FloatVector represents a genome as a slice of floating-point numbers.
type FloatVector []float64

// Permutation represents a genome as a slice of integers, usually forming a permutation of 0 to N-1.
type Permutation []int

// Clone functions for convenience

func (b BitString) Clone() BitString {
	c := make(BitString, len(b))
	copy(c, b)
	return c
}

func (f FloatVector) Clone() FloatVector {
	c := make(FloatVector, len(f))
	copy(c, f)
	return c
}

func (p Permutation) Clone() Permutation {
	c := make(Permutation, len(p))
	copy(c, p)
	return c
}
