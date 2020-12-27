// v1.4. Change thinking from x y to row column to solve problems.
package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	term "github.com/nsf/termbox-go"
)

var score = 0
var level = 1
var lives = 1
var screwdrivers = 1
var teleports = 1
var board = [10][20]string{} //Game Board Array. Flip from columns rows to 10 rows (y) 20 columns (x)
var row, col, randrow, randcol, whorow, whocol, dalekrow, dalekcol int
var numrows = 10 // Number of rows in the gameboard
var numcols = 20 // Number of columns in the gameboard. Typical spacing is twice as high as wide
var daleks = 5
var alivedaleks = 5
var dalekcount = 5
var dalekscoord = [50][3]int{} // Max 50 Daleks with columns/properties of row,column,Dead.
// Start with 5 Daleks and increment in each level

func reset() {
	term.Sync() // cosmestic purpose
}

func initboard() { // Fill entire board with dots to define the playing field
	for row = 0; row < numrows; row++ {
		for col = 0; col < numcols; col++ {
			board[row][col] = "."
		}
	}
}

func initdaleks() {
	daleks = level * 5
	dalekcount = daleks
	var i = 0
	for i < daleks {
		var seednum = rand.NewSource(time.Now().UnixNano())
		var rnum = rand.New(seednum)
		randrow = rnum.Intn(numrows)
		randcol = rnum.Intn(numcols)
		dalekscoord[i][0] = randrow // Row (y)
		dalekscoord[i][1] = randcol // Column (x)
		dalekscoord[i][2] = 0       // Not Dead Yet
		board[randrow][randcol] = "D"
		i = i + 1
	}
}

func initwho() {
	// Randomize Doctor Who ("W")
	var seednum = rand.NewSource(time.Now().UnixNano())
	var rnum = rand.New(seednum)
	randrow = rnum.Intn(numrows)
	randcol = rnum.Intn(numcols)
	// Check to see if we landed on  Dalek
	if board[randrow][randcol] == "D" {
		rnum = rand.New(seednum)
		randrow = rnum.Intn(numrows)
		randcol = rnum.Intn(numcols)
		board[randrow][randcol] = "W"
	} else {
		board[randrow][randcol] = "W"
	}
	board[randrow][randcol] = "W"
	whorow = randrow
	whocol = randcol
}

func drawboard() { // Draw the board again after changing something
	if dalekcount == 0 { // Killed all Daleks, level up
		level = level + 1
		lives = lives + 1
		screwdrivers = screwdrivers + 1
		teleports = teleports + 1

		initboard()
		initdaleks()
		initwho()
		// fmt.Println("YOU WIN!!, ESC=Quit")
	}
	for row = 0; row < numrows; row++ {
		for col = 0; col < numcols; col++ {
			fmt.Printf("%s", board[row][col])
			if col == (numcols - 1) {
				fmt.Printf("\n")
			}
		}
	}
	fmt.Println("Arrows=Move, .=Wait, t=Teleport, s=Screwdriver, ESC=Quit")
	fmt.Printf("Lives: %d Sonic Screwdrivers: %d Teleports: %d Daleks: %d\n", lives, screwdrivers, teleports, dalekcount)
	fmt.Printf("Level: %d Score: %d\n", level, score)
}

func dalekonwho() int {
	var crash = 0
	var i = 0
	for i < daleks {
		if dalekscoord[i][0] == whorow && dalekscoord[i][1] == whocol {
			crash = 1
		}
		i = i + 1
	}
	return crash
}

func dalekondalek(daleknum int) int { // Not sure if this is completely accurate yet. Might need to move ALL Daleks and then check crashes
	// Move the check out of the for loop to after the for loop
	var crash = 0
	var i = 0
	for i < daleks {
		if dalekscoord[i][0] == dalekscoord[daleknum][0] && dalekscoord[i][1] == dalekscoord[daleknum][1] && i != daleknum {
			crash = 1
			dalekscoord[i][2] = 1
			dalekscoord[daleknum][2] = 1
			board[dalekscoord[i][0]][dalekscoord[i][1]] = "*"
		}
		i = i + 1
	}
	return crash
}

func countdaleks() int {
	var i = 0
	var remaining = daleks
	for i < daleks {
		if dalekscoord[i][2] == 1 {
			remaining = remaining - 1
		}
		i = i + 1
	}
	return remaining
}

func blastdaleks() {
	var i = 0
	for i < daleks {
		var dalekrow = dalekscoord[i][0]
		var dalekcol = dalekscoord[i][1]
		if ((dalekrow == whorow) && (dalekcol == whocol+1) || dalekcol == whocol-1) ||
			((dalekcol == whocol) && (dalekrow == whorow+1) || dalekrow == whorow-1) {
			dalekscoord[i][2] = 1
			board[dalekrow][dalekcol] = "*"
			dalekcount = dalekcount - 1
			score = score + 10
		}
		i = i + 1
	}
}

func movedaleks() { // Move unless already dead. ([3]==1)
	var i = 0
	for i < daleks {
		if dalekscoord[i][2] == 0 { // Not dead yet. OK to move
			dalekrow = dalekscoord[i][0] // Save current Dalek Row (y) Coord
			dalekcol = dalekscoord[i][1] // Save current Dalek Column (x) Coord
			// Test if we need to move the Dalek row closer to Doctor Who
			if dalekrow > whorow { // Decrement Dalek Row
				board[dalekrow][dalekcol] = "."
				dalekscoord[i][0] = dalekrow - 1
			} else if dalekrow < whorow { // Increment Dalek Row
				board[dalekrow][dalekcol] = "."
				dalekscoord[i][0] = dalekrow + 1
			} else { // dalekrow must be = whorow, don't change Dalek Row
				dalekscoord[i][0] = dalekrow
			}

			if dalekcol > whocol { // Decrement Dalek Column
				board[dalekrow][dalekcol] = "."
				dalekscoord[i][1] = dalekcol - 1
			} else if dalekcol < whocol { // Increment Dalek Column
				board[dalekrow][dalekcol] = "."
				dalekscoord[i][1] = dalekcol + 1
			} else { // dalekcol must be = whocol, don't change col
				dalekscoord[i][1] = dalekcol
			}

			if dalekonwho() == 1 { // Dalek landed on Doctor Who. Game Over
				lives = lives - 1
				if lives == 0 {
					fmt.Println("\nPoor Doctor Who :(\n")
					os.Exit(0)
				}
			}

			if dalekondalek(i) == 1 {
				dalekcount = countdaleks()
				dalekscoord[i][2] = 1                             // Still not Dead
				board[dalekscoord[i][0]][dalekscoord[i][1]] = "*" // Set the new board position for the Dalek
				score = score + 40
				// drawboard()
				// fmt.Println("\nDaleks Crashed!\n")
				// os.Exit(0)
			} else {
				dalekscoord[i][2] = 0                             // Still not Dead
				board[dalekscoord[i][0]][dalekscoord[i][1]] = "D" // Set the new board position for the Dalek
			}

		}
		i = i + 1
	}
}

//******************************************************************************************************************

func main() {

	// Initialize Board
	initboard()
	initdaleks()
	initwho()

	err := term.Init()
	if err != nil {
		panic(err)
	}

	defer term.Close()

	fmt.Println("Welcome to Daleks Game")

	//Print initial game board
	drawboard()

keyPressListenerLoop:
	for {
		switch ev := term.PollEvent(); ev.Type {
		case term.EventKey:
			switch ev.Key {
			case term.KeyEsc:
				break keyPressListenerLoop
			case term.KeyArrowUp:
				reset()
				fmt.Println("Arrow Up pressed")
				board[whorow][whocol] = "."
				if whorow == 0 {
					whorow = (numrows - 1)
				} else {
					whorow = whorow - 1
				}
				board[whorow][whocol] = "W"
				movedaleks() // Move Daleks closer to Doctor Who
				drawboard()

			case term.KeyArrowDown:
				reset()
				fmt.Println("Arrow Down pressed")
				board[whorow][whocol] = "."
				if whorow == (numrows - 1) {
					whorow = 0
				} else {
					whorow = whorow + 1
				}
				board[whorow][whocol] = "W"
				movedaleks() // Move Daleks closer to Doctor Who
				drawboard()

			case term.KeyArrowLeft:
				reset()
				fmt.Println("Arrow Left pressed")
				board[whorow][whocol] = "."
				if whocol == 0 {
					whocol = (numcols - 1)
				} else {
					whocol = whocol - 1
				}
				board[whorow][whocol] = "W"
				movedaleks() // Move Daleks closer to Doctor Who
				drawboard()

			case term.KeyArrowRight:
				reset()
				fmt.Println("Arrow Right pressed")
				board[whorow][whocol] = "."
				if whocol == (numcols - 1) {
					whocol = 0
				} else {
					whocol = whocol + 1
				}
				board[whorow][whocol] = "W"
				movedaleks() // Move Daleks closer to Doctor Who
				drawboard()

			default:
				reset()
				switch ev.Ch {
				case 46: // "." Stand
					fmt.Println("Stand")
					board[whorow][whocol] = "W"
					movedaleks() // Move Daleks closer to Doctor Who
					drawboard()

				case 115: // "s" Sonic Screwdriver
					fmt.Println("Sonic Screwdriver")
					if screwdrivers > 0 {
						blastdaleks()
						screwdrivers = screwdrivers - 1
					}
					movedaleks() // Move Daleks closer to Doctor Who
					drawboard()

				case 116: // "t" Teleport
					if teleports > 0 {
						fmt.Println("Teleport")
						board[whorow][whocol] = "."
						seednum := rand.NewSource(time.Now().UnixNano())
						rnum := rand.New(seednum)
						whorow = rnum.Intn(numrows)
						whocol = rnum.Intn(numcols)
						board[whorow][whocol] = "W"
						teleports = teleports - 1
						movedaleks() // Move Daleks closer to Doctor Who
						drawboard()
					} else {
						fmt.Println("(No Teleports Left)")
						movedaleks() // Move Daleks closer to Doctor Who
						drawboard()
					}

				default:
					//fmt.Println("Unknown Key")
				}
			}
		case term.EventError:
			panic(ev.Err)
		}
	}
}
