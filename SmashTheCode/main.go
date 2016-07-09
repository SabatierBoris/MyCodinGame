package main

import "fmt"

//import "os"

/**
 * Auto-generated code below aims at helping you parse
 * the standard input according to the problem statement.
 **/

type Dot struct {
	Value int
	Up    *Dot
	Down  *Dot
	Left  *Dot
	Right *Dot
}

func (d *Dot) SetUp(up *Dot) {
	d.Up = up
	if up.Down != d {
		up.SetDown(d)
	}
}

func (d *Dot) SetDown(down *Dot) {
	d.Down = down
	if down.Up != d {
		down.SetUp(d)
	}
}

func (d *Dot) SetLeft(left *Dot) {
	d.Left = left
	if left.Right != d {
		left.SetRight(d)
	}
}

func (d *Dot) SetRight(right *Dot) {
	d.Right = right
	if right.Left != d {
		right.SetLeft(d)
	}
}

type BoardGrid struct {
	columns []Dot
}

func NewBoardGrid(size_x int, size_y int) *BoardGrid {
	var tops []Dot
	tops = make([]Dot, size_x, size_x)
	return nil
}

func (b *BoardGrid) initialize()

func main() {
	for {
		for i := 0; i < 8; i++ {
			// colorA: color of the first block
			// colorB: color of the attached block
			var colorA, colorB int
			fmt.Scan(&colorA, &colorB)
		}
		var score1 int
		fmt.Scan(&score1)

		for i := 0; i < 12; i++ {
			var row string
			fmt.Scan(&row)
		}
		var score2 int
		fmt.Scan(&score2)

		for i := 0; i < 12; i++ {
			// row: One line of the map ('.' = empty, '0' = skull block, '1' to '5' = colored block)
			var row string
			fmt.Scan(&row)
		}

		// fmt.Fprintln(os.Stderr, "Debug messages...")
		fmt.Printf("0\n") // "x": the column in which to drop your blocks
	}
}
