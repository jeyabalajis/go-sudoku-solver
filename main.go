package main

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	mySudoku := newSudokuFromFile("tests/hardest_3.txt")

	mySudoku.print()
	// log.Println("start filling")
	start := time.Now()
	solvedSudoku, solved, iterations, err := solve(mySudoku)
	elapsed := time.Since(start)
	fmt.Println("solved:", solved)
	fmt.Println("error:", err)
	fmt.Println("total iterations: ", iterations)
	fmt.Println("elapsed time: ", elapsed)
	solvedSudoku.print()

	// chanSudokuSolve := make(chan sudokuChannel)
	// wg := new(sync.WaitGroup)

	// cellPotentials := []int{2, 3, 4, 5, 6, 7, 8, 9}
	// for _, potential := range cellPotentials {
	// 	wg.Add(1)

	// 	_sudokuTemp := make(sudoku, len(mySudoku))
	// 	copy(_sudokuTemp, mySudoku)
	// 	go _fillWrapper(_sudokuTemp, 0, 0, potential, wg, &chanSudokuSolve)
	// }

	// go func(wg *sync.WaitGroup, c chan sudokuChannel) {
	// 	log.Println("waiting")
	// 	wg.Wait()
	// 	log.Println("done waiting")
	// 	close(c)
	// }(wg, chanSudokuSolve)
	// for r := range chanSudokuSolve {
	// 	_sudokuInter := r.intermediate

	// 	fmt.Println("*****")
	// 	_sudokuInter.print()
	// 	fmt.Println("*****")
	// }
}

func _fillWrapper(sudokuIn sudoku, rowID int, colID int, fillVal int, wg *sync.WaitGroup, c *chan sudokuChannel) {
	defer wg.Done()

	_sudokuOut := sudokuIn.copy()

	done := make(chan struct{})
	go func() {
		_sudokuOut[rowID][colID] = fillVal
		done <- struct{}{}
	}()
	<-done

	*c <- sudokuChannel{intermediate: _sudokuOut, solved: true, iteration: 0, err: errors.New("done")}
	log.Println("sent solution")
}
