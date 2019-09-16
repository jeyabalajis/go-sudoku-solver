package main

import (
	"fmt"
	"gosudokusolver/sudoku"
	"time"
)

func main() {
	mySudoku := sudoku.NewSudokuFromFile("tests/hardest_3.txt")

	mySudoku.Print()
	start := time.Now()
	solvedSudoku, solved, iterations, err := sudoku.Solve(mySudoku)
	elapsed := time.Since(start)
	fmt.Println("solved:", solved)
	fmt.Println("error:", err)
	fmt.Println("total iterations: ", iterations)
	fmt.Println("elapsed time: ", elapsed)
	solvedSudoku.Print()
}
