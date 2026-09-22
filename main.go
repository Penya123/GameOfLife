package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	width  = 80
	height = 40
	alive  = "█"
	dead   = " "
)

func main() {
	// Capturar Ctrl+C para restaurar el cursor antes de salir
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print("\033[?25h") // Mostrar cursor
		os.Exit(0)
	}()

	fmt.Print("\033[?25l") // Ocultar cursor
	fmt.Print("\033[2J")   // Limpiar pantalla
	defer fmt.Print("\033[?25h")

	grid := initGrid()

	for {
		fmt.Print("\033[H") // Mover cursor a la posición inicial (arriba a la izquierda)
		drawGrid(grid)
		grid = nextGeneration(grid)
		time.Sleep(80 * time.Millisecond) // Velocidad de la animación
	}
}

func initGrid() [][]bool {
	grid := make([][]bool, height)
	for i := range grid {
		grid[i] = make([]bool, width)
		for j := range grid[i] {
			// 20% de probabilidad de que una célula nazca viva
			grid[i][j] = rand.Float32() < 0.20
		}
	}
	return grid
}

func drawGrid(grid [][]bool) {
	var output string
	for _, row := range grid {
		for _, cell := range row {
			if cell {
				// Puedes añadir colores ANSI aquí, por ejemplo verde: "\033[32m█\033[0m"
				output += alive
			} else {
				output += dead
			}
		}
		output += "\n"
	}
	fmt.Print(output)
}

func nextGeneration(grid [][]bool) [][]bool {
	next := make([][]bool, height)
	for i := range grid {
		next[i] = make([]bool, width)
		for j := range grid[i] {
			neighbors := countNeighbors(grid, i, j)

			if grid[i][j] {
				// Reglas 1, 2 y 3: Sobrevive si tiene 2 o 3 vecinos
				next[i][j] = neighbors == 2 || neighbors == 3
			} else {
				// Regla 4: Reproducción
				next[i][j] = neighbors == 3
			}
		}
	}
	return next
}

func countNeighbors(grid [][]bool, x, y int) int {
	count := 0
	// Revisar los 8 vecinos (incluso en los bordes usando módulo para un efecto de mapa infinito o toroide)
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
