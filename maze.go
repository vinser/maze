package maze

import (
	"fmt"
)

// Cell represents the type of a single cell in the maze.
type Cell rune

const (
	// Wall is an impassable cell.
	Wall Cell = '█'
	// Path is a walkable cell.
	Path Cell = ' '
	// Start is the starting point of the maze.
	Start Cell = 'S'
	// End is the ending point of the maze.
	End Cell = 'E'
	// SolutionPath is a cell on the solved path.
	SolutionPath Cell = '.'
)

// Point represents a coordinate in the maze.
type Point struct {
	X, Y int
}

// Maze represents the maze structure.
type Maze struct {
	width  int
	height int
	grid   [][]Cell
	start  Point
	end    Point
	door   Point

	// den dimensions
	denWidth  int
	denHeight int
	denStartX int
	denStartY int
}

// adjustToOdd ensures a dimension is odd by incrementing it if it's even and positive.
func adjustToOdd(dim int) int {
	if dim > 0 && dim%2 == 0 {
		return dim + 1
	}
	return dim
}

// validateAndAdjustDimensions checks and adjusts maze and den dimensions.
func validateAndAdjustDimensions(width, height, denWidth, denHeight int) (int, int, int, int, error) {
	if width <= 0 || height <= 0 {
		return 0, 0, 0, 0, fmt.Errorf("width and height must be positive")
	}
	if denWidth < 0 || denHeight < 0 {
		return 0, 0, 0, 0, fmt.Errorf("den dimensions must be non-negative")
	}

	// Adjust dimensions to be odd for proper maze structure.
	adjWidth := adjustToOdd(width)
	adjHeight := adjustToOdd(height)
	adjDenWidth := adjustToOdd(denWidth)
	adjDenHeight := adjustToOdd(denHeight)

	// Ensure den fits within the maze walls (a 1-cell border on each side).
	if adjDenWidth > 0 && adjDenWidth >= adjWidth-2 {
		return 0, 0, 0, 0, fmt.Errorf("den width (%d) is too large for the maze width (%d)", adjDenWidth, adjWidth)
	}
	if adjDenHeight > 0 && adjDenHeight >= adjHeight-2 {
		return 0, 0, 0, 0, fmt.Errorf("den height (%d) is too large for the maze height (%d)", adjDenHeight, adjHeight)
	}

	return adjWidth, adjHeight, adjDenWidth, adjDenHeight, nil
}

// calculateDenPosition determines the top-left corner of the den, ensuring grid alignment.
func calculateDenPosition(width, denWidth, height, denHeight int) (int, int) {
	// Center the den.
	denStartX := (width - denWidth) / 2
	denStartY := (height - denHeight) / 2

	// Adjust for grid alignment to prevent "double walls".
	// Den boundaries should be on odd coordinates to align with maze paths.
	if denStartX > 0 && denStartX%2 == 0 {
		denStartX--
	}
	if denStartY > 0 && denStartY%2 == 0 {
		denStartY--
	}
	return denStartX, denStartY
}

// New creates a new maze with the given width and height.
// To ensure a border and clear pathways, dimensions will be adjusted to be odd.
func New(width, height, denWidth, denHeight int) (*Maze, error) {
	adjWidth, adjHeight, adjDenWidth, adjDenHeight, err := validateAndAdjustDimensions(width, height, denWidth, denHeight)
	if err != nil {
		return nil, err
	}

	m := &Maze{
		width:     adjWidth,
		height:    adjHeight,
		denWidth:  adjDenWidth,
		denHeight: adjDenHeight,
	}

	if m.denWidth > 0 && m.denHeight > 0 {
		m.denStartX, m.denStartY = calculateDenPosition(m.width, m.denWidth, m.height, m.denHeight)
	}

	m.initializeGrid()

	return m, nil
}

// initializeGrid creates the grid and carves out the den area.
func (m *Maze) initializeGrid() {
	m.grid = make([][]Cell, m.height)
	for i := range m.grid {
		m.grid[i] = make([]Cell, m.width)
		for j := range m.grid[i] {
			// Pre-carve the den by setting its area to Path, otherwise set to Wall.
			if m.IsInsideDen(Point{X: j, Y: i}) {
				m.grid[i][j] = Path
			} else {
				m.grid[i][j] = Wall
			}
		}
	}

}

// IsInsideDen checks if a given point is within the boundaries of the central den.
func (m *Maze) IsInsideDen(p Point) bool {
	if m.denWidth <= 0 || m.denHeight <= 0 {
		return false
	}
	return p.X >= m.denStartX && p.X < m.denStartX+m.denWidth &&
		p.Y >= m.denStartY && p.Y < m.denStartY+m.denHeight
}

// IsAdjacentToDen checks if a point is directly next to a den cell, but not inside it.
func (m *Maze) IsAdjacentToDen(p Point) bool {
	if m.denWidth <= 0 || m.denHeight <= 0 {
		return false
	}
	// A point inside the den is not considered "adjacent".
	if m.IsInsideDen(p) {
		return false
	}

	// Check the four cardinal neighbors of the point.
	for _, dir := range []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
		if m.IsInsideDen(Point{X: p.X + dir.X, Y: p.Y + dir.Y}) {
			return true
		}
	}
	return false
}

// Start returns the maze's starting point.
func (m *Maze) Start() Point {
	return m.start
}

// End returns the maze's ending point.
func (m *Maze) End() Point {
	return m.end
}

// Width returns the maze's width.
func (m *Maze) Width() int {
	return m.width
}

// Height returns the maze's height.
func (m *Maze) Height() int {
	return m.height
}

// DenWidth returns the maze's den width.
func (m *Maze) DenWidth() int {
	return m.denWidth
}

// DenHeight returns the maze's den height.
func (m *Maze) DenHeight() int {
	return m.denHeight
}

// DenStartX returns the maze's den startX
func (m *Maze) DenStartX() int {
	return m.denStartX
}

// DenStartY returns the maze's den startY
func (m *Maze) DenStartY() int {
	return m.denStartY
}

// Door returns the maze's den door point.
func (m *Maze) Door() Point {
	return m.door
}

// Cell returns the cell type at a given coordinate.
// It returns the cell and true if the point is within bounds, otherwise it returns a zero value and false.
func (m *Maze) Cell(x, y int) (Cell, bool) {
	if x < 0 || x >= m.width || y < 0 || y >= m.height {
		return 0, false
	}
	return m.grid[y][x], true
}
