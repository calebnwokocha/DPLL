package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Clause []int // A clause is a slice of integers representing literals
type CNF []Clause // CNF is a conjunction of clauses

// UnitPropagation simplifies the CNF by assigning values for unit clauses
func UnitPropagation(cnf CNF, assignment map[int]bool) (CNF, bool) {
	for {
		unitFound := false
		for _, clause := range cnf {
			if len(clause) == 1 { // Found a unit clause
				unit := clause[0]
				unitFound = true
				value := unit > 0
				variable := abs(unit)
				assignment[variable] = value
				cnf = assign(cnf, variable, value)
				break
			}
		}
		if !unitFound {
			break
		}
	}
	for _, clause := range cnf {
		if len(clause) == 0 {
			return cnf, false // Conflict detected
		}
	}
	return cnf, true
}

// PureLiteralElimination simplifies CNF by assigning values for pure literals
func PureLiteralElimination(cnf CNF, assignment map[int]bool) CNF {
	literalCount := make(map[int]int)
	for _, clause := range cnf {
		for _, literal := range clause {
			literalCount[literal]++
		}
	}
	for literal, count := range literalCount {
		if count > 0 && literalCount[-literal] == 0 { // Pure literal found
			value := literal > 0
			variable := abs(literal)
			assignment[variable] = value
			cnf = assign(cnf, variable, value)
		}
	}
	return cnf
}

// Assign simplifies the CNF given a variable assignment
func assign(cnf CNF, variable int, value bool) CNF {
	newCNF := CNF{}
	for _, clause := range cnf {
		newClause := Clause{}
		skipClause := false
		for _, literal := range clause {
			if literal == variable && value || literal == -variable && !value {
				skipClause = true
				break
			} else if literal != variable && literal != -variable {
				newClause = append(newClause, literal)
			}
		}
		if !skipClause {
			newCNF = append(newCNF, newClause)
		}
	}
	return newCNF
}

// DPLL implements the main algorithm
func DPLL(cnf CNF, assignment map[int]bool) bool {
	// Apply unit propagation
	cnf, ok := UnitPropagation(cnf, assignment)
	if !ok {
		return false // Conflict detected
	}

	// Apply pure literal elimination
	cnf = PureLiteralElimination(cnf, assignment)

	// Check if all clauses are satisfied
	if len(cnf) == 0 {
		return true // Satisfiable
	}

	// Select the next variable to assign (heuristic: first literal in the first clause)
	var variable int
	for _, clause := range cnf {
		if len(clause) > 0 {
			variable = abs(clause[0])
			break
		}
	}

	// Try assigning true
	assignment[variable] = true
	if DPLL(assign(cnf, variable, true), assignment) {
		return true
	}

	// Backtrack and try assigning false
	assignment[variable] = false
	return DPLL(assign(cnf, variable, false), assignment)
}

// Helper function: absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// CompleteAssignment ensures all variables have an assignment
func CompleteAssignment(cnf CNF, assignment map[int]bool) map[int]bool {
	for _, clause := range cnf {
		for _, literal := range clause {
			variable := abs(literal)
			if _, exists := assignment[variable]; !exists {
				assignment[variable] = true // Default arbitrary assignment
			}
		}
	}
	return assignment
}

// ParseCNF parses user input into a CNF
func ParseCNF(input string) CNF {
	cnf := CNF{}
	clauses := strings.Split(input, " AND ")
	for _, clause := range clauses {
		clause = strings.Trim(clause, "() ")
		literals := strings.Split(clause, " OR ")
		c := Clause{}
		for _, literal := range literals {
			num, _ := strconv.Atoi(strings.TrimSpace(literal))
			c = append(c, num)
		}
		cnf = append(cnf, c)
	}
	return cnf
}

// ValidateCNF ensures the formula is in correct CNF format
func ValidateCNF(input string) bool {
	clauses := strings.Split(input, " AND ")
	for _, clause := range clauses {
		clause = strings.TrimSpace(clause)
		if len(clause) < 2 || clause[0] != '(' || clause[len(clause)-1] != ')' {
			return false
		}
		literals := strings.Split(clause[1:len(clause)-1], " OR ")
		if len(literals) == 0 {
			return false
		}
		for _, literal := range literals {
			literal = strings.TrimSpace(literal)
			if _, err := strconv.Atoi(literal); err != nil {
				return false
			}
		}
	}
	return true
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the Interactive DPLL SAT Solver")
	fmt.Println("Input your CNF formula using the format: (1 OR -2) AND (-1 OR 3) AND (2 OR -3)")
	fmt.Println("Type 'exit' to quit the program.")

	for {
		fmt.Print("\nEnter your formula: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Check for exit condition
		if strings.ToLower(input) == "exit" {
			fmt.Println("Exiting the program. Goodbye!")
			break
		}

		// Validate input
		if !ValidateCNF(input) {
			fmt.Println("Invalid CNF format. Please use the format: (literal1 OR literal2) AND (literal3 OR ... )")
			continue
		}

		// Parse input into CNF
		cnf := ParseCNF(input)

		// Solve using DPLL
		assignment := make(map[int]bool)
		if DPLL(cnf, assignment) {
			assignment = CompleteAssignment(cnf, assignment)
			fmt.Println("SATISFIABLE with assignment:", assignment)
		} else {
			fmt.Println("UNSATISFIABLE")
		}
	}
}
