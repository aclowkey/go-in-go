package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Utility consts
var cellOffsets = [][]int{
	{-1, 0},
	{1, 0},
	{0, 1},
	{0, -1},
}

// Piece deinfes whether it's black, white or empty
type Piece int

const (
	// Empty is a place which isn't occupied
	Empty Piece = iota
	// White is the first player
	White
	// Black is the second player
	Black
)

// Cell has a piece which occupies it, and a number of liberties available to it
type Cell struct {
	piece   Piece
	liberty int
}

func (piece Piece) String() string {
	names := [...]string{
		" ",
		"O",
		"X",
	}
	if piece < Empty || piece > Black {
		return "?"
	}
	return names[piece]
}

func (cell Cell) String(printLiberty bool) string {
	if cell.piece < Empty || cell.piece > Black {
		return "Unknown"
	}
	liberty := " "
	if printLiberty {
		liberty = strconv.Itoa(cell.liberty)
	}
	return fmt.Sprintf("%s %s   ", cell.piece.String(), liberty)

}

// Grid is a matrix of cells
type Grid [][]Cell

// Board is responsible for containing the cells, and history
type Board struct {
	size            int
	data            Grid
	moves           int
	boardHistory    BoardQueue
	movementHistroy MovementQueue
}

// MakeBoard constructs a board of size size*size
func MakeBoard(size int) *Board {
	data := make(Grid, size)
	for y := range data {
		data[y] = make([]Cell, size)
		for x := range data[y] {
			// Adjust available liberties initially
			liberties := 4
			if x == 0 || x == size-1 {
				liberties--
			}
			if y == 0 || y == size-1 {
				liberties--
			}
			// Create the cell
			data[y][x] = Cell{
				Empty,
				liberties,
			}
		}
	}
	return &Board{
		size,
		data,
		0,
		*MakeBoardQueue(3),
		*MakeMovementQueue(3),
	}
}

// Move defines where a player placed a piece in form of x, y
type Move struct {
	x     int
	y     int
	piece Piece
}

func (move *Move) String() string {
	return fmt.Sprintf("%s to (%d, %d)", move.piece.String(), move.x, move.y)
}

// Move contains the logic of validating the move and changing the board in accordance
func (board *Board) Move(move *Move) (err error) {
	// First check: Is the cell empty?
	if board.data[move.x][move.y].piece != Empty {
		err = fmt.Errorf("cell (%d, %d) is occupied", move.x, move.y)
		return
	}
	placePiece := true
	if board.data[move.x][move.y].liberty == 0 {
		// Maybe it kills something and allows for liberty
		placePiece = false
	}
	// Updating liberty of neighbours
	for i := range cellOffsets {
		newX, newY := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(newX, newY) {
			continue
		}
		cell := &board.data[newX][newY]
		cell.liberty--
	}
	// Checking neighbour for either a kill, additional liberty
	for i := range cellOffsets {
		newX, newY := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(newX, newY) {
			continue
		}
		cell := &board.data[newX][newY]
		// If neighbour is not empty, enemy, and has no liberty
		// check if it's connected to a piece with liberty
		if cell.piece != move.piece && cell.piece != Empty && cell.liberty == 0 {
			if board.KillConfirm(nil, Move{newX, newY, cell.piece}) {
				placePiece = true
				board.Kill(Move{newX, newY, cell.piece})
			}
		}

		if !placePiece && cell.piece == move.piece {
			if cell.liberty > 0 {
				placePiece = true
			} else {
				noLiberty := board.KillConfirm(nil, *move)
				if !noLiberty {
					placePiece = true
				}
			}
		}
	}
	if placePiece {
		// Place the piece
		board.data[move.x][move.y].piece = move.piece
		board.moves++
		board.movementHistroy.Enqueue(move)
		var historialBoard = make(Grid, board.size)
		copy(historialBoard, board.data)
		board.boardHistory.Enqueue(&historialBoard)
	} else {
		// Piece couldn't be placed so give the neighbours their liberties
		for i := range cellOffsets {
			newX, newY := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
			if !board.Inbounds(newX, newY) {
				continue
			}
			cell := &board.data[newX][newY]
			cell.liberty++
		}
		err = fmt.Errorf("cell (%d, %d) has no liberty", move.x, move.y)
		return
	}
	return nil
}

// KillConfirm checks if the piece at the move doesn't have any liberty connected to it
func (board *Board) KillConfirm(visited [][]bool, move Move) bool {
	// Initilizing visit array
	if visited == nil {
		visited := make([][]bool, board.size)
		for x := range visited {
			visited[x] = make([]bool, board.size)
			for y := range visited[x] {
				visited[x][y] = false
			}
		}
		return board.KillConfirm(visited, move)
	}
	// Look for neighbouring allies
	for i := range cellOffsets {
		newX, newY := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(newX, newY) || board.data[newX][newY].piece != move.piece {
			continue
		}
		cell := board.data[newX][newY]
		if cell.piece != move.piece {
			continue
		}
		if visited[newX][newY] {
			continue
		}
		visited[newX][newY] = true
		// Piece at newX, new_ y is an allie, does it have liberty?
		if cell.liberty > 0 {
			return false // No kill!
		}
		// Maybe it has an allie with liberty
		if !board.KillConfirm(visited, Move{newX, newY, move.piece}) {
			return false // No Kill!
		}

	}
	return true // Piece has no more liberty!
}

// Kill is emptying the piece/s connected to the move
func (board *Board) Kill(move Move) {
	board.data[move.x][move.y].piece = Empty
	for i := range cellOffsets {
		newX, newY := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(newX, newY) {
			continue
		}
		cell := &board.data[newX][newY]
		cell.liberty++
	}

	for i := range cellOffsets {
		newX, newY := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(newX, newY) || board.data[newX][newY].piece != move.piece {
			continue
		}
		board.Kill(Move{newX, newY, move.piece})
	}
}

func (board *Board) PrintHistory(printLiberty bool) {
	for i := 0; i < 3; i++ {
		fmt.Printf(board.movementHistroy.data[i].String() + "\n")
		fmt.Printf(board.String(false, i))
	}
}

func (board *Board) String(printLiberty bool, history int) string {
	var str strings.Builder
	str.WriteString("=========         Move #" + strconv.Itoa(history) + "    ===================\n")
	for x := 0; x < board.size; x++ {
		str.WriteString("  " + strconv.Itoa(x) + "   ")
	}
	var grid = *board.boardHistory.data[history]
	str.WriteString(PrintGrid(printLiberty, &grid))
	str.WriteString("====================================================\n")
	return str.String()
}

func PrintGrid(printLiberty bool, grid *Grid) string {
	var str strings.Builder
	str.WriteString("\n----------------------------------------------------\n")
	size := len(*grid)
	for y := 0; y < size; y++ {
		str.WriteString(strconv.Itoa(y) + "|")
		for x := 0; x < size; x++ {
			str.WriteString((*grid)[x][y].String(printLiberty))
		}
		str.WriteString("\n")
	}
	str.WriteString("r-\n")
	return str.String()
}

func (board *Board) Inbounds(x int, y int) bool {
	return x >= 0 && x < board.size && y >= 0 && y < board.size
}

func (board *Board) SafeMove(x int, y int, piece Piece) {
	err := board.Move(&Move{x, y, piece})
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
