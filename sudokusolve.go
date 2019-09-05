package main

import (
	"errors"
	"fmt"
)

// Take an unsolved sudoku input and return a solved sudoku output
func solve(sudokuIn sudoku, iter ...int) (sudokuOut sudoku, solved bool, err error) {

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

	if iter != nil {
		iteration = iter[0]
	}

	fmt.Println("<<<Iteration: ", iteration)

	for {
		// If the sudoku is solved, exit out of the routine
		if sudokuOut.solved() {
			break
		}

		if iteration >= 1000 {
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
			c := make(chan int)
			go sudokuOut.fillEligibleNumber(_cell, c)
			myResult := <-c
			if myResult == -1 {
				fmt.Println("incorrect sudoku")
				return sudokuOut, sudokuOut.solved(), errors.New("incorrect sudoku")
			}
		}

		// If no cells have been reduced, there is no point in repeating, start brute force
		if sudokuOut.unfilledCount() >= unfilledCount {
			// fmt.Println("start brute force attack")
			breakLoop := false
			for _, _cell := range mapResults {
				// Pick each eligible number, fill it and see if it works
				for eligNum, val := range _cell.eligibleNumbers {
					if val {
						sudokuCopy := sudokuOut.copy()

						sudokuOut[_cell.rowID][_cell.colID] = eligNum
						_, _solved, _err := solve(sudokuOut, iteration)

						if _solved || _err == nil {
							breakLoop = true
						} else {
							if _err != nil && _err.Error() == "incorrect sudoku" {
								// rollback the assignment and continue searching
								sudokuOut = sudokuCopy.copy()
							}
						}
					}

					if breakLoop {
						break
					}
				}

				if breakLoop {
					break
				}
			}
		}
	}

	return sudokuOut, sudokuOut.solved(), nil
}
