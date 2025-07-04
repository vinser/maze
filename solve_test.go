package maze

import (
	"testing"
)

func TestSolve(t *testing.T) {
	t.Run("Simple solvable maze", func(t *testing.T) {
		m := &Maze{
			width:  5,
			height: 3,
			grid: [][]Cell{
				{Wall, Wall, Wall, Wall, Wall},
				{Wall, Start, Path, End, Wall},
				{Wall, Wall, Wall, Wall, Wall},
			},
			start: Point{X: 1, Y: 1},
			end:   Point{X: 3, Y: 1},
		}

		path, found := m.Solve()
		if !found {
			t.Fatal("Expected to find a path, but did not")
		}
		expectedPath := []Point{{1, 1}, {2, 1}, {3, 1}}
		if len(path) != len(expectedPath) {
			t.Fatalf("Expected path length of %d, got %d", len(expectedPath), len(path))
		}
		for i, p := range expectedPath {
			if path[i] != p {
				t.Errorf("Path point %d is incorrect. Expected %+v, got %+v", i, p, path[i])
			}
		}
	})

	t.Run("Unsolvable maze", func(t *testing.T) {
		m := &Maze{
			width:  5,
			height: 3,
			grid: [][]Cell{
				{Wall, Wall, Wall, Wall, Wall},
				{Wall, Start, Wall, End, Wall},
				{Wall, Wall, Wall, Wall, Wall},
			},
			start: Point{X: 1, Y: 1},
			end:   Point{X: 3, Y: 1},
		}

		path, found := m.Solve()
		if found {
			t.Error("Expected not to find a path, but one was found")
		}
		if path != nil {
			t.Errorf("Expected nil path for unsolvable maze, got %+v", path)
		}
	})

	t.Run("Start equals End", func(t *testing.T) {
		m := &Maze{
			width:  3,
			height: 3,
			grid: [][]Cell{
				{Wall, Wall, Wall},
				{Wall, Start, Wall},
				{Wall, Wall, Wall},
			},
			start: Point{X: 1, Y: 1},
			end:   Point{X: 1, Y: 1},
		}

		path, found := m.Solve()
		if !found {
			t.Fatal("Expected to find a path, but did not")
		}
		if len(path) != 1 {
			t.Fatalf("Expected path length of 1, got %d", len(path))
		}
		if path[0] != m.start {
			t.Errorf("Expected path to contain only the start point, got %+v", path)
		}
	})
}
