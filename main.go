package main

import "fmt"

func main() {
	mySudoku := newSudokuFromFile("tests/simple_3.txt")

	mySudoku.print()
	solvedSudoku, solved, err := solve(mySudoku)

	solvedSudoku.print()
	fmt.Println(solved)
	fmt.Println(err)
}
