package maze_test

import (
	"fmt"
	"log"

	"github.com/vinser/maze"
)

// ExampleNew demonstrates how to create a new maze.
// It shows a simple case without a den.
func ExampleNew() {
	// Create a new 15x7 maze.
	m, err := maze.New(15, 7, 0, 0)
	if err != nil {
		log.Fatalf("failed to create maze: %v", err)
	}

	// Generate the maze paths with a fixed seed for reproducibility.
	err = m.Generate(42, nil, nil, nil, "", 0.5)
	if err != nil {
		log.Fatalf("failed to generate maze: %v", err)
	}

	// Print the maze structure.
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			cell, _ := m.Cell(x, y)
			fmt.Printf("%c", cell)
		}
		fmt.Println()
	}

	// Output:
	// ███████████████
	// █         █   █
	// ███ ███ ███ █ █
	// █   █S  █   █ █
	// █ ███████ ███ █
	// █         █E  █
	// ███████████████
}

// ExampleMaze_Solve demonstrates how to solve a generated maze and display the solution.
func ExampleMaze_Solve() {
	// Create a new 15x7 maze.
	m, err := maze.New(15, 7, 0, 0)
	if err != nil {
		log.Fatalf("failed to create maze: %v", err)
	}

	// Generate the maze with a fixed seed for a predictable structure.
	err = m.Generate(42, nil, nil, nil, "", 0.5)
	if err != nil {
		log.Fatalf("failed to generate maze: %v", err)
	}

	// Solve the maze to get the solution path.
	solution, found := m.Solve()
	if !found {
		fmt.Println("No solution found.")
		return
	}

	// Create a map of solution points for easy lookup when rendering.
	solutionPoints := make(map[maze.Point]bool)
	for _, p := range solution {
		solutionPoints[p] = true
	}

	// Print the maze, overlaying the solution path.
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			p := maze.Point{X: x, Y: y}
			cell, _ := m.Cell(x, y)

			// If the point is part of the solution and not the start/end, draw it.
			if solutionPoints[p] && cell != maze.Start && cell != maze.End {
				fmt.Printf("%c", maze.SolutionPath)
			} else {
				fmt.Printf("%c", cell)
			}
		}
		fmt.Println()
	}

	// Output:
	// ███████████████
	// █  .....  █...█
	// ███.███.███.█.█
	// █...█S..█...█.█
	// █.███████.███.█
	// █.........█E..█
	// ███████████████
}
