package main

import (
	"errors"
	"fmt"
)

type Game struct {
	Board Board   `json:"board"`
	Turn  Piece   `json:"turn"`
	Komi  float32 `json:"komi"`
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
	fmt.Printf(game.Board.String(false))
	fmt.Printf("%s's turn: ", game.Turn)
	_, err = fmt.Scanf("%d %d", &x, &y)
	if err != nil {
		return nil, errors.New("Invalid move: should be: x, y")
	}
	return &Move{x, y, game.Turn}, nil

}

func (game *Game) Move(move *Move) (MoveResult, error) {
	if move.piece != game.Turn {
		return Illegal, errors.New("Not your turn")
	}
	err := game.Board.Move(move)
	if err != nil {
		// TODO komi r
		return Illegal, err
	}
	if game.Turn == White {
		game.Turn = Black
	} else {
		game.Turn = White
	}
	return Ok, nil

}
func (game *Game) Start() {
	gameOver := false
	for !gameOver {
		move, err := game.getMove()
		if err != nil {
			fmt.Printf("Invalid move: %s", err.Error())
			continue
		}
		result, err := game.Move(move)
		if result != Ok {
			fmt.Printf("Illegal move. Try again!\n")
			continue
		}
	}
}
