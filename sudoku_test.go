package main

import (
	"fmt"
	"strconv"
	"testing"
)

func TestSplitString(t *testing.T) {
	strSlice := splitString("0,0,0,2,6,0,7,0,1")
	fmt.Println(strSlice)
}

func TestStringToRow(t *testing.T) {
	strSlice := splitString("0,0,0,2,6,0,7,0,1")
	myRow := _convertStringToRow(strSlice)
	myRow.print()
}

func TestNewSudokuFromFile(t *testing.T) {
	mySudoku := newSudokuFromFile("tests/simple_1.txt")

	if len(mySudoku) != 9 {
		t.Errorf("Expected 9 rows, but got %d", len(mySudoku))
	}

	for _, r := range mySudoku {
		if len(r) != 9 {
			t.Errorf("Expected 9 columns, but got %d", len(r))
		}
	}
	mySudoku.print()
}

func TestStandardMap(t *testing.T) {
	myMap := standardMap()
	fmt.Println(myMap)

	for key, val := range myMap {
		if !val {
			t.Errorf("Expected true, but got %s", strconv.FormatBool(val))
		}

		if key > 9 || key < 1 {
			t.Errorf("Only 1 to 9 allowed but got %d", key)
		}
	}
}

func TestGetEligibleList(t *testing.T) {
	myRow := make(row, 9)
	myRow[2] = 5
	myRow[7] = 8
	myRow[8] = 9
	myMap := myRow._getEligibleMap()
	fmt.Println(myMap)
}

func TestRowComplete(t *testing.T) {
	myRow := make(row, 9)
	myRow[2] = 5
	myRow[7] = 8
	myRow[8] = 9
	if myRow.complete() {
		t.Errorf("Expected Map to be incomplete, but got complete")
	}
}

func TestNonRepeatingRow(t *testing.T) {
	myRow := make(row, 9)
	for i := 1; i <= 9; i++ {
		myRow[i-1] = i
	}
	if !myRow.nonRepeating() {
		t.Errorf("Expected row to be composed of non repeating numbers")
	}

	myRow[1] = 1

	if myRow.nonRepeating() {
		t.Errorf("Expected row to be composed of one repeating number")
	}
}

func TestSudokuSolved(t *testing.T) {
	mySudoku := newSudokuFromFile("tests/simple_1.txt")
	if mySudoku.solved() {
		t.Errorf("Expected sudoku NOT to be solved, but got solved")
	}

	mySudoku = newSudokuFromFile("tests/simple_2.txt")
	if !mySudoku.solved() {
		t.Errorf("Expected sudoku to be SOLVED, but got NOT solved")
	}
}
