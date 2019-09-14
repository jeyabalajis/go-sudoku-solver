package main

import (
	"errors"
	"sync"
	"sync/atomic"
)

var globalCounter = new(int32)

// Take an unsolved sudoku input and return a solved sudoku output
func solve(sudokuIn sudoku) (sudoku, bool, int, error) {

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

	var iteration int
	sudokuOut := sudokuIn.copy()

	// mapResults := make([]cell, 0)
	unfilledCount := 0

	for {

		// If the sudoku is solved, exit out of the routine
		if sudokuOut.solved() {
			break
		}

		unfilledCount = sudokuOut.unfilledCount()

		// iteration++
		iteration = int(atomic.AddInt32(globalCounter, 1))

		if iteration >= 10000000 {
			break
		}

		// run across all cells and perform map reduce
		// cells with a single potential will end up getting filled
		for rowID, row := range sudokuOut {
			for colID, col := range row {
				if col == 0 {
					_cell := sudokuOut.mapEligibleNumbers(rowID, colID)
					_result := sudokuOut.reduceAndFillEligibleNumber(_cell)
					if _result == -1 {
						// incorrect sudoku. return to caller
						// this scenario occurs when solve is called recursively with a guess
						return sudokuOut, sudokuOut.solved(), iteration, errors.New("incorrect sudoku")
					}
				}
			}
		}

		// If the sudoku is solved, exit out of the routine
		if sudokuOut.solved() {
			break
		}

		// If no cells have been reduced, there is no point in repeating
		// Map eligible numbers for each cell and pick a cell with the least number of eligible numbers
		if sudokuOut.unfilledCount() >= unfilledCount {
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

			// Walk through all cells and pick the cell with the least number of eligible numbers
			var cellToEvaluate cell
			potentialsRange := []int{2, 3, 4, 5, 6, 7, 8, 9}
			for _, _potential := range potentialsRange {
				if _, ok := potentials[_potential]; ok {
					cellToEvaluate = potentials[_potential]
					break
				}
			}

			// Pick each eligible number, fill it and see if it works
			// Do this CONCURRENTLY to save time
			chanSudokuSolve := make(chan sudokuChannel)
			wg := new(sync.WaitGroup)

			for eligNum, eligible := range cellToEvaluate.eligibleNumbers {

				if eligible {
					wg.Add(1)

					// Call the solve function recursively, but as a go routine thread so that it executes asynchronously
					go func(sudokuIn sudoku, rowID int, colID int, fillVal int, wg *sync.WaitGroup, c *chan sudokuChannel) {
						defer wg.Done()

						_sudokuOut := sudokuIn.copy()

						_sudokuOut.fill(rowID, colID, fillVal)

						sudokuInter, _solved, _iteration, _err := solve(_sudokuOut)
						*c <- sudokuChannel{intermediate: sudokuInter, solved: _solved, iteration: _iteration, err: _err}
					}(sudokuOut, cellToEvaluate.rowID, cellToEvaluate.colID, eligNum, wg, &chanSudokuSolve)
				}
			}

			// wait for the threads to be done & close channel once all threads are done
			go func(wg *sync.WaitGroup, c chan sudokuChannel) {
				wg.Wait()
				close(c)
			}(wg, chanSudokuSolve)

			// collect the results and look for the right guess
			for r := range chanSudokuSolve {
				_sudokuInter := r.intermediate
				_solved := r.solved
				_err := r.err
				iteration = iteration + r.iteration

				if _solved {
					return _sudokuInter, _solved, iteration, _err
				}

				if _err.Error() == "incorrect sudoku" {
					// This combination is invalid. drop it
				} else {
					// not solved, but the guess is correct. try from beginning
					sudokuOut = _sudokuInter.copy()
					break
				}
			}

		}
	}

	return sudokuOut, sudokuOut.solved(), iteration, errors.New("done")
}
