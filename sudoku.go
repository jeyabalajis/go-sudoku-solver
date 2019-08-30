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

func (s sudoku) solved() bool {
	/*A sudoku is considered solved when
	- there are no empty cells (i.e. cells with number zero)
	- all rows, columns and bounded box contain numbers from 1 to 9
	- there are no overlapping numbers in rows, columns or bounded box
	*/

	myColumns := make(map[int]row)
	myBoundedBoxes := make(map[int]row)

	// traverse the sudoku once to collect rows, columns and bounded boxes
	for rowID, row := range s {
		for colID, col := range row {
			// No number shall be zero in a solved sudoku
			if col == 0 {
				return false
			}

			// collect column values belonging to the same column id in a separate row
			myColumns[colID] = append(myColumns[colID], col)

			// collect column values belonging to the same bounded box into a separate row
			bbID := _getBoundedBoxIndex(rowID, colID)
			myBoundedBoxes[bbID] = append(myBoundedBoxes[bbID], col)
		}

		if !(row.complete() && row.nonRepeating()) {
			return false
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

// getEligibleMap gets a map of eligible numbers for a particular position
// in the sudoku puzzle
func (s sudoku) getEligibleMap(rowID int, colID int) EligibleNumbers {
	myMap := standardMap()

	// first get the row, column and bounded box corresponding to the position

	// scan row, column and bounded box to eliminate already present numbers

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
