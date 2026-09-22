package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/term"
)

const (
	alive = "█"
	dead  = " "
)

var colors = map[string]string{
	"green":   "\033[32m",
	"red":     "\033[31m",
	"blue":    "\033[34m",
	"magenta": "\033[35m",
	"cyan":    "\033[36m",
	"white":   "\033[37m",
}

const resetColor = "\033[0m"

func main() {
	// Flags actualizadas
	colorFlag := flag.String("color", "white", "Define el color de las células (green, red, blue, magenta, cyan, white)")
	delayFlag := flag.Duration("delay", 80*time.Millisecond, "Tiempo de espera entre cada generación (ej. 100ms, 0.5s)")
	flag.Parse()

	selectedColor, exists := colors[*colorFlag]
	if !exists {
		selectedColor = colors["white"] // Fallback al nuevo default
	}

	fd := int(os.Stdout.Fd())
	width, height, err := term.GetSize(fd)
	if err != nil {
		width, height = 80, 24
	}
	height--

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print(resetColor)
		fmt.Print("\033[?25h")
		os.Exit(0)
	}()

	fmt.Print("\033[?25l")
	fmt.Print("\033[2J")
	defer fmt.Print("\033[?25h")

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	grid := initGrid(width, height, rng)

	for {
		fmt.Print("\033[H")
		drawGrid(grid, selectedColor)
		grid = nextGeneration(grid, width, height)
		// Usamos el valor de la flag de duración directamente
		time.Sleep(*delayFlag)
	}
}

func initGrid(width, height int, rng *rand.Rand) [][]bool {
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
		for j := range grid[i] {
			grid[i][j] = rng.Float32() < 0.20
		}
	}
	return grid
}

func drawGrid(grid [][]bool, colorCode string) {
	output := colorCode
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				output += alive
			} else {
				output += dead
			}
		}
		output += "\n"
	}
	fmt.Print(output)
}

func nextGeneration(grid [][]bool, width, height int) [][]bool {
	next := make([][]bool, height)
	for i := range grid {
		next[i] = make([]bool, width)
		for j := range grid[i] {
			neighbors := countNeighbors(grid, i, j, width, height)

			if grid[i][j] {
				next[i][j] = neighbors == 2 || neighbors == 3
			} else {
				next[i][j] = neighbors == 3
			}
		}
	}
	return next
}

func countNeighbors(grid [][]bool, x, y, width, height int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			nx := (x + i + height) % height
			ny := (y + j + width) % width
			if grid[nx][ny] {
				count++
			}
		}
	}
	return count
}
