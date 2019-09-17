package sudoku

import (
	"fmt"
)

// Row is a slice of integers. Used for representing a row, column or a bounded box of a sudoku of length 9
type Row []int

// Sudoku is a slice of Rows. Used for representing a sudoku puzzle, solved or otherwise
type Sudoku []Row

// EligibleNumbers - A Hashmap to keep track of which numbers are eligible to be filled in a column
type EligibleNumbers map[int]bool

// RowLength is a constant that represents the length of a sudoku row
const RowLength int = 9

// ColLength is a constant that represents the length of a sudoku column
const ColLength int = 9

// Cell is a structure that contains a sudoku position (RowID, ColID) and the eligible numbers that can be filled in it
type Cell struct {
	RowID           int
	ColID           int
	EligibleNumbers EligibleNumbers
}

// Channel is a structure that is used to communicate the results of a solve run that is run concurrently
type Channel struct {
	Intermediate Sudoku
	Solved       bool
	Err          error
}

func _getBoundedBoxIndex(rowID int, colID int) int {
	bbIndex := 0
	switch {
	case rowID < 3:
		switch {
		case colID < 3:
			bbIndex = 0
			break
		case colID < 6:
			bbIndex = 1
			break
		case colID < 9:
			bbIndex = 2
			break
		}
		break
	case rowID < 6:
		switch {
		case colID < 3:
			bbIndex = 3
			break
		case colID < 6:
			bbIndex = 4
			break
		case colID < 9:
			bbIndex = 5
			break
		}
		break
	case rowID < 9:
		switch {
		case colID < 3:
			bbIndex = 6
			break
		case colID < 6:
			bbIndex = 7
			break
		case colID < 9:
			bbIndex = 8
			break
		}
		break
	}
	return bbIndex
}

func _getBoundedBoxBoundaries(rowID int, colID int) (int, int, int, int) {
	var RowLowerBoundary int
	var RowUpperBoundary int
	var colLowerBoundary int
	var colUpperBoundary int

	switch {
	case rowID <= 2 && colID <= 2:
		RowLowerBoundary = 0
		RowUpperBoundary = 2
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowID <= 5 && colID <= 2:
		RowLowerBoundary = 3
		RowUpperBoundary = 5
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowID <= 8 && colID <= 2:
		RowLowerBoundary = 6
		RowUpperBoundary = 8
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowID <= 2 && colID <= 5:
		RowLowerBoundary = 0
		RowUpperBoundary = 2
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowID <= 5 && colID <= 5:
		RowLowerBoundary = 3
		RowUpperBoundary = 5
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowID <= 8 && colID <= 5:
		RowLowerBoundary = 6
		RowUpperBoundary = 8
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowID <= 2 && colID <= 8:
		RowLowerBoundary = 0
		RowUpperBoundary = 2
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	case rowID <= 5 && colID <= 8:
		RowLowerBoundary = 3
		RowUpperBoundary = 5
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	case rowID <= 8 && colID <= 8:
		RowLowerBoundary = 6
		RowUpperBoundary = 8
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	}

	return RowLowerBoundary, RowUpperBoundary, colLowerBoundary, colUpperBoundary
}

// Solved validates whether a sudoku is fully solved or not
func (s Sudoku) Solved() bool {
	/*A Sudoku is considered solved when
	- there are no empty Cells (i.e. Cells with number zero)
	- all Rows, columns and bounded box contain numbers from 1 to 9 (i.e. complete)
	- there are no repeating numbers in Rows, columns or bounded box (i.e. nonRepeating)
	*/

	myColumns := make(map[int]Row)
	myBoundedBoxes := make(map[int]Row)

	// traverse the Sudoku once to collect Rows, columns and bounded boxes
	for rowID, Row := range s {

		if !(Row.complete() && Row.nonRepeating()) {
			return false
		}

		for colID, col := range Row {

			// collect column values belonging to the same column id in a separate Row
			myColumns[colID] = append(myColumns[colID], col)

			// collect column values belonging to the same bounded box into a separate Row
			bbID := _getBoundedBoxIndex(rowID, colID)
			myBoundedBoxes[bbID] = append(myBoundedBoxes[bbID], col)
		}
	}

	if len(myColumns) > 0 {
		for _, Row := range myColumns {
			if !(Row.complete() && Row.nonRepeating()) {
				return false
			}
		}
	}

	if len(myBoundedBoxes) > 0 {
		for _, Row := range myBoundedBoxes {
			if !(Row.complete() && Row.nonRepeating()) {
				return false
			}
		}
	}

	return true
}

// Copy is a deep copy solution for Sudoku structure which is array of array of int
func (s Sudoku) Copy() Sudoku {
	mySudoku := make(Sudoku, 0)
	done := make(chan struct{})

	go func() {
		for _, _Row := range s {
			myRow := make(Row, 0)
			for _, _col := range _Row {
				myRow = append(myRow, _col)
			}
			mySudoku = append(mySudoku, myRow)
		}
		done <- struct{}{}
	}()
	<-done
	return mySudoku
}

func (s Sudoku) getRow(rowID int) Row {
	return s[rowID]
}

func (s Sudoku) getColumn(colID int) Row {
	var myColumn Row

	for _, Row := range s {
		for colIndex, col := range Row {
			if colID == colIndex {
				myColumn = append(myColumn, col)
			}
		}
	}
	return myColumn
}

// Print prints the rows of a sudoku
func (s Sudoku) Print() {
	for i, col := range s {
		fmt.Println(i, col)
	}
}

func (s Sudoku) getBoundedBox(rowID int, colID int) Row {
	var myBB Row
	RowMin, RowMax, colMin, colMax := _getBoundedBoxBoundaries(rowID, colID)

	for RowIndex, Row := range s {
		for colIndex, col := range Row {
			if (RowIndex >= RowMin && RowIndex <= RowMax) && (colIndex >= colMin && colIndex <= colMax) {
				myBB = append(myBB, col)
			}
		}
	}
	return myBB
}

// MapEligibleNumbers maps eligible numbers that can be filled for a particular position in a sudoku
func (s Sudoku) MapEligibleNumbers(rowID int, colID int) Cell {
	eligibleNumsMap := s._getEligibleMap(rowID, colID)
	return Cell{RowID: rowID, ColID: colID, EligibleNumbers: eligibleNumsMap}
}

// ReduceAndFillEligibleNumber updates a specific Cell if the Cell contains only one eligible number
func (s Sudoku) ReduceAndFillEligibleNumber(ec Cell) int {
	myMap := ec.EligibleNumbers
	rowID := ec.RowID
	colID := ec.ColID

	// ec.eligibleNumbers.Print()
	eligNum := myMap.getSingularEligibleNumber()

	if eligNum >= 1 && eligNum <= 9 {
		s.Fill(rowID, colID, eligNum)
	}

	return eligNum
}

// Fill updates a particular cell, represented by rowID and colID with the value passed
func (s Sudoku) Fill(rowID int, colID int, val int) {
	done := make(chan struct{})

	go func(s Sudoku) {
		s[rowID][colID] = val
		done <- struct{}{}
	}(s)

	<-done
}

// UnfilledCount returns the total cells that are still yet to be filled in a sudoku
func (s Sudoku) UnfilledCount() (unfilled int) {
	done := make(chan struct{})

	go func() {
		for _, Row := range s {
			for _, col := range Row {
				if col == 0 {
					unfilled++
				}
			}
		}
		done <- struct{}{}
	}()

	<-done
	return unfilled
}

// getEligibleMap gets a map of eligible numbers for a particular position
// in the Sudoku puzzle
func (s Sudoku) _getEligibleMap(rowID int, colID int) EligibleNumbers {
	myMap := standardMap()

	// first get the Row, column and bounded box corresponding to the position
	myRow := s.getRow(rowID)
	myColumn := s.getColumn(colID)
	myBoundedBox := s.getBoundedBox(rowID, colID)

	// scan Row, column and bounded box to eliminate already present numbers
	for _, col := range myRow {
		if col != 0 {
			myMap[col] = false
		}
	}

	for _, col := range myColumn {
		if col != 0 {
			myMap[col] = false
		}
	}

	for _, col := range myBoundedBox {
		if col != 0 {
			myMap[col] = false
		}
	}

	// return the resultant eligible numbers map
	return myMap
}

// make a standard map with all numbers from 1 to 9 as eligible
func standardMap() EligibleNumbers {
	stdMap := make(EligibleNumbers)
	for i := 1; i <= 9; i++ {
		stdMap[i] = true
	}
	return stdMap
}

func (r Row) _getEligibleMap() EligibleNumbers {
	stdMap := standardMap()
	// make the numbers already present as NOT eligible in the map
	for _, col := range r {
		if col != 0 {
			stdMap[col] = false
		}
	}
	return stdMap
}

// nonRepeating validates whether a Row is composed of non repeated numbers
func (r Row) nonRepeating() bool {
	myMap := make(map[int]int)
	for _, col := range r {
		myMap[col] = myMap[col] + 1
		if myMap[col] > 1 {
			return false
		}
	}
	return true
}

// complete validates whether a Row contains all numbers and no zeros
func (r Row) complete() bool {
	e := r._getEligibleMap()
	for _, val := range e {
		if val {
			return false
		}
	}
	return true
}

func (en EligibleNumbers) getSingularEligibleNumber() (eligNum int) {
	eligNumCount := 0
	falseCount := 0

	for key, val := range en {
		if val {
			eligNumCount++
			eligNum = key
		} else {
			falseCount++
		}
	}

	// If exactly one number is eligible, send the corresponding number
	if eligNumCount == 1 {
		return eligNum
	}

	// If no numbers are eligible, then send a different signal -1
	if falseCount == 9 {
		return -1
	}

	// If more than one number is eligible, send 0
	return 0
}

// Print prints the eligible numbers in a map
func (en EligibleNumbers) Print() {
	ea := make([]int, 0)
	for key, val := range en {
		if val {
			ea = append(ea, key)
		}
	}
	fmt.Println(ea)
}

// GetList converts eligible numbers map into a list
func (en EligibleNumbers) GetList() []int {
	ea := make([]int, 0)
	for key, val := range en {
		if val {
			ea = append(ea, key)
		}
	}
	return ea
}
