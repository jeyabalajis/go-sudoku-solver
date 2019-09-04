package main

import "fmt"

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

		fire up reducer for all 81 columns concurrently to fill up the cell if there is a singular eligible number

		in the reducer, each concurrent thread will be working on non overlapping columns and fill up the cell
	*/

	sudokuOut = sudokuIn
	mapResults := make([]cell, 0)
	unfilledCount := 0
	iteration := 0

	for {
		// If the sudoku is solved, exit out of the routine
		if sudokuOut.solved() {
			break
		}

		unfilledCount = sudokuOut.unfilledCount()
		c := make(chan cell)

		iteration++

		// call map function concurrently for all the cells
		for rowID, row := range sudokuOut {
			for colID, col := range row {
				if col == 0 {
					go sudokuOut.mapEligibleNumbers(rowID, colID, c)
					myCell := <-c
					mapResults = append(mapResults, myCell)
				}
			}
		}

		// call reduce/fill function concurrently for all the cells
		for _, _cell := range mapResults {
			c := make(chan bool)
			go sudokuOut.fillEligibleNumber(_cell, c)
			myResult := <-c
			if myResult {
				fmt.Println(_cell.rowID, _cell.colID, myResult)
			}
		}

		// If no cells have been reduced, there is no point in repeating, start brute force
		if sudokuOut.unfilledCount() >= unfilledCount {
			fmt.Println("giving up!")
			break
		}
		fmt.Println("iteration: ", iteration)
	}

	return sudokuOut, sudokuOut.solved(), nil
}
