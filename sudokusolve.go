package main

import (
	"errors"
	"log"
	"sync"
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
	// mapResults := make([]cell, 0)
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

			chanSudokuSolve := make(chan sudokuChannel)
			wg := new(sync.WaitGroup)

			// Pick each eligible number, fill it and see if it works
			for eligNum, eligible := range cellToEvaluate.eligibleNumbers {

				if eligible {

					sudokuOut.fill(cellToEvaluate.rowID, cellToEvaluate.colID, eligNum)

					wg.Add(1)
					go _solveWrapper(sudokuOut, 0, cellToEvaluate, eligNum, wg, &chanSudokuSolve)

					// go _solveConcurrent(sudokuOut, chanSudokuSolve, iteration)

					// r := <-chanSudokuSolve
					// sudokuInter := r.intermediate
					// _solved := r.solved
					// _err := r.err
					// iteration = iteration + r.iteration

					// if _solved {
					// 	//fmt.Println("solved. return to caller")
					// 	return sudokuInter, _solved, iteration, _err
					// }

					// if _err != nil {
					// 	// This combination is invalid. drop it
					// } else {
					// 	//fmt.Println("not solved, but the guess is correct. try from beginning")
					// 	sudokuOut = make(sudoku, len(sudokuInter))
					// 	copy(sudokuOut, sudokuInter)
					// 	break
					// }
				}
			}

			go func(wg *sync.WaitGroup, c chan sudokuChannel) {
				log.Println("waiting")
				wg.Wait()
				log.Println("done waiting")
				close(c)
			}(wg, chanSudokuSolve)

			// wg.Wait()
			// close(chanSudokuSolve)

			// collect the results and look for the right guess
			log.Println("look at the results")
			for r := range chanSudokuSolve {
				sudokuInter := r.intermediate
				_solved := r.solved
				_err := r.err
				iteration = iteration + r.iteration
				_cell := r.cellMutated
				_valueOption := r.valueOption

				log.Println(_cell.rowID, _cell.colID, _valueOption, _solved, _err.Error())
				if _solved {
					// fmt.Println("solved. return to caller")
					return sudokuInter, _solved, iteration, _err
				}

				if _err.Error() == "incorrect sudoku" {
					// This combination is invalid. drop it
				} else {
					//fmt.Println("not solved, but the guess is correct. try from beginning")
					sudokuOut = make(sudoku, len(sudokuInter))
					copy(sudokuOut, sudokuInter)
					break
				}
			}

		}
	}

	//fmt.Println("finally going back")
	return sudokuOut, sudokuOut.solved(), iteration, errors.New("done")
}

func _solveWrapper(sudokuIn sudoku, iter int, cE cell, mutatedValue int, wg *sync.WaitGroup, c *chan sudokuChannel) {
	defer wg.Done()
	sudokuInter, _solved, _iteration, _err := solve(sudokuIn, iter)
	*c <- sudokuChannel{intermediate: sudokuInter, solved: _solved, iteration: _iteration, err: _err, cellMutated: cE, valueOption: mutatedValue}
	log.Println("sent solution")
}
