package maze_test

import (
	"testing"

	"github.com/vinser/maze"
)

func TestGenerate(t *testing.T) {
	testCases := []struct {
		name       string
		width      int
		height     int
		denWidth   int
		denHeight  int
		startPoint func(m *maze.Maze) *maze.Point // Use a function to get context-aware points
		endPoint   func(m *maze.Maze) *maze.Point
		doorSide   string
		expectErr  bool
		postCheck  func(t *testing.T, m *maze.Maze)
	}{
		{
			name:      "Successful generation produces solvable maze",
			width:     21,
			height:    21,
			denWidth:  5,
			denHeight: 5,
			expectErr: false,
			postCheck: func(t *testing.T, m *maze.Maze) {
				startCell, _ := m.Cell(m.Start().X, m.Start().Y)
				if startCell != maze.Start {
					t.Errorf("Expected start cell to be 'S', got '%c'", startCell)
				}
				endCell, _ := m.Cell(m.End().X, m.End().Y)
				if endCell != maze.End {
					t.Errorf("Expected end cell to be 'E', got '%c'", endCell)
				}
				_, found := m.Solve()
				if !found {
					t.Error("Generated maze should be solvable, but it is not")
				}
			},
		},
		{
			name:      "Invalid start point inside den",
			width:     21,
			height:    21,
			denWidth:  5,
			denHeight: 5,
			startPoint: func(m *maze.Maze) *maze.Point {
				// A point guaranteed to be inside the den
				return &maze.Point{X: m.DenStartX() + 1, Y: m.DenStartY() + 1}
			},
			expectErr: true,
		},
		{
			name:   "Successful generation with specified start and end",
			width:  21,
			height: 21,
			startPoint: func(m *maze.Maze) *maze.Point {
				return &maze.Point{X: 1, Y: 1}
			},
			endPoint: func(m *maze.Maze) *maze.Point {
				return &maze.Point{X: 19, Y: 19}
			},
			expectErr: false,
			postCheck: func(t *testing.T, m *maze.Maze) {
				if m.Start() != (maze.Point{X: 1, Y: 1}) {
					t.Errorf("Expected start point to be {1, 1}, got %+v", m.Start())
				}
				if m.End() != (maze.Point{X: 19, Y: 19}) {
					t.Errorf("Expected end point to be {19, 19}, got %+v", m.End())
				}
				_, found := m.Solve()
				if !found {
					t.Error("Generated maze with specified start/end should be solvable")
				}
			},
		},
		{
			name:   "Start and end points are the same",
			width:  21,
			height: 21,
			startPoint: func(m *maze.Maze) *maze.Point {
				return &maze.Point{X: 5, Y: 5}
			},
			endPoint: func(m *maze.Maze) *maze.Point {
				return &maze.Point{X: 5, Y: 5}
			},
			expectErr: true,
		},
		{
			name:      "Guaranteed door creation on top side",
			width:     21,
			height:    21,
			denWidth:  5,
			denHeight: 5,
			doorSide:  "top",
			expectErr: false,
			postCheck: func(t *testing.T, m *maze.Maze) {
				door := m.Door()
				doorCell, _ := m.Cell(door.X, door.Y)
				if doorCell != maze.Path {
					t.Fatalf("Expected door cell to be a Path, got '%c'", doorCell)
				}
				// Check that it connects the den to the maze
				p1 := maze.Point{X: door.X, Y: door.Y - 1} // Outside den
				p2 := maze.Point{X: door.X, Y: door.Y + 1} // Inside den
				cell1, _ := m.Cell(p1.X, p1.Y)
				cell2, _ := m.Cell(p2.X, p2.Y)
				if cell1 != maze.Path {
					t.Errorf("Expected cell outside door to be a Path, but it's '%c'", cell1)
				}
				if cell2 != maze.Path {
					t.Errorf("Expected cell inside door to be a Path, but it's '%c'", cell2)
				}
				if !m.IsInsideDen(p2) {
					t.Errorf("Expected cell inside door to be in the den, but it's not")
				}
			},
		},
		{
			name:      "Random door creation",
			width:     21,
			height:    21,
			denWidth:  5,
			denHeight: 5,
			expectErr: false,
			postCheck: func(t *testing.T, m *maze.Maze) {
				door := m.Door()
				if door.X == 0 && door.Y == 0 {
					t.Fatal("Expected a door to be created, but door point is zero")
				}
				doorCell, _ := m.Cell(door.X, door.Y)
				if doorCell != maze.Path {
					t.Errorf("Expected door cell to be a Path, got '%c'", doorCell)
				}
				// A valid door must connect a den cell to a maze path cell
				p1_h := maze.Point{X: door.X - 1, Y: door.Y}
				p2_h := maze.Point{X: door.X + 1, Y: door.Y}
				isHorizDoor := m.IsInsideDen(p1_h) != m.IsInsideDen(p2_h)
				p1_v := maze.Point{X: door.X, Y: door.Y - 1}
				p2_v := maze.Point{X: door.X, Y: door.Y + 1}
				isVertDoor := m.IsInsideDen(p1_v) != m.IsInsideDen(p2_v)
				if !isHorizDoor && !isVertDoor {
					t.Errorf("Door at %+v does not connect den to maze path", door)
				}
			},
		},
		{
			name:      "Door side error - too close to edge",
			width:     7,
			height:    7,
			denWidth:  3,
			denHeight: 3,
			doorSide:  "top",
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := maze.New(tc.width, tc.height, tc.denWidth, tc.denHeight)
			if err != nil {
				t.Fatalf("Failed to create maze for test case: %v", err)
			}

			var start *maze.Point
			if tc.startPoint != nil {
				start = tc.startPoint(m)
			}
			var end *maze.Point
			if tc.endPoint != nil {
				end = tc.endPoint(m)
			}

			// Use a fixed seed for reproducibility
			err = m.Generate(1, start, end, nil, tc.doorSide, 0.5)

			if tc.expectErr {
				if err == nil {
					t.Errorf("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("Expected no error, but got: %v", err)
				}
				if tc.postCheck != nil {
					tc.postCheck(t, m)
				}
			}
		})
	}
}
