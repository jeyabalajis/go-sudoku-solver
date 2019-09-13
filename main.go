package main

import (
	"fmt"
	"time"
)

func main() {
	mySudoku := newSudokuFromFile("tests/hard_1.txt")

	mySudoku.print()
	start := time.Now()
	solvedSudoku, solved, iterations, err := solve(mySudoku)
	elapsed := time.Since(start)
	fmt.Println("solved:", solved)
	fmt.Println("error:", err)
	fmt.Println("total iterations: ", iterations)
	fmt.Println("elapsed time: ", elapsed)
	solvedSudoku.print()
}
