package main

import (
	"errors"
	"fmt"
)

type Game struct {
	board Board
	turn  Piece
	komi  float32
}

type MoveResult int

const (
	Ok = iota
	Illegal
	Komi
	GameOver
)

type GameResult int

const (
	WhiteWins = iota
	BlackWins
	Draw
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

func (game *Game) getMove() (move *Move, err error) {
	var x int
	var y int
	fmt.Printf(game.board.String(false))
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
		if result != Ok {
			fmt.Printf("Illegal move. Try again!\n")
			continue
		}
	}
}
