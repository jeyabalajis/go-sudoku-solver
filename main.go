package main

import "fmt"

func main() {
	mySudoku := newSudokuFromFile("tests/hard_2.txt")

	mySudoku.print()
	solvedSudoku, solved, iterations, err := solve(mySudoku)

	fmt.Println("solved:", solved)
	fmt.Println("error:", err)
	fmt.Println("total iterations: ", iterations)
	solvedSudoku.print()
}
