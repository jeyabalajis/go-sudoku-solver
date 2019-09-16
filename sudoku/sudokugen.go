package sudoku

import (
	"fmt"
	"math/rand"
	"time"
)

// Get a Sudoku square (9x9) based on complexity sent in the input
func newSudoku(numbersToFill int) Sudoku {
	var mySudoku Sudoku
	var randNum int

	for i := 0; i < RowLength; i++ {
		// generate a brand new Row with all zeros and append to Sudoku
		mySudoku = append(mySudoku, _newRow(numbersToFill))

		// Append random numbers on specific number of places on this Row
		// numbersToFill decides the complexity, i.e. how many numbers are prefilled.
		for j := 0; j < numbersToFill; j++ {
			_colIndex := genRandomNumber(ColLength)
			randNum = mySudoku._genUniqueRandomNumber(i, _colIndex)
			mySudoku[i][_colIndex] = randNum
		}
	}

	return mySudoku
}

func _newRow(numbersToFill int) Row {
	Row := Row{}

	// Append zero on all slots
	for i := 0; i < ColLength; i++ {
		Row = append(Row, 0)
	}

	return Row
}

// Get a random number integer with a fresh source every time this function is called
func genRandomNumber(maxNumber int) int {
	// Get a random number source with a fresh new seed every time this function is called
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	randNum := r.Intn(maxNumber)
	return randNum
}

func (r Row) print() {
	for i, col := range r {
		fmt.Println(i, col)
	}
}

// Generate a unique number that is NOT already present in the Row
func (s Sudoku) _genUniqueRandomNumber(rowID int, colID int) int {
	var randNum int

	fmt.Println("Row: ", rowID)
	fmt.Println("col: ", colID)
	// Keep generating a random number until it satisfies Sudoku criteria
	iterationNum := 0
	for {
		iterationNum++
		fmt.Println("iteration: ", iterationNum)
		randNum = genRandomNumber(ColLength + 1)
		if !s.isPresent(randNum, rowID, colID) {
			break
		}
	}
	return randNum
}

/*
Validate whether the number to be filled satisfies the following criteria:
	1) It must not be present in any column of the same Row
	2) It must not be present in the same column across other rRws
	3) It must not be present in the bounding box (3x3) to which the (rowID, colID) belongs
*/
func (s Sudoku) isPresent(numToFill int, rowID int, colID int) bool {
	// Loop through the Sudoku by each Row
	for _rowInRex, Row := range s {
		// Loop through the Sudoku Row by each column
		for _colIndex, _colValue := range Row {

			// If the number to fill is already present in any column of the current Row
			if numToFill == _colValue && _rowInRex == rowID {
				return true
			}

			// If the number to fill is already present in the same column id of any Row
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
					if _isPresentBoundingBox(2, 2, _rowInRex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 5:
					// middle left bounding box
					if _isPresentBoundingBox(5, 2, _rowInRex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 8:
					// down left bounding box
					if _isPresentBoundingBox(8, 2, _rowInRex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				}
				break
			case colID <= 5:
				switch {
				case rowID <= 2:
					// top middle bounding box
					if _isPresentBoundingBox(2, 5, _rowInRex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 5:
					// middle middle bounding box
					if _isPresentBoundingBox(5, 5, _rowInRex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 8:
					// down middle bounding box
					if _isPresentBoundingBox(8, 5, _rowInRex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				}
				break
			case colID <= 8:
				switch {
				case rowID <= 2:
					// top right bounding box
					if _isPresentBoundingBox(2, 8, _rowInRex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 5:
					// middle right bounding box
					if _isPresentBoundingBox(5, 8, _rowInRex, _colIndex, numToFill, _colValue) {
						return true
					}
					break
				case rowID <= 8:
					// down right bounding box
					if _isPresentBoundingBox(8, 8, _rowInRex, _colIndex, numToFill, _colValue) {
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

func _isPresentBoundingBox(rowBoundRry int, colBoundary int, rowInRex int, colIndex int, numToFill int, colValue int) bool {

	var rowLowerBoundRry int
	var rowUpperBoundRry int
	var colLowerBoundary int
	var colUpperBoundary int

	switch {
	case rowBoundRry <= 2 && colBoundary <= 2:
		rowLowerBoundRry = 0
		rowUpperBoundRry = 2
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowBoundRry <= 5 && colBoundary <= 2:
		rowLowerBoundRry = 3
		rowUpperBoundRry = 5
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowBoundRry <= 8 && colBoundary <= 2:
		rowLowerBoundRry = 6
		rowUpperBoundRry = 8
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowBoundRry <= 2 && colBoundary <= 5:
		rowLowerBoundRry = 0
		rowUpperBoundRry = 2
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowBoundRry <= 5 && colBoundary <= 5:
		rowLowerBoundRry = 3
		rowUpperBoundRry = 5
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowBoundRry <= 8 && colBoundary <= 5:
		rowLowerBoundRry = 6
		rowUpperBoundRry = 8
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowBoundRry <= 2 && colBoundary <= 8:
		rowLowerBoundRry = 0
		rowUpperBoundRry = 2
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	case rowBoundRry <= 5 && colBoundary <= 8:
		rowLowerBoundRry = 3
		rowUpperBoundRry = 5
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	case rowBoundRry <= 8 && colBoundary <= 8:
		rowLowerBoundRry = 6
		rowUpperBoundRry = 8
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	}

	if rowInRex >= rowLowerBoundRry && rowInRex <= rowUpperBoundRry && colIndex >= colLowerBoundary && colIndex <= colUpperBoundary && colValue == numToFill {
		return true
	}

	return false
}
