package main

import (
	"fmt"
	"net"

	log "github.com/cloudflare/cfssl/log"
)

const port = 9060

func handleConnection(conn net.Conn) {
	log.Debugf("Accepted a connection from %s\n", conn.RemoteAddr())
	defer conn.Close()
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		log.Errorf("Cannot read ")
	}
	log.Debugf("Got message: %s\n", string(data[:n]))
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
