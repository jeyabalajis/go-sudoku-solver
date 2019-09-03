package main

type row []int
type sudoku []row

// EligibleNumbers - A Hashmap to keep track of which numbers are eligible to be filled in a column
type EligibleNumbers map[int]bool

const rowLength int = 9
const colLength int = 9

type cell struct {
	rowID           int
	colID           int
	eligibleNumbers EligibleNumbers
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
	var rowLowerBoundary int
	var rowUpperBoundary int
	var colLowerBoundary int
	var colUpperBoundary int

	switch {
	case rowID <= 2 && colID <= 2:
		rowLowerBoundary = 0
		rowUpperBoundary = 2
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowID <= 5 && colID <= 2:
		rowLowerBoundary = 3
		rowUpperBoundary = 5
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowID <= 8 && colID <= 2:
		rowLowerBoundary = 6
		rowUpperBoundary = 8
		colLowerBoundary = 0
		colUpperBoundary = 2
		break
	case rowID <= 2 && colID <= 5:
		rowLowerBoundary = 0
		rowUpperBoundary = 2
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowID <= 5 && colID <= 5:
		rowLowerBoundary = 3
		rowUpperBoundary = 5
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowID <= 8 && colID <= 5:
		rowLowerBoundary = 6
		rowUpperBoundary = 8
		colLowerBoundary = 3
		colUpperBoundary = 5
		break
	case rowID <= 2 && colID <= 8:
		rowLowerBoundary = 0
		rowUpperBoundary = 2
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	case rowID <= 5 && colID <= 8:
		rowLowerBoundary = 3
		rowUpperBoundary = 5
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	case rowID <= 8 && colID <= 8:
		rowLowerBoundary = 6
		rowUpperBoundary = 8
		colLowerBoundary = 6
		colUpperBoundary = 8
		break
	}

	return rowLowerBoundary, rowUpperBoundary, colLowerBoundary, colUpperBoundary
}

func (s sudoku) solved() bool {
	/*A sudoku is considered solved when
	- there are no empty cells (i.e. cells with number zero)
	- all rows, columns and bounded box contain numbers from 1 to 9 (i.e. complete)
	- there are no repeating numbers in rows, columns or bounded box (i.e. nonRepeating)
	*/

	myColumns := make(map[int]row)
	myBoundedBoxes := make(map[int]row)

	// traverse the sudoku once to collect rows, columns and bounded boxes
	for rowID, row := range s {

		if !(row.complete() && row.nonRepeating()) {
			return false
		}

		for colID, col := range row {

			// collect column values belonging to the same column id in a separate row
			myColumns[colID] = append(myColumns[colID], col)

			// collect column values belonging to the same bounded box into a separate row
			bbID := _getBoundedBoxIndex(rowID, colID)
			myBoundedBoxes[bbID] = append(myBoundedBoxes[bbID], col)
		}
	}

	if len(myColumns) > 0 {
		for _, row := range myColumns {
			if !(row.complete() && row.nonRepeating()) {
				return false
			}
		}
	}

	if len(myBoundedBoxes) > 0 {
		for _, row := range myBoundedBoxes {
			if !(row.complete() && row.nonRepeating()) {
				return false
			}
		}
	}

	return true
}

func (s sudoku) getRow(rowID int) row {
	return s[rowID]
}

func (s sudoku) getColumn(colID int) row {
	var myColumn row

	for _, row := range s {
		for colIndex, col := range row {
			if colID == colIndex {
				myColumn = append(myColumn, col)
			}
		}
	}
	return myColumn
}

func (s sudoku) getBoundedBox(rowID int, colID int) row {
	var myBB row
	rowMin, rowMax, colMin, colMax := _getBoundedBoxBoundaries(rowID, colID)

	for rowIndex, row := range s {
		for colIndex, col := range row {
			if (rowIndex >= rowMin && rowIndex <= rowMax) && (colIndex >= colMin && colIndex <= colMax) {
				myBB = append(myBB, col)
			}
		}
	}
	return myBB
}

// getEligibleMap gets a map of eligible numbers for a particular position
// in the sudoku puzzle
func (s sudoku) getEligibleMap(rowID int, colID int) EligibleNumbers {
	myMap := standardMap()

	// first get the row, column and bounded box corresponding to the position
	myRow := s.getRow(rowID)
	myColumn := s.getColumn(colID)
	myBoundedBox := s.getBoundedBox(rowID, colID)

	// scan row, column and bounded box to eliminate already present numbers
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

	// return the resultany eligible numbers map
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

func (r row) _getEligibleMap() EligibleNumbers {
	stdMap := standardMap()
	// make the numbers already present as NOT eligible in the map
	for _, col := range r {
		if col != 0 {
			stdMap[col] = false
		}
	}
	return stdMap
}

// nonRepeating validates whether a row is composed of non repeated numbers
func (r row) nonRepeating() bool {
	myMap := make(map[int]int)
	for _, col := range r {
		myMap[col] = myMap[col] + 1
		if myMap[col] > 1 {
			return false
		}
	}
	return true
}

// complete validates whether a row contains all numbers and no zeros
func (r row) complete() bool {
	e := r._getEligibleMap()
	for _, val := range e {
		if val {
			return false
		}
	}
	return true
}
