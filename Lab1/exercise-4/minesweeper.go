package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

func generateMinesweeperBoard(rows int, cols int, mines int) [][]string {
	board := make([][]string, rows)
	for i := range board {
		board[i] = make([]string, cols)
	}
	//invalid inputs
	if rows*cols < mines {
		fmt.Println("Invalid Board!")
		return board
	}

	//initialize board without any mines
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			board[i][j] = "."
		}
	}

	for i := 0; i < mines; i++ {
		//randomize mine location
		mineRow := rand.Intn(rows)
		mineCol := rand.Intn(cols)

		//insert mine into the board
		board[mineRow][mineCol] = "*"
	}

	//insert numbers ranging from 1 to 8
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if board[i][j] != "*" {
				adjacentMines := 0

				//iterate through adjacent tiles
				for r := i - 1; r <= i+1; r++ {
					for c := j - 1; c <= j+1; c++ {
						//make sure the location is not out of bound
						if (r >= 0 && r < rows) && (c >= 0 && c < cols) && (r != i || c != j) {
							if board[r][c] == "*" {
								adjacentMines++
							}

						}

					}
				}
				//if the number of adjacent mines > 0 then set number for that tile
				if adjacentMines > 0 {
					board[i][j] = strconv.Itoa(adjacentMines)
				}
			}
		}
	}

	return board
}

func main() {
	board := generateMinesweeperBoard(20, 25, 99)

	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			fmt.Print(board[i][j])
		}
		fmt.Println()
	}
}
