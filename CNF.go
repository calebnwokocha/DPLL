package main

import (
	"fmt"
	"strings"
)

// Define types for the propositional logic formulas
type Formula interface {
	ToCNF() string
}

type Variable string
type Conjunction struct {
	Left  Formula
	Right Formula
}
type Disjunction struct {
	Left  Formula
	Right Formula
}
type Negation struct {
	Operand Formula
}
type Implication struct {
	Left  Formula
	Right Formula
}
type Equivalence struct {
	Left  Formula
	Right Formula
}
type XOR struct {
	Left  Formula
	Right Formula
}

// ToCNF method for Variable
func (v Variable) ToCNF() string {
	return string(v)
}

// ToCNF method for Conjunction
func (c *Conjunction) ToCNF() string {
	return fmt.Sprintf("(%s ^ %s)", c.Left.ToCNF(), c.Right.ToCNF())
}

// ToCNF method for Disjunction with Switching Variables
func (d *Disjunction) ToCNF() string {
	// Recursively convert both parts
	leftCNF := d.Left.ToCNF()
	rightCNF := d.Right.ToCNF()

	// Apply the switching variable technique if both sides are conjunctions
	if strings.Contains(leftCNF, "^") && strings.Contains(rightCNF, "^") {
		// Introduce a new variable 'Z' (can be any fresh name)
		Z := "Z"
		// Use the switching variable technique (Z -> P) ^ (~Z -> Q)
		return fmt.Sprintf("(%s -> %s) ^ (~%s -> %s)", Z, leftCNF, Z, rightCNF)
	}

	// If not, just return the disjunction
	return fmt.Sprintf("(%s v %s)", leftCNF, rightCNF)
}

// ToCNF method for Negation
func (n *Negation) ToCNF() string {
	operandCNF := n.Operand.ToCNF()

	// Apply De Morgan's law
	if strings.HasPrefix(operandCNF, "~") {
		// Double negation
		return operandCNF[1:]
	}

	// Apply De Morgan's law: ~(P ^ Q) -> ~P v ~Q and ~(P v Q) -> ~P ^ ~Q
	if strings.Contains(operandCNF, "^") {
		// ~(P ^ Q) -> ~P v ~Q
		parts := strings.Split(operandCNF, "^")
		return fmt.Sprintf("~(%s) v ~(%s)", strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	} else if strings.Contains(operandCNF, "v") {
		// ~(P v Q) -> ~P ^ ~Q
		parts := strings.Split(operandCNF, "v")
		return fmt.Sprintf("~(%s) ^ ~(%s)", strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}

	// If it's just a variable, return the negation
	return "~" + operandCNF
}

// ToCNF method for Implication
func (i *Implication) ToCNF() string {
	// P -> Q is equivalent to ~P v Q
	return fmt.Sprintf("(~%s v %s)", i.Left.ToCNF(), i.Right.ToCNF())
}

// ToCNF method for Equivalence
func (e *Equivalence) ToCNF() string {
	// P <-> Q is equivalent to (P ^ Q) v (~P ^ ~Q)
	return fmt.Sprintf("(%s ^ %s) v (~%s ^ ~%s)", e.Left.ToCNF(), e.Right.ToCNF(), e.Left.ToCNF(), e.Right.ToCNF())
}

// ToCNF method for XOR
func (x *XOR) ToCNF() string {
	// P xor Q is equivalent to (P ^ ~Q) v (~P ^ Q)
	return fmt.Sprintf("(%s ^ ~%s) v (~%s ^ %s)", x.Left.ToCNF(), x.Right.ToCNF(), x.Left.ToCNF(), x.Right.ToCNF())
}

func main() {
	// Example 1: (P1 ^ P2) v (Q1 ^ Q2)
	formula1 := &Disjunction{
		Left: &Conjunction{
			Left:  Variable("P1"),
			Right: Variable("P2"),
		},
		Right: &Conjunction{
			Left:  Variable("Q1"),
			Right: Variable("Q2"),
		},
	}

	// Example 2: (A ^ B) v C
	formula2 := &Disjunction{
		Left: &Conjunction{
			Left:  Variable("A"),
			Right: Variable("B"),
		},
		Right: Variable("C"),
	}

	// Example 3: (P1 -> P2) ^ (Q1 -> Q2)
	formula3 := &Conjunction{
		Left: &Implication{
			Left:  Variable("P1"),
			Right: Variable("P2"),
		},
		Right: &Implication{
			Left:  Variable("Q1"),
			Right: Variable("Q2"),
		},
	}

	// Example 4: ~(P ^ Q)
	formula4 := &Negation{
		Operand: &Conjunction{
			Left:  Variable("P"),
			Right: Variable("Q"),
			},
	}

	// Print CNF conversions for the examples
	fmt.Println("Example 1 CNF:", formula1.ToCNF())
	fmt.Println("Example 2 CNF:", formula2.ToCNF())
	fmt.Println("Example 3 CNF:", formula3.ToCNF())
	fmt.Println("Example 4 CNF:", formula4.ToCNF())
}
