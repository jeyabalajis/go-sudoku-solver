# go-sudoku-solver
A __concurrent__ & __recursive__ sudoku solver written in golang!

# Motivation
This is a small project that showcases the power of go routines by solving a sudoku puzzle. This solver solves
arguably [the world's hardest sudoku ever](https://curiosity.com/topics/a-finnish-mathematician-claimed-that-this-is-the-most-difficult-sudoku-puzzle-in-the-world-curiosity/) in around 3 seconds.

# Solution Approach

The solver employs two different sub-approaches to solve a sudoku

1. The solver first maps the set of eligible numbers (_that can be filled_) for each cell. It also _fills_ the cells which have just a single eligible number. This is a **map and reduce** performed on each cell.
2. If each iteration keeps reducing the number of _unfilled_ cells, the solver keeps repeating **map and reduce** until the sudoku is solved.


3. When the number of _unfilled_ cells stops reducing, the solver switches to a brute force **trial and error** approach
4. Under this approach, the cell with *least number of eligible numbers* is identified.
5. The solver fires **a go routine as a recursive call to itself** for each of these eligible numbers after filling the cell with the number.
6. In the next iteration, the next cell with least number of eligible numbers is picked up, so on and so forth
7. This keeps repeating until either the sudoku is solved or a total of 10000000 (ten million) iterations are exhausted.

> The approach results in a mega decision tree with each guess on a cell as a node. There is only one path which correctly solves
> the sudoku. When the solver hits this path, it returns the solved sudoku.

> The sudoku is represented as a slice of int slice. A custom Copy method is written to perform deep copy of a sudoku. This is important
> since each thread of solver must work on a distinct copy of a sudoku to avoid **data race**. 

# Results

Here's a snapshot of the solver run on arguably the world's hardest sudoku ever!

![Sudoku Solver results](/images/sudoku_run_results.png)