package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/jeyabalajis/go-sudoku-solver/sudoku"
)

func TestNewSudokuFromFile(t *testing.T) {
	mySudoku := sudoku.NewSudokuFromFile("tests/medium_1.txt")

	if len(mySudoku) != 9 {
		t.Errorf("Expected 9 rows, but got %d", len(mySudoku))
	}

	for _, r := range mySudoku {
		if len(r) != 9 {
			t.Errorf("Expected 9 columns, but got %d", len(r))
		}
	}
	mySudoku.Print()
}

func TestSudokuSolved(t *testing.T) {
	mySudoku := sudoku.NewSudokuFromFile("tests/simple_1.txt")
	if mySudoku.Solved() {
		t.Errorf("Expected sudoku NOT to be solved, but got solved")
	}

	mySudoku = sudoku.NewSudokuFromFile("tests/simple_2.txt")
	if !mySudoku.Solved() {
		t.Errorf("Expected sudoku to be SOLVED, but got NOT solved")
	}
}

func TestSolve(t *testing.T) {
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
	if !solved {
		t.Errorf("Expected sudoku to be solved")
	}
}

func TestErrValidate(t *testing.T) {
	myErr := errors.New("incorrect sudoku")

	fmt.Println(myErr)
	if myErr.Error() == "incorrect sudoku" {

	} else {
		t.Errorf("Expected the error to match")
	}
}
