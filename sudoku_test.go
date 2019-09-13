package main

import (
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"
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

func TestGetRow(t *testing.T) {
	mySudoku := newSudokuFromFile("tests/simple_2.txt")

	testRow := row{1, 5, 2, 4, 8, 9, 3, 7, 6}
	valRow := mySudoku.getRow(0)

	for index, val := range valRow {
		if val != testRow[index] {
			t.Errorf("Correct row not retrieved")
		}
	}
}

func TestGetColumn(t *testing.T) {
	mySudoku := newSudokuFromFile("tests/simple_2.txt")

	testRow := row{5, 3, 6, 8, 9, 4, 1, 2, 7}
	valRow := mySudoku.getColumn(1)

	for index, val := range valRow {
		if val != testRow[index] {
			t.Errorf("Correct column not retrieved")
		}
	}
}

func TestGetBoundexBox(t *testing.T) {
	mySudoku := newSudokuFromFile("tests/simple_2.txt")

	testRow := row{3, 7, 6, 8, 4, 1, 2, 9, 5}
	valRow := mySudoku.getBoundedBox(0, 7)

	for index, val := range valRow {
		if val != testRow[index] {
			t.Errorf("Correct boundex box not retrieved")
		}
	}

	testRow = row{1, 2, 4, 7, 6, 3, 8, 9, 5}
	valRow = mySudoku.getBoundedBox(3, 5)

	for index, val := range valRow {
		if val != testRow[index] {
			t.Errorf("Correct bounded box not retrieved")
		}
	}
}

func TestMapEligibleNumbers(t *testing.T) {
	mySudoku := newSudokuFromFile("tests/simple_2.txt")

	myCell := mySudoku.mapEligibleNumbers(3, 5)

	for _, val := range myCell.eligibleNumbers {
		if val {
			t.Errorf("Expected no numbers to be eligible on a solved sudoku")
		}
	}

	mySudoku = newSudokuFromFile("tests/simple_1.txt")

	myCell = mySudoku.mapEligibleNumbers(0, 0)

	for index, val := range myCell.eligibleNumbers {
		if val {
			if index == 3 || index == 4 || index == 5 {
			} else {
				t.Errorf("Expected only eligible numbers to be 3, 4 or 5")
			}
		}
	}

	mySudoku = newSudokuFromFile("tests/simple_1.txt")

	myCell = mySudoku.mapEligibleNumbers(0, 1)

	for index, val := range myCell.eligibleNumbers {
		if val {
			if index == 3 {
			} else {
				t.Errorf("Expected only eligible numbers to be 3, 4 or 5")
			}
		}
	}
}

func TestUnfilledCount(t *testing.T) {
	mySudoku := newSudokuFromFile("tests/simple_2.txt")
	if mySudoku.unfilledCount() != 0 {
		t.Errorf("Expected no unfilled count on a solved sudoku")
	}

	mySudoku = newSudokuFromFile("tests/simple_1.txt")
	if mySudoku.unfilledCount() == 0 {
		t.Errorf("Expected unfilled count on an unsolved sudoku")
	}

}

func TestGetSingularEligibleNumber(t *testing.T) {
	myMap := standardMap()
	for key := range myMap {
		myMap[key] = false
		if key == 8 {
			myMap[key] = true
		}
	}

	singVal := myMap.getSingularEligibleNumber()
	if singVal != 8 {
		t.Errorf("Expected eligible number to be 8")
	}

	myMap = standardMap()
	singVal = myMap.getSingularEligibleNumber()
	if singVal != 0 {
		t.Errorf("Expected eligible number to be 0")
	}

	myMap = standardMap()
	for key := range myMap {
		myMap[key] = false
	}

	singVal = myMap.getSingularEligibleNumber()
	if singVal != -1 {
		t.Errorf("Expected no numbers to be eligible")
	}

}

func TestSolve(t *testing.T) {
	mySudoku := newSudokuFromFile("tests/hardest_3.txt")

	mySudoku.print()
	start := time.Now()
	solvedSudoku, solved, iterations, err := solve(mySudoku)
	elapsed := time.Since(start)
	fmt.Println("solved:", solved)
	fmt.Println("error:", err)
	fmt.Println("total iterations: ", iterations)
	fmt.Println("elapsed time: ", elapsed)
	solvedSudoku.print()
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
