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

	sudokuOut = sudokuIn.copy()
	mapResults := make([]cell, 0)
	unfilledCount := 0
	iteration := 0

	if iter != nil {
		iteration = iter[0]
	}

	for {
		// If the sudoku is solved, exit out of the routine
		if sudokuOut.solved() {
			break
		}

		if iteration >= 1000 {
			break
		}

		if iteration%50 == 0 {
			sudokuOut.print()
		}

		unfilledCount = sudokuOut.unfilledCount()

		// cMap := make(chan cell)
		// cReduce := make(chan int)

		iteration++

		fmt.Println("<<<Iteration>>>: ", iteration)

		fmt.Println("map")
		// call map function concurrently for all the cells
		mapResults = make([]cell, 0)
		for rowID, row := range sudokuOut {
			for colID, col := range row {
				if col == 0 {
					_cell := sudokuOut.mapEligibleNumbers(rowID, colID)
					mapResults = append(mapResults, _cell)
				}
			}
		}

		fmt.Println("reduce")
		// call reduce/fill function concurrently for all the cells
		for _, _cell := range mapResults {
			fmt.Println(_cell)
			_result := sudokuOut.fillEligibleNumber(_cell)
			if _result == -1 {
				fmt.Println("incorrect sudoku")
				return sudokuOut, sudokuOut.solved(), errors.New("incorrect sudoku")
			}
		}

		// If no cells have been reduced, there is no point in repeating, start brute force
		if sudokuOut.unfilledCount() >= unfilledCount {
			fmt.Println("start brute force attack")
			breakLoop := false
			for _, _cell := range mapResults {
				// Pick each eligible number, fill it and see if it works
				for eligNum, val := range _cell.eligibleNumbers {
					if val {
						sudokuCopy := sudokuOut.copy()

						sudokuOut[_cell.rowID][_cell.colID] = eligNum
						_, _solved, _err := solve(sudokuOut, iteration)

						if _err != nil && _err.Error() == "incorrect sudoku" {
							// rollback the assignment and continue searching
							sudokuOut = sudokuCopy.copy()
						}

						if _solved {
							breakLoop = true
						}

						if _err == nil {
							breakLoop = true
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
