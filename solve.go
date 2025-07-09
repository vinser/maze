package maze

// Solve finds the shortest path from Start to End using Breadth-First Search (BFS)
// and returns it as a slice of points.
// It returns the path and true if a path is found, otherwise it returns nil and false.
func (m *Maze) Solve() ([]Point, bool) {
	// Queue for BFS, starting with the start point
	queue := []Point{m.start}

	// visited grid to prevent cycles and redundant checks
	visited := make([][]bool, m.height)
	for i := range visited {
		visited[i] = make([]bool, m.width)
	}
	visited[m.start.Y][m.start.X] = true

	// parent map to reconstruct the path
	parent := make(map[Point]Point)

	pathFound := false
	head := 0
	for head < len(queue) {
		// Dequeue the current point
		current := queue[head]
		head++

		// If we reached the end, stop searching
		if current == m.end {
			pathFound = true
			break
		}

		// Explore neighbors (Up, Down, Left, Right)
		for _, dir := range []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
			next := Point{X: current.X + dir.X, Y: current.Y + dir.Y}

			// Check if the neighbor is within bounds
			if next.X < 0 || next.X >= m.width || next.Y < 0 || next.Y >= m.height {
				continue
			}

			// Check if the neighbor is a walkable path and hasn't been visited
			cell := m.grid[next.Y][next.X]
			if cell != Wall && !visited[next.Y][next.X] {
				visited[next.Y][next.X] = true
				parent[next] = current
				queue = append(queue, next)
			}
		}
	}

	// If a path was found, backtrack to reconstruct it.
	if pathFound {
		// Reconstruct the path by walking backwards from the end.
		// First, determine the length to pre-allocate the slice.
		pathLen := 1
		p := m.end
		for p != m.start {
			p = parent[p]
			pathLen++
		}

		// Allocate the slice and fill it from back to front, which avoids a separate reverse step.
		fullPath := make([]Point, pathLen)
		p = m.end
		for i := pathLen - 1; i >= 0; i-- {
			fullPath[i] = p
			if i > 0 { // The start point has no parent in the map.
				p = parent[p]
			}
		}
		return fullPath, true
	}
	return nil, false
}
