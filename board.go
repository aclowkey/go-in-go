package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Utility consts
var cellOffsets [][]int = [][]int{
	{-1, 0},
	{1, 0},
	{0, 1},
	{0, -1},
}

// Structs
type Piece int

const (
	Empty Piece = iota
	White
	Black
)

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
	} else {
		return names[piece]
	}
}

func (cell Cell) String(printLiberty bool) string {
	if cell.piece < Empty || cell.piece > Black {
		return "Unknown"
	} else {
		liberty := " "
		if printLiberty {
			liberty = strconv.Itoa(cell.liberty)
		}
		return fmt.Sprintf("%s %s   ", cell.piece.String(), liberty)
	}
}

type Board struct {
	size            int
	data            [][]Cell
	moves           int
	boardHistory    BoardQueue
	movementHistroy MovementQueue
}

func MakeBoard(size int) *Board {
	data := make([][]Cell, size)
	for y := range data {
		data[y] = make([]Cell, size)
		for x := range data[y] {
			// Adjust available liberties initially
			liberties := 4
			if x == 0 || x == size-1 {
				liberties -= 1
			}
			if y == 0 || y == size-1 {
				liberties -= 1
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

type Move struct {
	x     int
	y     int
	piece Piece
}

func (move *Move) String() string {
	return fmt.Sprintf("%s to (%d, %d)", move.piece.String(), move.x, move.y)
}

func (board *Board) Move(move *Move) (err error) {
	// First check: Is the cell empty?
	if board.data[move.x][move.y].piece != Empty {
		err = fmt.Errorf("Cell (%d, %d) is occupied!", move.x, move.y)
		return
	}
	placePiece := true
	if board.data[move.x][move.y].liberty == 0 {
		// Maybe it kills something and allows for liberty
		placePiece = false
	}
	// Updating liberty of neighbours
	for i := range cellOffsets {
		new_x, new_y := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(new_x, new_y) {
			continue
		}
		cell := &board.data[new_x][new_y]
		cell.liberty -= 1
	}

	for i := range cellOffsets {
		new_x, new_y := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(new_x, new_y) {
			continue
		}
		cell := &board.data[new_x][new_y]
		// If neighbour is not empty, enemy, and has no liberty
		// check if it's connected to a piece with liberty
		if cell.piece != move.piece && cell.piece != Empty && cell.liberty == 0 {
			if board.KillConfirm(nil, Move{new_x, new_y, cell.piece}) {
				placePiece = true
				board.Kill(Move{new_x, new_y, cell.piece})
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
		board.moves += 1
		board.movementHistroy.Enqueue(move)
		var historialBoard [][]Cell = make([][]Cell, board.size)
		copy(historialBoard, board.data)
		board.boardHistory.Enqueue(&historialBoard)
	} else {
		// Piece couldn't be placed so give the neighbours their liberties
		for i := range cellOffsets {
			new_x, new_y := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
			if !board.Inbounds(new_x, new_y) {
				continue
			}
			cell := &board.data[new_x][new_y]
			cell.liberty += 1
		}
		err = fmt.Errorf("Cell (%d, %d) has no liberty!", move.x, move.y)
		return
	}
	return nil
}

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
		new_x, new_y := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(new_x, new_y) || board.data[new_x][new_y].piece != move.piece {
			continue
		}
		cell := board.data[new_x][new_y]
		if cell.piece != move.piece {
			continue
		}
		if visited[new_x][new_y] {
			continue
		}
		visited[new_x][new_y] = true
		// Piece at new_x, new_ y is an allie, does it have liberty?
		if cell.liberty > 0 {
			return false // No kill!
		} else {
			// Maybe it has an allie with liberty
			if !board.KillConfirm(visited, Move{new_x, new_y, move.piece}) {
				return false // No Kill!
			}
		}
	}
	return true // Piece has no more liberty!
}

func (board *Board) Kill(move Move) {
	board.data[move.x][move.y].piece = Empty
	for i := range cellOffsets {
		new_x, new_y := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(new_x, new_y) {
			continue
		}
		cell := &board.data[new_x][new_y]
		cell.liberty += 1
	}

	for i := range cellOffsets {
		new_x, new_y := move.x+cellOffsets[i][0], move.y+cellOffsets[i][1]
		if !board.Inbounds(new_x, new_y) || board.data[new_x][new_y].piece != move.piece {
			continue
		}
		board.Kill(Move{new_x, new_y, move.piece})
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
	var cells [][]Cell = *board.boardHistory.data[history]
	str.WriteString(PrintCells(board.size, printLiberty, cells))
	str.WriteString("====================================================\n")
	return str.String()
}

func PrintCells(size int, printLiberty bool, cells [][]Cell) string {
	var str strings.Builder
	str.WriteString("\n----------------------------------------------------\n")
	for y := 0; y < size; y++ {
		str.WriteString(strconv.Itoa(y) + "|")
		for x := 0; x < size; x++ {
			str.WriteString(cells[x][y].String(printLiberty))
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
