package main

import (
	"fmt"
	"os"
	"strings"
)

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
)

type Position struct {
	x, y      int
	direction Direction
}

type Grid struct {
	cells [][]rune
	rows  int
	cols  int
}

func (d Direction) turnRight() Direction {
	return (d + 1) % 4
}

func (p Position) move() Position {
	switch p.direction {
	case Up:
		return Position{p.x, p.y - 1, p.direction}
	case Right:
		return Position{p.x + 1, p.y, p.direction}
	case Down:
		return Position{p.x, p.y + 1, p.direction}
	case Left:
		return Position{p.x - 1, p.y, p.direction}
	}
	return p
}

func parseGrid(input string) *Grid {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	rows := len(lines)
	if rows == 0 {
		return nil
	}
	cols := len(lines[0])
	cells := make([][]rune, rows)

	for i, line := range lines {
		cells[i] = []rune(line)
	}

	return &Grid{cells, rows, cols}
}

func (g *Grid) isValidPosition(p Position) bool {
	return p.x >= 0 && p.x < g.cols && p.y >= 0 && p.y < g.rows
}

func findStartingPosition(grid *Grid) Position {
	for y := 0; y < grid.rows; y++ {
		for x := 0; x < grid.cols; x++ {
			if grid.cells[y][x] == '^' {
				return Position{x, y, Up}
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Starting position not found\n")
	return Position{}
}

func makeGridCopy(grid *Grid) *Grid {
	cells := make([][]rune, grid.rows)
	for y := 0; y < grid.rows; y++ {
		cells[y] = make([]rune, grid.cols)
		copy(cells[y], grid.cells[y])
	}
	return &Grid{cells, grid.rows, grid.cols}
}

func findUniquePositions(grid *Grid) int {
	pos := findStartingPosition(grid)

	visited := make(map[string]bool)

	visited[fmt.Sprintf("%d,%d", pos.x, pos.y)] = true

	for {
		nextPos := pos.move()

		if !grid.isValidPosition(nextPos) {
			return len(visited)
		}

		if grid.cells[nextPos.y][nextPos.x] == '#' {
			pos.direction = pos.direction.turnRight()
			continue
		}

		pos = nextPos
		visited[fmt.Sprintf("%d,%d", pos.x, pos.y)] = true
	}
}

func runTrial(grid *Grid) bool {
	pos := findStartingPosition(grid)

	visited := make(map[string]bool)

	visited[fmt.Sprintf("%d,%d,%d", pos.x, pos.y, pos.direction)] = true

	for {
		nextPos := pos.move()

		if !grid.isValidPosition(nextPos) {
			return false
		}

		// check if we've been here before in this direction
		if visited[fmt.Sprintf("%d,%d,%d", nextPos.x, nextPos.y, nextPos.direction)] {
			// we're in a loop
			return true
		}

		if grid.cells[nextPos.y][nextPos.x] == '#' {
			pos.direction = pos.direction.turnRight()
			continue
		}

		pos = nextPos
		visited[fmt.Sprintf("%d,%d,%d", pos.x, pos.y, pos.direction)] = true
	}
}

func runTrials(grid *Grid) int {

	numLoops := 0

	// place obstacles at each point on the grid
	for y := 0; y < grid.rows; y++ {
		fmt.Fprintf(os.Stderr, "Running trial %d/%d\n", y+1, grid.rows)
		for x := 0; x < grid.cols; x++ {

			// check if already an obstacle
			if grid.cells[y][x] == '#' || grid.cells[y][x] == '^' {
				continue
			}

			grid.cells[y][x] = '#'

			result := runTrial(grid)

			if result {
				numLoops++
			}

			// remove obstacle
			grid.cells[y][x] = '.'
		}
	}

	return numLoops
}

func main() {
	data, err := os.ReadFile("input.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v", err)
		os.Exit(1)
	}

	grid := parseGrid(string(data))
	if grid == nil {
		fmt.Fprintf(os.Stderr, "Error parsing input.txt: %v\n", err)
		os.Exit(1)
	}

	// move the guard through the grid and count the unique positions
	// uniquePositions := findUniquePositions(grid)
	// fmt.Println("Unique positions:", uniquePositions)

	// run trials and count the number of loops
	fmt.Printf("Running trials...\n")
	numLoops := runTrials(grid)
	fmt.Println("Number of loops:", numLoops)
}
