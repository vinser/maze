package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/vinser/maze" // The maze generation library.
)

func main() {
	// Default size, can be overridden by command-line arguments
	width := flag.Int("width", 41, "The width of the maze")
	height := flag.Int("height", 21, "The height of the maze")
	denWidth := flag.Int("denWidth", 0, "The width of the central den. Set to 0 for no den.")
	denHeight := flag.Int("denHeight", 0, "The height of the central den. Set to 0 for no den.")
	seed := flag.Int64("seed", 0, "Seed for the random number generator. If 0, uses current time.")
	startX := flag.Int("startX", 0, "The X coordinate for the generation start point. If 0, a random point is chosen.")
	startY := flag.Int("startY", 0, "The Y coordinate for the generation start point. If 0, a random point is chosen.")
	doorX := flag.Int("doorX", 0, "The X coordinate for the den door. If 0, a random door is chosen.")
	doorY := flag.Int("doorY", 0, "The Y coordinate for the den door. If 0, a random door is chosen.")
	doorSide := flag.String("doorSide", "", "Side for the den door (top, bottom, left, right). Overrides --doorX/Y.")
	bias := flag.Float64("bias", 0.5, "Bias for straight corridors (0.0 to 1.0). 0 is random, 1 always goes straight if possible.")
	solveRatio := flag.Float64("solveRatio", -1.0, "The fraction of the solution path to display (0.0 to 1.0). If not set, maze is not solved.")
	flag.Parse()

	// Create a new maze instance
	m, err := maze.New(*width, *height, *denWidth, *denHeight)
	if err != nil {
		log.Fatalf("Error creating maze: %v", err)
	}

	// Prepare parameters for generation
	genSeed := *seed
	if genSeed == 0 {
		genSeed = time.Now().UnixNano()
	}

	var startPoint *maze.Point
	if *startX > 0 && *startY > 0 {
		startPoint = &maze.Point{X: *startX, Y: *startY}
	}

	var doorPoint *maze.Point
	// doorSide takes precedence over doorX/Y
	if *doorSide == "" && *doorX > 0 && *doorY > 0 {
		doorPoint = &maze.Point{X: *doorX, Y: *doorY}
	}

	// Generate the maze paths
	if err := m.Generate(genSeed, startPoint, doorPoint, *doorSide, *bias); err != nil {
		log.Fatalf("Error generating maze: %v", err)
	}

	var solutionPath []maze.Point
	// If the solveRatio flag is set, solve the maze
	if *solveRatio >= 0.0 {
		if *solveRatio > 1.0 {
			log.Fatalf("solveRatio must be between 0.0 and 1.0")
		}
		path, found := m.Solve()
		if !found {
			fmt.Println("No solution could be found for the maze.")
		} else {
			solutionPath = path
		}
	}

	// Print the generated maze to the console
	fmt.Println(renderMaze(m, solutionPath, *solveRatio))
}

// renderMaze builds the string representation of the maze.
// It takes the maze structure and overlays the solution path based on the ratio.
func renderMaze(m *maze.Maze, path []maze.Point, ratio float64) string {
	// Create a map of solution points for quick lookup.
	solutionPoints := make(map[maze.Point]bool)
	if path != nil && ratio >= 0.0 {
		pathLength := len(path)
		if pathLength > 0 {
			// Determine how many points of the path to show, starting after the 'S'.
			// math.Ceil ensures that for any ratio > 0, at least one step is shown.
			pointsToShow := int(math.Ceil(float64(pathLength-1) * ratio))
			for i := 1; i <= pointsToShow && i < pathLength; i++ {
				solutionPoints[path[i]] = true
			}
		}
	}

	var sb strings.Builder
	for y := 0; y < m.Height(); y++ {
		for x := 0; x < m.Width(); x++ {
			p := maze.Point{X: x, Y: y}
			cell, _ := m.Cell(x, y)
			// If the point is part of the solution path and is a normal path cell, draw it.
			if solutionPoints[p] && cell == maze.Path {
				sb.WriteRune(rune(maze.SolutionPath))
			} else {
				sb.WriteRune(rune(cell))
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
