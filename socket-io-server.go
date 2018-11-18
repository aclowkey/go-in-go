package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	log "github.com/cloudflare/cfssl/log"
	"github.com/googollee/go-socket.io"
)

type SocketIOServer struct {
	port int
}

func (server *SocketIOServer) Start() {
	socketServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	gameSessions := make(map[string]bool)
	socketServer.On("connection", func(so socketio.Socket) {
		var gameID string
		playerRole := "white"
		gameIDparams, ok := so.Request().URL.Query()["gameID"]
		if !ok && len(gameID) != 1 {
			gameID = strconv.Itoa(r1.Int())
			gameSessions[gameID] = true
			log.Infof("Game %s created \n", gameID)
		} else {
			gameID = gameIDparams[0]
			_, found := gameSessions[gameID]
			if !found {
				log.Errorf("Player tried to joined room which doesn't exist!")
				return
			}
			playerRole = "black"
		}
		playerID := strconv.Itoa(r1.Int())
		so.Join(gameID)
		log.Infof("Player %s (%s) join game (%s)", playerRole, playerID, gameID)
		// log.Infof("Request: %s\n", )
		// so.Emit("joined_game", fmt.Sprintf("Joined game %s", gameID))
	})

	http.Handle("/", socketServer)
	log.Infof("Listening for socket-io on :%d\n", server.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil))
}
