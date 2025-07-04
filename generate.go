package maze

import (
	"fmt"
	"math/rand"
)

// Generate creates the maze paths using an iterative randomized depth-first search.
// It takes a seed for reproducibility, an optional start point, and a bias
// that controls the straightness of corridors.
func (m *Maze) Generate(seed int64, start *Point, door *Point, doorSide string, bias float64) error {
	r := rand.New(rand.NewSource(seed))
	var generationStart Point

	// 1. Choose a starting point.
	if start != nil {
		// Use the provided start point after validation.
		if start.X <= 0 || start.X >= m.width-1 || start.Y <= 0 || start.Y >= m.height-1 || start.X%2 == 0 || start.Y%2 == 0 {
			return fmt.Errorf("invalid start point: %+v. must be within maze bounds and have odd coordinates", *start)
		}
		if m.IsInsideDen(*start) {
			return fmt.Errorf("invalid start point: %+v. cannot start generation inside the den", *start)
		}
		generationStart = *start
	} else {
		// Choose a random starting point (must be on a path cell, so odd coordinates).
		// Loop until we find a point that is not inside the den.
		for {
			startX := r.Intn((m.width-1)/2)*2 + 1
			startY := r.Intn((m.height-1)/2)*2 + 1
			p := Point{X: startX, Y: startY}
			if !m.IsInsideDen(p) {
				generationStart = p
				break
			}
		}
	}

	// 2. Run the generation algorithm.
	m.runDFS(r, generationStart, bias)

	// 3. If a den exists, create a single door to connect it to the maze.
	if err := m.connectDen(r, door, doorSide); err != nil {
		return err
	}

	// 4. Set the Start and End points for the maze.
	m.placeStartAndEnd(generationStart, start != nil)

	return nil
}

// runDFS executes the iterative depth-first search algorithm to carve the maze paths.
func (m *Maze) runDFS(r *rand.Rand, start Point, bias float64) {
	var stack []Point

	current := start
	m.grid[current.Y][current.X] = Path
	stack = append(stack, current)

	for len(stack) > 0 {
		current = stack[len(stack)-1]

		neighbors := m.findValidNeighbors(current)

		if len(neighbors) > 0 {
			next := chooseBiasedNeighbor(neighbors, stack, bias, r)

			// Carve a path between the current cell and the neighbor
			wallToRemove := Point{
				X: current.X + (next.X-current.X)/2,
				Y: current.Y + (next.Y-current.Y)/2,
			}
			m.grid[wallToRemove.Y][wallToRemove.X] = Path
			m.grid[next.Y][next.X] = Path

			stack = append(stack, next)
		} else {
			// If no unvisited neighbors, backtrack by popping from the stack
			stack = stack[:len(stack)-1]
		}
	}
}

// findValidNeighbors finds all unvisited neighbors of a point that can be carved into.
func (m *Maze) findValidNeighbors(p Point) []Point {
	var neighbors []Point
	directions := []Point{{X: 0, Y: -2}, {X: 0, Y: 2}, {X: -2, Y: 0}, {X: 2, Y: 0}}

	for _, dir := range directions {
		next := Point{X: p.X + dir.X, Y: p.Y + dir.Y}

		// Check if the neighbor is a valid, unvisited cell that doesn't breach the den.
		if next.X > 0 && next.X < m.width-1 && next.Y > 0 && next.Y < m.height-1 && m.grid[next.Y][next.X] == Wall {
			wallBetween := Point{X: p.X + dir.X/2, Y: p.Y + dir.Y/2}
			if m.IsInsideDen(wallBetween) || m.IsAdjacentToDen(next) {
				continue
			}
			neighbors = append(neighbors, next)
		}
	}
	return neighbors
}

// chooseBiasedNeighbor selects a neighbor from a list, applying a bias to continue in a straight line.
func chooseBiasedNeighbor(neighbors []Point, stack []Point, bias float64, r *rand.Rand) Point {
	// Determine the last direction of travel.
	var lastDirection Point
	if len(stack) > 1 {
		previous := stack[len(stack)-2]
		current := stack[len(stack)-1]
		lastDirection = Point{X: current.X - previous.X, Y: current.Y - previous.Y}
	}

	// Check if moving straight is a valid option.
	var straightOption *Point
	for i := range neighbors {
		n := neighbors[i]
		current := stack[len(stack)-1]
		if (n.X-current.X) == lastDirection.X && (n.Y-current.Y) == lastDirection.Y {
			straightOption = &neighbors[i]
			break
		}
	}

	if straightOption != nil && r.Float64() < bias {
		return *straightOption
	}

	// Otherwise, pick a random neighbor from the available options.
	return neighbors[r.Intn(len(neighbors))]
}

// placeStartAndEnd determines and sets the Start and End points on the maze grid.
func (m *Maze) placeStartAndEnd(generationStart Point, useProvidedStart bool) {
	if useProvidedStart {
		m.start = generationStart
	} else {
		// If no start point was provided, find the longest path in the maze.
		// The start of the longest path is the point farthest from the generation start.
		m.start, _ = m.findFarthestPoint(generationStart)
	}

	// The end of the longest path is the point farthest from our new start point.
	m.end, _ = m.findFarthestPoint(m.start)

	// Place Start and End markers on the grid.
	m.grid[m.start.Y][m.start.X] = Start
	m.grid[m.end.Y][m.end.X] = End
}

// connectDen finds all possible walls that can be turned into a door
// between the maze and the den, and randomly picks one to open.
func (m *Maze) connectDen(r *rand.Rand, userDoor *Point, doorSide string) error {
	if m.denWidth <= 0 || m.denHeight <= 0 {
		return nil // No den to connect.
	}

	// Handle specified door side (e.g., "top", "bottom").
	if doorSide != "" {
		return m.connectDenAtSide(doorSide)
	}

	// If a specific door location is provided, validate and use it.
	if userDoor != nil {
		return m.connectDenAtPoint(*userDoor)
	}

	return m.connectRandomDenDoor(r)
}

// connectDenAtSide connects the den to the maze at the center of a specified wall.
func (m *Maze) connectDenAtSide(doorSide string) error {
	var door, neighbor Point
	switch doorSide {
	case "top":
		door = Point{X: m.denStartX + m.denWidth/2, Y: m.denStartY - 1}
		neighbor = Point{X: door.X, Y: door.Y - 1}
	case "bottom":
		door = Point{X: m.denStartX + m.denWidth/2, Y: m.denStartY + m.denHeight}
		neighbor = Point{X: door.X, Y: door.Y + 1}
	case "left":
		door = Point{X: m.denStartX - 1, Y: m.denStartY + m.denHeight/2}
		neighbor = Point{X: door.X - 1, Y: door.Y}
	case "right":
		door = Point{X: m.denStartX + m.denWidth, Y: m.denStartY + m.denHeight/2}
		neighbor = Point{X: door.X + 1, Y: door.Y}
	default:
		return fmt.Errorf("invalid door side: %s. use 'top', 'bottom', 'left', or 'right'", doorSide)
	}

	// Check if the door and neighbor are within bounds.
	if door.X <= 0 || door.X >= m.width-1 || door.Y <= 0 || door.Y >= m.height-1 ||
		neighbor.X <= 0 || neighbor.X >= m.width-1 || neighbor.Y <= 0 || neighbor.Y >= m.height-1 {
		return fmt.Errorf("cannot place door on side '%s': too close to maze edge", doorSide)
	}

	// If the neighbor is already a path, we just need to open the door.
	if m.grid[neighbor.Y][neighbor.X] == Path {
		m.grid[door.Y][door.X] = Path
		m.door = door
		return nil
	}

	// Otherwise, we need to carve a path from the neighbor to the nearest maze path.
	if err := m.carvePathToNearest(neighbor); err != nil {
		return fmt.Errorf("failed to connect door on side '%s': %w", doorSide, err)
	}

	// Finally, open the door itself.
	m.grid[door.Y][door.X] = Path
	m.door = door
	return nil
}

// connectDenAtPoint connects the den to the maze at a user-specified point.
func (m *Maze) connectDenAtPoint(userDoor Point) error {
	// 1. Must be a wall within the maze's inner boundaries.
	if userDoor.X <= 0 || userDoor.X >= m.width-1 || userDoor.Y <= 0 || userDoor.Y >= m.height-1 || m.grid[userDoor.Y][userDoor.X] != Wall {
		return fmt.Errorf("invalid door location at %+v: not a valid wall position", userDoor)
	}

	// 2. Must be on the boundary of the den, connecting an inner path to an outer path.
	// Check for a horizontal connection: Path-Wall-Path
	p1_h := Point{X: userDoor.X - 1, Y: userDoor.Y}
	p2_h := Point{X: userDoor.X + 1, Y: userDoor.Y}
	if m.grid[p1_h.Y][p1_h.X] == Path && m.grid[p2_h.Y][p2_h.X] == Path && m.IsInsideDen(p1_h) != m.IsInsideDen(p2_h) {
		m.grid[userDoor.Y][userDoor.X] = Path
		m.door = userDoor
		return nil
	}

	// Check for a vertical connection: Path-Wall-Path
	p1_v := Point{X: userDoor.X, Y: userDoor.Y - 1}
	p2_v := Point{X: userDoor.X, Y: userDoor.Y + 1}
	if m.grid[p1_v.Y][p1_v.X] == Path && m.grid[p2_v.Y][p2_v.X] == Path && m.IsInsideDen(p1_v) != m.IsInsideDen(p2_v) {
		m.grid[userDoor.Y][userDoor.X] = Path
		m.door = userDoor
		return nil
	}

	return fmt.Errorf("invalid door location at %+v: does not connect the den to a maze path", userDoor)
}

// connectRandomDenDoor finds all possible walls that can be turned into a door
// and randomly picks one to open.
func (m *Maze) connectRandomDenDoor(r *rand.Rand) error {
	var potentialDoors []Point

	// Iterate through the grid to find walls that separate the den from the maze path.
	for y := 1; y < m.height-1; y++ {
		for x := 1; x < m.width-1; x++ {
			// We are looking for a Wall cell to serve as a door.
			if m.grid[y][x] != Wall {
				continue
			}

			// Check for a horizontal separation: Path-Wall-Path
			p1_h := Point{X: x - 1, Y: y}
			p2_h := Point{X: x + 1, Y: y}
			if m.grid[p1_h.Y][p1_h.X] == Path && m.grid[p2_h.Y][p2_h.X] == Path {
				// If one side is in the den and the other isn't, it's a valid door.
				if m.IsInsideDen(p1_h) != m.IsInsideDen(p2_h) {
					potentialDoors = append(potentialDoors, Point{X: x, Y: y})
					continue // Found a candidate, move to the next wall cell.
				}
			}

			// Check for a vertical separation: Path-Wall-Path
			p1_v := Point{X: x, Y: y - 1}
			p2_v := Point{X: x, Y: y + 1}
			if m.grid[p1_v.Y][p1_v.X] == Path && m.grid[p2_v.Y][p2_v.X] == Path {
				if m.IsInsideDen(p1_v) != m.IsInsideDen(p2_v) {
					potentialDoors = append(potentialDoors, Point{X: x, Y: y})
				}
			}
		}
	}

	if len(potentialDoors) > 0 {
		// Pick a random door from all possibilities and open it.
		door := potentialDoors[r.Intn(len(potentialDoors))]
		m.grid[door.Y][door.X] = Path
		m.door = door
	}

	return nil // It's not an error if no potential doors are found.
}

// carvePathToNearest finds the closest maze path from a starting point (through walls)
// and carves a corridor to connect them.
func (m *Maze) carvePathToNearest(start Point) error {
	// This function uses BFS to find the nearest Path cell, exploring only through Wall cells.
	queue := []Point{start}
	visited := make(map[Point]bool)
	visited[start] = true
	parent := make(map[Point]Point)

	var targetPath Point
	pathFound := false

	head := 0
head_loop:
	for head < len(queue) {
		current := queue[head]
		head++

		// Explore neighbors
		for _, dir := range []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
			next := Point{X: current.X + dir.X, Y: current.Y + dir.Y}

			// Check bounds
			if next.X <= 0 || next.X >= m.width-1 || next.Y <= 0 || next.Y >= m.height-1 {
				continue
			}

			// If we found an existing maze path, we're done searching.
			if m.grid[next.Y][next.X] == Path {
				parent[next] = current
				targetPath = next
				pathFound = true
				break head_loop // Exit the outer loop
			}

			// Otherwise, if it's a wall we haven't visited, add it to the queue.
			if m.grid[next.Y][next.X] == Wall {
				if _, ok := visited[next]; !ok {
					visited[next] = true
					parent[next] = current
					queue = append(queue, next)
				}
			}
		}
	}

	if !pathFound {
		return fmt.Errorf("no path found to connect the door to the maze")
	}

	// Backtrack from the target path to the start point, carving a path.
	p := targetPath
	for p != start {
		p = parent[p]
		m.grid[p.Y][p.X] = Path
	}

	return nil
}

// findFarthestPoint performs a BFS from a given start point to find the
// cell that is the farthest away along the maze paths.
// It returns the farthest point and its distance.
func (m *Maze) findFarthestPoint(start Point) (farthestPoint Point, maxDistance int) {
	queue := []Point{start}
	// distances map also serves as the visited set
	distances := make(map[Point]int)
	distances[start] = 0

	farthestPoint = start
	maxDistance = 0

	head := 0
	for head < len(queue) {
		current := queue[head]
		head++

		// Explore neighbors
		for _, dir := range []Point{{0, -1}, {0, 1}, {-1, 0}, {1, 0}} {
			next := Point{X: current.X + dir.X, Y: current.Y + dir.Y}

			// Check if the neighbor is a valid path and hasn't been visited.
			if next.X > 0 && next.X < m.width-1 && next.Y > 0 && next.Y < m.height-1 && m.grid[next.Y][next.X] != Wall {
				if _, visited := distances[next]; !visited {
					dist := distances[current] + 1
					distances[next] = dist
					queue = append(queue, next)

					// Update the farthest point only if it's not inside the den.
					// Also ensure it's not on the den's wall (i.e., the door).
					if dist > maxDistance && !m.IsInsideDen(next) && !m.IsAdjacentToDen(next) {
						maxDistance = dist
						farthestPoint = next
					}
				}
			}
		}
	}
	return farthestPoint, maxDistance
}
