package main

import (
	"errors"
)

// Take an unsolved sudoku input and return a solved sudoku output
func solve(sudokuIn sudoku, iter ...int) (sudokuOut sudoku, solved bool, iteration int, err error) {

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
	// var input string
	sudokuOut = sudokuIn.copy()
	mapResults := make([]cell, 0)
	unfilledCount := 0

	if iter != nil {
		iteration = iter[0]
	} else {
		iteration = 0
	}

	for {

		// If the sudoku is solved, exit out of the routine
		if sudokuOut.solved() {
			//fmt.Println("sudoku solved!")
			break
		}

		if iteration >= 1000000 {
			break
		}

		unfilledCount = sudokuOut.unfilledCount()

		// cMap := make(chan cell)
		// cReduce := make(chan int)

		iteration++

		//fmt.Println("<<<Iteration & unfilled>>>: ", iteration, sudokuOut.unfilledCount())

		// fmt.Scanln(&input)

		// call map function concurrently for all the cells
		// mapResults = make([]cell, 0)
		for rowID, row := range sudokuOut {
			for colID, col := range row {
				if col == 0 {
					_cell := sudokuOut.mapEligibleNumbers(rowID, colID)
					_result := sudokuOut.fillEligibleNumber(_cell)
					if _result == -1 {
						//fmt.Println("incorrect sudoku. return to caller")
						return sudokuOut, sudokuOut.solved(), iteration, errors.New("incorrect sudoku")
					}
				}
			}
		}

		// If the sudoku is solved, exit out of the routine
		if sudokuOut.solved() {
			//fmt.Println("sudoku solved!")
			break
		}

		// If no cells have been reduced, there is no point in repeating, start brute force
		if sudokuOut.unfilledCount() >= unfilledCount {
			//fmt.Println("start brute force attack")

			mapResults = make([]cell, 0)
			for rowID, row := range sudokuOut {
				for colID, col := range row {
					if col == 0 {
						_cell := sudokuOut.mapEligibleNumbers(rowID, colID)
						mapResults = append(mapResults, _cell)
					}
				}
			}

			stopSearching := false
			for _, _cell := range mapResults {
				// Pick each eligible number, fill it and see if it works
				for eligNum, val := range _cell.eligibleNumbers {
					if val {
						// fmt.Printf("try out value %v in index (%v, %v)", eligNum, _cell.rowID, _cell.colID)
						sudokuCopy := make(sudoku, len(sudokuOut))
						copy(sudokuCopy, sudokuOut)

						sudokuOut[_cell.rowID][_cell.colID] = eligNum
						sudokuInter, _solved, _iteration, _err := solve(sudokuOut, iteration)
						iteration = _iteration

						if _solved {
							//fmt.Println("solved. return to caller")
							return sudokuInter, _solved, iteration, _err
						}

						if _err != nil {
							//fmt.Println("incorrect sudoku, try the next one")
							// rollback the assignment and continue searching
							sudokuOut = make(sudoku, len(sudokuCopy))
							copy(sudokuOut, sudokuCopy)
						} else {
							//fmt.Println("not solved, but the guess is correct. try from beginning")
							sudokuOut = make(sudoku, len(sudokuCopy))
							copy(sudokuOut, sudokuInter)
							stopSearching = true
						}
					}

					// sudokuOut.print()

					if stopSearching {
						// //fmt.Println("break search")
						break
					}
				}

				// //fmt.Println("break search")
				break
			}
		}
	}

	//fmt.Println("finally going back")
	return sudokuOut, sudokuOut.solved(), iteration, nil
}
