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
		3) repeat step 1 and 2 as long as the number of unfilled reduces in each iteration
		4) If the unfilled cells are not reducing anymore, do the following
			pick the cell with the least number of potentials:
			fire multiple threads concurrently with each of these potentials filled in the cell
			do this recursively.

			I.e. once a cell is filled with a potential, a recursive call is made to solve function,
			which fills the next potential and so on. There can be only one of two outcomes at the top most level:
			(a) the sudoku is solved
			(b) this combination is invalid, in which case, this guess is abandoned

			at intermediate levels, there can be one of two outcomes:
			(a) the sudoku is partially solved, in which case, this guessing comtinues
			(b) this combination is invalid, in which case, this guess is abandoned

	*/

	sudokuOut = sudokuIn.copy()
	unfilledCount := 0

	if iter != nil {
		iteration = iter[0]
	} else {
		iteration = 0
	}

	// fmt.Println(iteration)

	for {

		// If the sudoku is solved, exit out of the routine
		if sudokuOut.solved() {
			//fmt.Println("sudoku solved!")
			break
		}

		if iteration >= 10000000 {
			break
		}

		unfilledCount = sudokuOut.unfilledCount()

		// cMap := make(chan cell)
		// cReduce := make(chan int)

		iteration++

		//fmt.Println("<<<Iteration & unfilled>>>: ", iteration, sudokuOut.unfilledCount())

		// sudokuOut.print()

		// call map function concurrently for all the cells
		// mapResults = make([]cell, 0)
		for rowID, row := range sudokuOut {
			for colID, col := range row {
				if col == 0 {
					_cell := sudokuOut.mapEligibleNumbers(rowID, colID)
					_result := sudokuOut.reduceAndFillEligibleNumber(_cell)
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

			potentials := make(map[int]cell)
			for rowID, row := range sudokuOut {
				for colID, col := range row {
					if col == 0 {
						_cell := sudokuOut.mapEligibleNumbers(rowID, colID)
						_potentialsLen := len(_cell.eligibleNumbers.getList())
						potentials[_potentialsLen] = _cell
					}
				}
			}

			// var input string
			// fmt.Println(potentials)
			// fmt.Scanln(&input)

			// Walk through all cells and group them by the number of potentials
			var cellToEvaluate cell
			potentialsRange := []int{2, 3, 4, 5, 6, 7, 8, 9}
			for _, _potential := range potentialsRange {
				if _, ok := potentials[_potential]; ok {
					cellToEvaluate = potentials[_potential]
					break
				}
			}

			// Pick each eligible number, fill it and see if it works
			for eligNum, eligible := range cellToEvaluate.eligibleNumbers {

				if eligible {

					sudokuCopy := sudokuOut.copy()

					sudokuOut.fill(cellToEvaluate.rowID, cellToEvaluate.colID, eligNum)

					_sudokuInter, _solved, _iteration, _err := solve(sudokuOut, iteration)

					if _solved {
						// fmt.Println("solved. return to caller")
						return _sudokuInter, _solved, _iteration, _err
					}

					if _err.Error() == "incorrect sudoku" {
						// This combination is invalid. rollback. try out the next option
						sudokuOut = sudokuCopy.copy()
					} else {
						//fmt.Println("not solved, but the guess is correct. try from beginning")
						sudokuOut = _sudokuInter.copy()
						sudokuOut.print()
						break
					}

				}
			}
		}
	}

	//fmt.Println("finally going back")
	return sudokuOut, sudokuOut.solved(), iteration, errors.New("done")
}
