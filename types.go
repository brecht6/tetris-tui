package main

const (
	boardWidth  = 10
	boardHeight = 20
)

type position struct {
	x int
	y int
}

type block struct {
	positions []position
	color     string
}

type model struct {
	blocks       [boardHeight][boardWidth]string
	currentBlock block
	score        int
	isGameOver   bool
	width        int
	height       int
}
