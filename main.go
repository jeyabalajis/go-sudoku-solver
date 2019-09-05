package main

import "fmt"

func main() {
	mySudoku := newSudokuFromFile("tests/hard_1.txt")

	mySudoku.print()
	solvedSudoku, solved, err := solve(mySudoku)

	solvedSudoku.print()
	fmt.Println(solved)
	fmt.Println(err)
}
