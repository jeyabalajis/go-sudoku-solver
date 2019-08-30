package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
	Utility to convert a fixed file into a 9x9 sudoku
*/

func _convertStringToRow(str []string) row {
	var myRow row

	for _, singleStr := range str {
		intVal, err := strconv.Atoi(singleStr)
		if err != nil {
			fmt.Println("non integer provided")
			intVal = 0
		}

		myRow = append(myRow, intVal)
	}

	return myRow
}

func splitString(str string) []string {
	cleaned := strings.Replace(str, ",", " ", -1)
	strSlice := strings.Fields(cleaned)
	return strSlice
}

func newSudokuFromFile(fileName string) sudoku {
	var mySudoku sudoku

	file, err := os.Open(fileName)

	if err != nil {
		log.Fatalf("failed opening file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		strSlice := splitString(scanner.Text())
		// fmt.Println(strSlice)
		myRow := _convertStringToRow(strSlice)
		mySudoku = append(mySudoku, myRow)
	}

	return mySudoku
}
