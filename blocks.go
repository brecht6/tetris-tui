package main

import (
	"math/rand"
)

var pieces = [7]block{
	{positions: []position{{4, 0}, {5, 0}, {6, 0}, {7, 0}}, color: "#00f0f0"}, // I
	{positions: []position{{4, 1}, {5, 1}, {4, 0}, {6, 1}}, color: "#0000f0"}, // J
	{positions: []position{{4, 1}, {5, 1}, {6, 0}, {6, 1}}, color: "#f0a000"}, // L
	{positions: []position{{4, 0}, {5, 0}, {4, 1}, {5, 1}}, color: "#f0f000"}, // O
	{positions: []position{{5, 1}, {5, 0}, {6, 0}, {4, 1}}, color: "#00f000"}, // S
	{positions: []position{{4, 1}, {5, 1}, {5, 0}, {6, 1}}, color: "#a000f0"}, // T
	{positions: []position{{4, 0}, {5, 1}, {5, 0}, {6, 1}}, color: "#f00000"}, // Z
}

func isGameOver(blocks [boardHeight][boardWidth]string, nextBlock block) bool {
	for _, pos := range nextBlock.positions {
		if !isFree(pos, blocks) {
			return true
		}
	}
	return false
}

func getRandomPiece() block {
	p := pieces[rand.Intn(len(pieces))]
	positions := make([]position, len(p.positions))
	copy(positions, p.positions)
	p.positions = positions
	return p
}

func rotateBlock(b block, blocks [boardHeight][boardWidth]string) []position {
	// if O block don't rotate
	if b.color == "#f0f000" {
		return b.positions
	}

	center := b.positions[1]

	newPositions := make([]position, len(b.positions))
	for i, pos := range b.positions {
		relX := pos.x - center.x
		relY := pos.y - center.y

		newX := center.x - relY
		newY := center.y + relX

		newPos := position{x: newX, y: newY}

		// check collision with walls
		if newPos.x < 0 || newPos.x >= boardWidth || newPos.y < 0 || newPos.y >= boardHeight {
			return b.positions
		}

		// check collision with other blocks
		if !isFree(newPos, blocks) {
			return b.positions
		}

		newPositions[i] = newPos
	}

	return newPositions
}

func isFree(pos position, blocks [boardHeight][boardWidth]string) bool {
	// check if position is empty string in blocks matrix
	if pos.y < 0 || pos.y >= len(blocks) || pos.x < 0 || pos.x >= len(blocks[pos.y]) {
		return false
	}

	return blocks[pos.y][pos.x] == ""
}

func moveBlockLeft(b block, blocks [boardHeight][boardWidth]string) []position {
	// check collision with other blocks
	for _, pos := range b.positions {
		newPos := position{x: pos.x - 1, y: pos.y}
		if !isFree(newPos, blocks) {
			return b.positions
		}
	}

	// check collision with walls
	minX := boardWidth
	for _, pos := range b.positions {
		if pos.x < minX {
			minX = pos.x
		}
	}

	if minX <= 0 {
		return b.positions
	}

	// move left
	for i := range b.positions {
		b.positions[i].x--
	}
	return b.positions
}

func moveBlockRight(b block, blocks [boardHeight][boardWidth]string) []position {
	// check collision with other blocks
	for _, pos := range b.positions {
		newPos := position{x: pos.x + 1, y: pos.y}
		if !isFree(newPos, blocks) {
			return b.positions
		}
	}

	// check collision with walls
	maxX := 0
	for _, pos := range b.positions {
		if pos.x > maxX {
			maxX = pos.x
		}
	}

	if maxX >= boardWidth-1 {
		return b.positions
	}

	// move right
	for i := range b.positions {
		b.positions[i].x++
	}
	return b.positions
}

func moveBlockDown(b block, blocks [boardHeight][boardWidth]string) []position {
	if hitBottom(b, blocks) {
		return b.positions
	}

	// move down
	for i := range b.positions {
		b.positions[i].y++
	}
	return b.positions
}

func hitBottom(b block, blocks [boardHeight][boardWidth]string) bool {
	// check collision with other blocks
	for _, pos := range b.positions {
		newPos := position{x: pos.x, y: pos.y + 1}
		if !isFree(newPos, blocks) {
			return true
		}
	}

	// check collision with floor
	maxY := 0
	for _, pos := range b.positions {
		if pos.y > maxY {
			maxY = pos.y
		}
	}

	return maxY >= boardHeight-1
}

func addBlockToMatrix(blocks [boardHeight][boardWidth]string, currentBlock block, checkLines bool) ([boardHeight][boardWidth]string, int) {
	// copy matrix
	matrix := blocks

	// add current block to matrix
	for _, pos := range currentBlock.positions {
		if pos.y >= 0 && pos.y < len(matrix) && pos.x >= 0 && pos.x < len(matrix[pos.y]) {
			matrix[pos.y][pos.x] = currentBlock.color
		}
	}

	if !checkLines {
		return matrix, 0
	}

	// unique lines of the current block
	seen := map[int]bool{}
	rows := []int{}
	for _, pos := range currentBlock.positions {
		if !seen[pos.y] {
			seen[pos.y] = true
			rows = append(rows, pos.y)
		}
	}

	// check for completed lines (only check unique lines of the current block)
	completedLines := []int{}
	for _, row := range rows {
		full := true
		for col := 0; col < boardWidth; col++ {
			if matrix[row][col] == "" {
				full = false
				break
			}
		}
		if full {
			completedLines = append(completedLines, row)
		}
	}

	// remove completed lines
	for _, row := range completedLines {
		for r := row; r > 0; r-- {
			matrix[r] = matrix[r-1]
		}
		matrix[0] = [boardWidth]string{}
	}

	return matrix, len(completedLines) * 100
}
