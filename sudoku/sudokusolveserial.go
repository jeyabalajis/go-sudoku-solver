package sudoku

import (
	"errors"
)

// Take an unsolved Sudoku input and return a solved Sudoku output
func solveSerial(sudokuIn Sudoku, iter ...int) (sudokuOut Sudoku, solved bool, iteration int, err error) {

	/*
		Solve the Sudoku puzzle as follows:
		1) mapper: find out potential numbers that can be filled for each unfilled column in each row by
			looking at the unfilled column from the perspective of the corresponding row, column and the bounded box
		2) reducer: scan through the row, column or bounding box and resolve the column value
		3) repeat step 1 and 2 as long as the number of unfilled reduces in each iteration
		4) If the unfilled Cells are not reducing anymore, do the following
			pick the Cell with the least number of potentials:
			fire multiple threads concurrently with each of these potentials filled in the Cell
			do this recursively.

			I.e. once a Cell is filled with a potential, a recursive call is made to solve function,
			which fills the next potential and so on. There can be only one of two outcomes at the top most level:
			(a) the Sudoku is solved
			(b) this combination is invalid, in which case, this guess is abandoned

			at intermediate levels, there can be one of two outcomes:
			(a) the Sudoku is partially solved, in which case, this guessing comtinues
			(b) this combination is invalid, in which case, this guess is abandoned

	*/

	sudokuOut = sudokuIn.Copy()
	UnfilledCount := 0

	if iter != nil {
		iteration = iter[0]
	} else {
		iteration = 0
	}

	// fmt.Println(iteration)

	for {

		// If the Sudoku is solved, exit out of the routine
		if sudokuOut.Solved() {
			//fmt.Println("Sudoku solved!")
			break
		}

		if iteration >= 10000000 {
			break
		}

		UnfilledCount = sudokuOut.UnfilledCount()

		// cMap := make(chan Cell)
		// cReduce := make(chan int)

		iteration++

		//fmt.Println("<<<Iteration & unfilled>>>: ", iteration, sudokuOut.UnfilledCount())

		// sudokuOut.Print(0

		// call map function concurrently for all the Cells
		// mapResults = make([]Cell, 0)
		for rowID, row := range sudokuOut {
			for colID, col := range row {
				if col == 0 {
					_cell := sudokuOut.MapEligibleNumbers(rowID, colID)
					_result := sudokuOut.ReduceAndFillEligibleNumber(_cell)
					if _result == -1 {
						//fmt.Println("incorrect Sudoku. return to caller")
						return sudokuOut, sudokuOut.Solved(), iteration, errors.New("incorrect Sudoku")
					}
				}
			}
		}

		// If the Sudoku is solved, exit out of the routine
		if sudokuOut.Solved() {
			//fmt.Println("Sudoku solved!")
			break
		}

		// If no Cells have been reduced, there is no point in repeating, start brute force
		if sudokuOut.UnfilledCount() >= UnfilledCount {
			//fmt.Println("start brute force attack")

			potentials := make(map[int]Cell)
			for rowID, row := range sudokuOut {
				for colID, col := range row {
					if col == 0 {
						_cell := sudokuOut.MapEligibleNumbers(rowID, colID)
						_potentialsLen := len(_cell.eligibleNumbers.GetList())
						potentials[_potentialsLen] = _cell
					}
				}
			}

			// var input string
			// fmt.Println(potentials)
			// fmt.Scanln(&input)

			// Walk through all Cells and group them by the number of potentials
			var cellToEvaluate Cell
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

					SudokuCopy := sudokuOut.Copy()

					sudokuOut.Fill(cellToEvaluate.rowID, cellToEvaluate.colID, eligNum)

					_sudokuInter, _solved, _iteration, _err := solveSerial(sudokuOut, iteration)

					if _solved {
						// fmt.Println("solved. return to caller")
						return _sudokuInter, _solved, _iteration, _err
					}

					if _err.Error() == "incorrect Sudoku" {
						// This combination is invalid. rollback. try out the next option
						sudokuOut = SudokuCopy.Copy()
					} else {
						//fmt.Println("not solved, but the guess is correct. try from beginning")
						sudokuOut = _sudokuInter.Copy()
						sudokuOut.Print()
						break
					}

				}
			}
		}
	}

	//fmt.Println("finally going back")
	return sudokuOut, sudokuOut.Solved(), iteration, errors.New("done")
}
