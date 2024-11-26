package main

import (
	"fmt"
	"strings"
)

// Node represents a node in the syntax tree of the logical expression.
type Node struct {
	Value string
	Left  *Node
	Right *Node
}

// parseExpression converts a propositional logic string into a syntax tree.
func parseExpression(expr string) *Node {
	expr = strings.TrimSpace(expr)
	// Handling parentheses
	if expr[0] == '(' && expr[len(expr)-1] == ')' {
		return parseExpression(expr[1 : len(expr)-1])
	}
	// Splitting logical operators
	operators := []string{"->", "<->", "|", "&", "!"}
	for _, op := range operators {
		pos := strings.LastIndex(expr, op)
		if pos != -1 {
			left := parseExpression(expr[:pos])
			right := parseExpression(expr[pos+len(op):])
			return &Node{Value: op, Left: left, Right: right}
		}
	}
	// Leaf node
	return &Node{Value: expr}
}

// eliminateImplications removes implications and equivalences.
func eliminateImplications(node *Node) *Node {
	if node == nil {
		return nil
	}
	node.Left = eliminateImplications(node.Left)
	node.Right = eliminateImplications(node.Right)
	switch node.Value {
	case "->": // A -> B ≡ !A | B
		node.Value = "|"
		node.Left = &Node{Value: "!", Left: eliminateImplications(node.Left)}
	case "<->": // A <-> B ≡ (!A | B) & (!B | A)
		a := eliminateImplications(node.Left)
		b := eliminateImplications(node.Right)
		node.Value = "&"
		node.Left = &Node{Value: "|", Left: &Node{Value: "!", Left: a}, Right: b}
		node.Right = &Node{Value: "|", Left: &Node{Value: "!", Left: b}, Right: a}
	}
	return node
}

// pushNegations applies De Morgan's laws to push negations inwards.
func pushNegations(node *Node) *Node {
	if node == nil {
		return nil
	}
	switch node.Value {
	case "!":
		if node.Left != nil && (node.Left.Value == "&" || node.Left.Value == "|") {
			// Apply De Morgan's Laws
			op := "|"
			if node.Left.Value == "&" {
				op = "&"
			}
			node = &Node{
				Value: op,
				Left:  pushNegations(&Node{Value: "!", Left: node.Left.Left}),
				Right: pushNegations(&Node{Value: "!", Left: node.Left.Right}),
			}
		} else if node.Left != nil && node.Left.Value == "!" {
			// Double negation elimination
			node = pushNegations(node.Left.Left)
		}
	default:
		node.Left = pushNegations(node.Left)
		node.Right = pushNegations(node.Right)
	}
	return node
}

// distributeOr distributes OR over AND to achieve CNF.
func distributeOr(node *Node) *Node {
	if node == nil {
		return nil
	}
	node.Left = distributeOr(node.Left)
	node.Right = distributeOr(node.Right)
	if node.Value == "|" {
		if node.Left != nil && node.Left.Value == "&" {
			return &Node{
				Value: "&",
				Left:  distributeOr(&Node{Value: "|", Left: node.Left.Left, Right: node.Right}),
				Right: distributeOr(&Node{Value: "|", Left: node.Left.Right, Right: node.Right}),
			}
		}
		if node.Right != nil && node.Right.Value == "&" {
			return &Node{
				Value: "&",
				Left:  distributeOr(&Node{Value: "|", Left: node.Left, Right: node.Right.Left}),
				Right: distributeOr(&Node{Value: "|", Left: node.Left, Right: node.Right.Right}),
			}
		}
	}
	return node
}

// toCNF converts a syntax tree to CNF.
func toCNF(node *Node) *Node {
	node = eliminateImplications(node)
	node = pushNegations(node)
	node = distributeOr(node)
	return node
}

// printExpression converts a syntax tree back to a string representation.
func printExpression(node *Node) string {
	if node == nil {
		return ""
	}
	if node.Left == nil && node.Right == nil {
		return node.Value
	}
	if node.Right == nil {
		return fmt.Sprintf("!(%s)", printExpression(node.Left))
	}
	return fmt.Sprintf("(%s %s %s)", printExpression(node.Left), node.Value, printExpression(node.Right))
}

// Main function
func main() {
	// Example input
	expression := "(A -> B) & (C | D) & (E -> F) & (G | H) | (I -> J) & (K | L)"
	fmt.Println("Original Expression:", expression)

	// Parse the expression into a syntax tree
	root := parseExpression(expression)

	// Convert to CNF
	cnfRoot := toCNF(root)

	// Print the CNF expression
	cnfExpression := printExpression(cnfRoot)
	fmt.Println("CNF Expression:", cnfExpression)
}
