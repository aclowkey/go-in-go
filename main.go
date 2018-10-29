package main

import (
	"errors"
	"fmt"
	"os"
)

type Game struct {
	board Board
	turn  Piece
	komi  float32
}

type MoveResult int

const (
	OK = iota
	ILLEGAL
	KOMI
	GAME_OVER
)

type GameResult int

const (
	WHITE_WINS = iota
	BLACK_WINS
	DRAW
)

func CreateGame(size int, komi float32) *Game {
	board := MakeBoard(size)
	game := &Game{
		*board,
		White,
		komi,
	}
	return game
}

func (game *Game) Move(move *Move) MoveResult {
	if move.piece != game.turn {
		fmt.Println("Not your turn!")
		return ILLEGAL
	}
	err := game.board.Move(move)
	if err != nil {
		// TODO komi r
		os.Exit(1)
		return ILLEGAL
	} else {
		if game.turn == White {
			game.turn = Black
		} else {
			game.turn = White
		}
		return OK
	}
}
func (game *Game) getMove() (move *Move, err error) {
	var x int
	var y int
	fmt.Printf(game.board.String(false, 0))
	fmt.Printf("%s's turn: ", game.turn)
	_, err = fmt.Scanf("%d %d", &x, &y)
	if err != nil {
		return nil, errors.New("Invalid move: should be: x, y")
	}
	return &Move{x, y, game.turn}, nil

}

func (game *Game) Start() {
	gameOver := false
	for !gameOver {
		move, err := game.getMove()
		if err != nil {
			fmt.Printf("Invalid move: %s", err.Error())
			continue
		}
		result := game.Move(move)
		if result != OK {
			fmt.Printf("Illegal move. Try again!\n")
			continue
		}

	}
}

func main() {
	game := CreateGame(9, 4.5)
	game.board.Move(&Move{3, 3, White})
	game.board.Move(&Move{2, 2, White})
	game.board.Move(&Move{1, 3, White})

	game.board.Move(&Move{1, 4, Black})
	// game.board.Move(&Move{2, 5, Black})
	// game.board.Move(&Move{3, 4, Black})

	// // Ko
	// game.board.Move(&Move{2, 3, Black}) // Move number 1
	// game.board.Move(&Move{2, 4, White}) // Move number 2
	// game.board.Move(&Move{2, 3, Black}) // Move number 3
	// A Ko happens if history[0] history[2]
	game.board.PrintHistory(true)
	// game.Start()

}
