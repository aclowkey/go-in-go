package main

import (
	"errors"
	"fmt"
	"net"

	log "github.com/cloudflare/cfssl/log"
)

type TcpServer struct {
	port int
}

func (server *TcpServer) start() {
	log.Level = log.LevelDebug
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", server.port))
	defer ln.Close()
	if err != nil {
		return
	}
	log.Infof("Go-in-go is ready: listening at %d\n", server.port)
	for {
		player1, err := ln.Accept()
		if err != nil {
			log.Fatalf("Connection issue: %+v\n", err)
		}
		player1.Write([]byte("Welcome player! Waiting for partner\n"))
		log.Infof("Player 1 joined. Waiting for player 2..")
		player2, err := ln.Accept()
		if err != nil {
			log.Fatalf("Connection issue: %+v\n", err)
		}
		handleConnection(player1, player2)

	}
}

func parsePosition(data string) (*Position, error) {
	var x int
	var y int
	_, err := fmt.Sscanf(data, "%d %d", &x, &y)
	if err != nil {
		return nil, errors.New("Invalid move: should be: x y")
	}
	return &Position{x, y}, nil
}

func getMove(conn net.Conn) (*Position, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		log.Errorf("Cannot read %+v\n", err)
	}
	data := string(buffer[:n])
	var position *Position
	position, err = parsePosition(data)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Should be 'x y' got: %s", data))
	}
	log.Debugf("Parsed position: %+v", *position)
	return position, nil
}

func handleConnection(player1 net.Conn, player2 net.Conn) {
	log.Debugf("%+v Vs %+v", player1.RemoteAddr(), player2.RemoteAddr())
	game := CreateGame(9, 4.5)
	defer player1.Close()
	defer player2.Close()
	var currentPlayer *net.Conn
	var otherPlayer *net.Conn
	for {
		// Assigning the proper player
		if game.Turn == White {
			currentPlayer = &player1
			otherPlayer = &player2
		} else {
			currentPlayer = &player2
			otherPlayer = &player1
		}

		// Displaying board
		player1.Write([]byte(game.Board.String(false)))
		player2.Write([]byte(game.Board.String(false)))

		// Instructing players
		(*currentPlayer).Write([]byte("0, " + game.Turn.String() + "'s turn\n"))
		(*otherPlayer).Write([]byte("0, Wait for your turn\n"))
		move, err := getMove(*currentPlayer)
		if err != nil {
			(*currentPlayer).Write([]byte(fmt.Sprintf("1, Invalid move: %s\n", err)))
			continue
		}
		log.Debugf("%s tried to move to %+v", game.Turn.String(), *move)
		result, err := game.Move(&Move{(*move).X, (*move).Y, game.Turn})
		if result != Ok {
			(*currentPlayer).Write([]byte(fmt.Sprintf("1, Invalid move: %v\n", err.Error())))
		}
	}
}
