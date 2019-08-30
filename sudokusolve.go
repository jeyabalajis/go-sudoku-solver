package main

// Take an unsolved sudoku input and return a solved sudoku output
func solve(sudokuIn sudoku) (sudokuOut sudoku, solved bool, err error) {

	/*
		Solve the sudoku puzzle as follows:
		1) mapper: find out potential numbers that can be filled for each unfilled column in each row by
			looking at the unfilled column from the perspective of the corresponding row, column and the bounded box
		2) reducer: scan through the row, column or bounding box and resolve the column value
		3) repeat step 1 and 2 until sudoku is fully solved or the program completes 100000 tries, whichever comes first

		concurrency:
		fire up mapper for all 81 columns concurrently since nothing is going to be filled. Just the eligible numbers
		would be mapped in this run for the column passed.

		fire up reducer by bounding boxes concurrently (total nine concurrent threads)
		fire up reducer by rows concurrently (total nine concurrent threads)
		fire up reducer by columns concurrently (total nine concurrent threads)

		in the reducer, each concurrent thread will be working on non overlapping columns and fill up / firm up the values
		and reduce eligible numbers
	*/

	sudokuOut = sudokuIn

	for {
		if sudokuOut.solved() {
			break
		}

	}

	return sudokuOut, sudokuOut.solved(), nil
}

// func (s sudoku) _map(rowID int, colID int) {
// 	eligibleNumbers := make(map[int]int)

// 	// traverse through the row and collect eligible numbers
// 	for _, col := range s[rowID] {

// 	}
// 	// traverse through the column and collect eligible numbers

// 	// traverse through the bounded box and collect eligible numbers

// 	// prepare the final set of eligible numbers and send it
// }
