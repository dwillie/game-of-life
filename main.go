package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func randomizeGrid(grid [][]int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))

	for x := range grid {
		for y := range grid[x] {
			val := r.Intn(2)
			//fmt.Printf("Setting %d, %d to %d\n", x, y, val)
			grid[x][y] = val
		}
	}
}

func renderGrid(grid [][]int) {
	var s string

	for x := range grid {
		fmt.Printf("|")
		for y := range grid[x] {
			if grid[x][y] > 0 {
				s = "X"
			} else {
				s = " "
			}
			fmt.Printf("%s", s)
		}
		fmt.Printf("|\n")
	}

	for _, _ = range grid[0] {
		fmt.Printf("-")
	}
	fmt.Printf("-\n")
}

func liveNeighborCount(x int, y int, grid [][]int) int {
	dxdy := [8][2]int{
		{-1, -1}, {0, -1}, {1, -1},
		{-1, 0}, {1, 0},
		{-1, 1}, {0, 1}, {1, 1},
	}

	count := 0
	maxX := len(grid) - 1
	maxY := len(grid[x]) - 1
	for _, d := range dxdy {
		dx := x + d[0]
		dy := y + d[1]
		if dx < 0 {
			dx = maxX
		} else if dx > maxX {
			dx = 0
		}
		if dy < 0 {
			dy = maxY
		} else if dy > maxY {
			dy = 0
		}

		count += grid[dx][dy]
	}

	return count
}

func copyGrid(grid [][]int) [][]int {
	newGrid := make([][]int, len(grid))

	for x := range grid {
		newGrid[x] = make([]int, len(grid[x]))
		for y := range grid[x] {
			newGrid[x][y] = grid[x][y]
		}
	}

	return newGrid
}

func updateGridSlice(oldGrid [][]int, newGrid [][]int, xLower int, xUpper int, barrier *sync.WaitGroup) {
	defer barrier.Done()

	for x := xLower; x < xUpper; x++ {
		for y := range oldGrid[x] {
			alive := oldGrid[x][y] == 1
			count := liveNeighborCount(x, y, oldGrid)

			if (alive && count == 2) || count == 3 {
				newGrid[x][y] = 1
			} else {
				newGrid[x][y] = 0
			}
		}
	}
}

func updateGrid(grid [][]int) {
	oldGrid := copyGrid(grid)

	gridWidth := len(grid)
	routines := 8
	routineChunkSize := gridWidth / routines
	barrier := sync.WaitGroup{}

	barrier.Add(routines)
	for i := 0; i < routines; i++ {
		go updateGridSlice(oldGrid, grid, i*routineChunkSize, (i*routineChunkSize)+routineChunkSize, &barrier)
	}

	barrier.Wait()
}

func main() {
	grid := make([][]int, 48)
	for x := range grid {
		grid[x] = make([]int, 202)
	}
	randomizeGrid(grid)

	for {
		renderGrid(grid)
		time.Sleep(1000 * 1000 * 25)
		updateGrid(grid)
	}
}
