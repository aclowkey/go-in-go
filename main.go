package main

import (
	"errors"
	"fmt"
	"net"

	log "github.com/cloudflare/cfssl/log"
)

const port = 9060

type Position struct {
	x int
	y int
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
	defer player1.Close()
	for {
		player1.Write([]byte("0, Your move\n"))
		player2.Write([]byte("0, Wait for your turn\n"))
		player1Position, err := getMove(player1)
		log.Debugf("Got position: %+v", *player1Position)
		if err != nil {
			player1.Write([]byte(fmt.Sprintf("1, Invalid move: %s\n", err)))
			continue
		}
		player1.Write([]byte("0, Wait for your turn"))
		player2.Write([]byte("0, Your move\n"))
		player2Position, err := getMove(player2)
		if err != nil {
			player1.Write([]byte(fmt.Sprintf("1, Invalid move: %s\n", err)))
			continue
		}
		log.Debugf("Got position: %+v", *player2Position)

	}
}

func main() {
	log.Level = log.LevelDebug
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	defer ln.Close()
	if err != nil {
		log.Fatalf("Failed to start server: %+v\n", err)
		return
	}
	log.Infof("Go-in-go is ready: listening at %d\n", port)
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
