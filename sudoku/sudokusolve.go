package sudoku

import (
	"errors"
	"sync"
	"sync/atomic"
)

var globalCounter = new(int32)

// Solve takes an unsolved sudoku and solves it. It returns a Sudoku, whether it is solved or not, total number of iterations and error.
// The error can either have a message done, which means everythingno error (or) an actual error message
func Solve(sudokuIn Sudoku) (Sudoku, bool, int, error) {

	/*
		Solve the Sudoku puzzle as follows:
		Perform the following for each cell of the sudoku puzzle
		1) mapper: find out potential numbers that can be filled for each unfilled column in each row by
			looking at the unfilled column from the perspective of the corresponding row, column and the bounded box
		2) reducer: scan through the row, column or bounding box and resolve the column value
		3) repeat step 1 and 2 as long as the number of unfilled reduces in each iteration


		If the unfilled cells are not reducing anymore, do the following

		1. Pick the cell with the least number of potentials:
		2. Fire multiple threads concurrently with each of these potentials filled in the cell.
			Do this recursively.
			I.e. once a cell is filled with a potential, a recursive call is made to solve function,
			which fills the next potential and so on. There can be only one of two outcomes at the top most level:
			(a) the Sudoku is solved
			(b) this combination is invalid, in which case, this guess is abandoned
			at intermediate levels, there can be one of two outcomes:
			(a) the Sudoku is partially solved, in which case, this guessing comtinues
			(b) this combination is invalid, in which case, this guess is abandoned
	*/

	var iteration int
	SudokuOut := sudokuIn.Copy()

	unfilledCount := 0

	for {

		if SudokuOut.Solved() {
			break
		}

		unfilledCount = SudokuOut.UnfilledCount()

		// iteration++
		iteration = int(atomic.AddInt32(globalCounter, 1))

		if iteration >= 10000000 {
			break
		}

		// run across all cells and perform map reduce
		// cells with a single potential number will end up getting filled
		for rowID, row := range SudokuOut {
			for colID, col := range row {
				if col == 0 {
					_cell := SudokuOut.MapEligibleNumbers(rowID, colID)
					_result := SudokuOut.ReduceAndFillEligibleNumber(_cell)
					if _result == -1 {
						// incorrect Sudoku. return to caller
						// this scenario occurs when solve is called recursively with a guess
						return SudokuOut, SudokuOut.Solved(), iteration, errors.New("incorrect Sudoku")
					}
				}
			}
		}

		// If the Sudoku is solved, exit out of the routine
		if SudokuOut.Solved() {
			break
		}

		// If no cells have been reduced, there is no point in repeating
		// Map eligible numbers for each cell and pick a cell with the least number of eligible numbers
		if SudokuOut.UnfilledCount() >= unfilledCount {
			potentials := make(map[int]Cell)
			for rowID, row := range SudokuOut {
				for colID, col := range row {
					if col == 0 {
						_cell := SudokuOut.MapEligibleNumbers(rowID, colID)
						_potentialsLen := len(_cell.eligibleNumbers.GetList())
						potentials[_potentialsLen] = _cell
					}
				}
			}

			// Walk through all cells and pick the cell with the least number of eligible numbers
			var cellToEvaluate Cell
			potentialsRange := []int{2, 3, 4, 5, 6, 7, 8, 9}
			for _, _potential := range potentialsRange {
				if _, ok := potentials[_potential]; ok {
					cellToEvaluate = potentials[_potential]
					break
				}
			}

			// Pick each eligible number, fill it and see if it works
			// Do this CONCURRENTLY to save time
			chanSudokuSolve := make(chan Channel)
			wg := new(sync.WaitGroup)

			for _, eligNum := range cellToEvaluate.eligibleNumbers.GetList() {

				wg.Add(1)

				// Call the solve function recursively, but as a go routine thread so that it executes asynchronously
				go func(sudokuIn Sudoku, rowID int, colID int, fillVal int, wg *sync.WaitGroup, c *chan Channel) {
					defer wg.Done()

					_SudokuOut := sudokuIn.Copy()

					_SudokuOut.Fill(rowID, colID, fillVal)

					sudokuInter, _solved, _iteration, _err := Solve(_SudokuOut)
					*c <- Channel{intermediate: sudokuInter, solved: _solved, iteration: _iteration, err: _err}
				}(SudokuOut, cellToEvaluate.rowID, cellToEvaluate.colID, eligNum, wg, &chanSudokuSolve)
			}

			// wait for the threads to be done & close channel once all threads are done
			go func(wg *sync.WaitGroup, c chan Channel) {
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

				if _err.Error() == "incorrect Sudoku" {
					// This combination is invalid. drop it
				} else {
					// not solved, but the guess is correct. try from beginning
					SudokuOut = _sudokuInter.Copy()
					break
				}
			}

		}
	}

	return SudokuOut, SudokuOut.Solved(), iteration, errors.New("done")
}
