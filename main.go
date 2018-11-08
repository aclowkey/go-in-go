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

func handleConnection(conn net.Conn) {
	log.Debugf("Accepted a connection from %s\n", conn.RemoteAddr())
	defer conn.Close()
	for {
		position, err := getMove(conn)
		if err != nil {
			conn.Write([]byte(fmt.Sprintf("1, Invalid move: %s\n", err)))
			continue
		}
		conn.Write([]byte("0\n"))
		log.Debugf("Got position: %+v", *position)
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
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalf("Connection issue: %+v\n", err)
		}
		go handleConnection(conn)
	}
}
