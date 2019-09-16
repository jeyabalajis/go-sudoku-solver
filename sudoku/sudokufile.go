package sudoku

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

/*
	Utility to convert a fixed file into a 9x9 Sudoku
*/

func _convertStringToRow(str []string) Row {
	var myRow Row

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

// NewSudokuFromFile reads the file with fileName from tests folder and converts it into a Sudoku
func NewSudokuFromFile(fileName string) Sudoku {
	var mySudoku Sudoku

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
