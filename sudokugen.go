package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Get a sudoku square (9x9) based on complexity sent in the input
func newSudoku(numbersToFill int) sudoku {
	var mySudoku sudoku
	var randNum int

	for i := 0; i < rowLength; i++ {
		// generate a brand new row with all zeros and append to sudoku
		mySudoku = append(mySudoku, _newRow(numbersToFill))

		// Append random numbers on specific number of places on this row
		// numbersToFill decides the complexity, i.e. how many numbers are prefilled.
		for j := 0; j < numbersToFill; j++ {
			_colIndex := genRandomNumber(colLength)
			randNum = mySudoku._genUniqueRandomNumber(i, _colIndex)
			mySudoku[i][_colIndex] = randNum
		}
	}

	return mySudoku
}

func _newRow(numbersToFill int) row {
	row := row{}

	// Append zero on all slots
	for i := 0; i < colLength; i++ {
		row = append(row, 0)
	}

	return row
}

// Get a random number integer with a fresh source every time this function is called
func genRandomNumber(maxNumber int) int {
	// Get a random number source with a fresh new seed every time this function is called
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randNum := r.Intn(maxNumber)
	return randNum
}

func (s sudoku) print() {
	for i, col := range s {
		fmt.Println(i, col)
	}
}

func (r row) print() {
	for i, col := range r {
		fmt.Println(i, col)
	}
}

// Generate a unique number that is NOT already present in the row
func (s sudoku) _genUniqueRandomNumber(rowID int, colID int) int {
	var randNum int

	fmt.Println("row: ", rowID)
	fmt.Println("col: ", colID)
	// Keep generating a random number until it satisfies sudoku criteria
	iterationNum := 0
	for {
		iterationNum++
		fmt.Println("iteration: ", iterationNum)
		randNum = genRandomNumber(colLength + 1)
		if !s.isPresent(randNum, rowID, colID) {
			break
		}
	}
	return randNum
}

/*
Validate whether the number to be filled satisfies the following criteria:
	1) It must not be present in any column of the same row
	2) It must not be present in the same column across other rows
	3) It must not be present in the bounding box (3x3) to which the (rowID, colID) belongs
*/
func (s sudoku) isPresent(numToFill int, rowID int, colID int) bool {
	// Loop through the sudoku by each row
	for _rowIndex, row := range s {
		// Loop through the sudoku row by each column
		for _colIndex, _colValue := range row {

			// If the number to fill is already present in any column of the current row
			if numToFill == _colValue && _rowIndex == rowID {
				return true
			}

			// If the number to fill is already present in the same column id of any row
			if numToFill == _colValue && _colIndex == colID {
				return true
			}

			// If the number to fill is already present in the 3x3 bounded box to which
			// (colId, rowID) belongs to
			switch {
			case colID <= 2:
				switch {
				case rowID <= 2:
					// top left bounding box
					if _isPresentBoundingBox(2, 2, _rowIndex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 5:
					// middle left bounding box
					if _isPresentBoundingBox(5, 2, _rowIndex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 8:
					// down left bounding box
					if _isPresentBoundingBox(8, 2, _rowIndex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				}
				break
			case colID <= 5:
				switch {
				case rowID <= 2:
					// top middle bounding box
					if _isPresentBoundingBox(2, 5, _rowIndex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 5:
					// middle middle bounding box
					if _isPresentBoundingBox(5, 5, _rowIndex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 8:
					// down middle bounding box
					if _isPresentBoundingBox(8, 5, _rowIndex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				}
				break
			case colID <= 8:
				switch {
				case rowID <= 2:
					// top right bounding box
					if _isPresentBoundingBox(2, 8, _rowIndex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 5:
					// middle right bounding box
					if _isPresentBoundingBox(5, 8, _rowIndex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 8:
					// down right bounding box
					if _isPresentBoundingBox(8, 8, _rowIndex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				}
				break
			}
		}
	}
	return false
}

func _isPresentBoundingBox(rowBoundary int, colBoundary int, rowIndex int, colIndex int, numToFill int, colValue int) bool {

	var rowLowerBoundary int
	var rowUpperBoundary int
	var colLowerBoundary int
	var colUpperBoundary int

	switch {
	case rowBoundary <= 2 && colBoundary <= 2:
		rowLowerBoundary = 0
		rowUpperBoundary = 2
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowBoundary <= 5 && colBoundary <= 2:
		rowLowerBoundary = 3
		rowUpperBoundary = 5
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowBoundary <= 8 && colBoundary <= 2:
		rowLowerBoundary = 6
		rowUpperBoundary = 8
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowBoundary <= 2 && colBoundary <= 5:
		rowLowerBoundary = 0
		rowUpperBoundary = 2
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowBoundary <= 5 && colBoundary <= 5:
		rowLowerBoundary = 3
		rowUpperBoundary = 5
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowBoundary <= 8 && colBoundary <= 5:
		rowLowerBoundary = 6
		rowUpperBoundary = 8
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowBoundary <= 2 && colBoundary <= 8:
		rowLowerBoundary = 0
		rowUpperBoundary = 2
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	case rowBoundary <= 5 && colBoundary <= 8:
		rowLowerBoundary = 3
		rowUpperBoundary = 5
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	case rowBoundary <= 8 && colBoundary <= 8:
		rowLowerBoundary = 6
		rowUpperBoundary = 8
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	}

	if rowIndex >= rowLowerBoundary && rowIndex <= rowUpperBoundary && colIndex >= colLowerBoundary && colIndex <= colUpperBoundary && colValue == numToFill {
		return true
	}

	return false
}
