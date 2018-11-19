package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"

	log "github.com/cloudflare/cfssl/log"
	"github.com/googollee/go-socket.io"
)

type SocketIOServer struct {
	port         int
	gameSessions map[string]*IOGameSession
}

type Player struct {
	name   string
	piece  Piece
	socket socketio.Socket
}

type IOGameSession struct {
	game    *Game
	player1 *Player
	player2 *Player
}

func (gameSession *IOGameSession) join(player *Player) (hasPlace bool) {
	hasPlace = true
	if gameSession.player1 == nil {
		gameSession.player1 = player
	} else if gameSession.player2 == nil {
		gameSession.player2 = player
	} else {
		hasPlace = false
	}
	return hasPlace
}

func (gameSession *IOGameSession) ready() bool {
	return gameSession.player1 != nil &&
		gameSession.player2 != nil
}

func (gameSession *IOGameSession) boardChanged(gameID string) {
	gameSession.player1.socket.Emit("board_changed", gameID)
	gameSession.player2.socket.Emit("board_changed", gameID)
}

func MakeSocketIOServer(port int) *SocketIOServer {
	gameSessions := make(map[string]*IOGameSession)
	return &SocketIOServer{
		port,
		gameSessions,
	}
}

func (server *SocketIOServer) Start() {
	socketServer, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Level = log.LevelDebug

	socketServer.On("connection", server.handleConnection)

	http.Handle("/", socketServer)
	http.HandleFunc("/game/", server.handleGame())
	log.Infof("Listening for socket-io on :%d\n", server.port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", server.port), nil))
}

func (server *SocketIOServer) handleGame() http.HandlerFunc {
	pattern, _ := regexp.Compile("^\\/game\\/(?P<GameID>[^\\/]+)?$")
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			matches := pattern.FindStringSubmatch(r.URL.String())
			if len(matches) != 2 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			gameID := matches[1]
			if gameID == "" {
				keys := make([]string, len(server.gameSessions))
				i := 0
				for k := range server.gameSessions {
					keys[i] = k
					i++
				}
				bytes, err := json.Marshal(keys)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(bytes)
			} else {
				session, ok := server.gameSessions[gameID]
				if !ok {
					w.WriteHeader(http.StatusNotFound)
					return
				}
				game := session.game
				bytes, err := json.Marshal(game)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.Write(bytes)
			}
		case http.MethodPut:
			r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
			gameID := strconv.Itoa(r1.Int())
			server.gameSessions[gameID] = &IOGameSession{
				CreateGame(9, 4.5), nil, nil,
			}
			type GameCreated struct {
				GameID string `json:"gameId"`
			}
			bytes, _ := json.Marshal(&GameCreated{gameID})
			w.Write(bytes)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}

func (server *SocketIOServer) handleConnection(so socketio.Socket) {
	// Handle user connecting to server by either joining a game
	var player *Player
	gameIDParams, ok := so.Request().URL.Query()["gameID"]
	if !ok || len(gameIDParams) != 1 {
		so.Emit("error", "Invalid request. Must provide gameID parameter")
	}
	gameID := gameIDParams[0]
	gameSession, found := server.gameSessions[gameID]
	if !found {
		log.Debugf("User attempted to join an invalid game: %s\n", gameID)
		so.Emit("error", "Cannot find game")
		return
	}
	player = &Player{"Player 2", Black, so}
	hasPlace := gameSession.join(player)
	if !hasPlace {
		log.Debugf("User attempted to join an already full game: %s\n", gameID)
		so.Emit("error", "Room already full")
		return
	}

	log.Debugf("[%s] Player joined game", gameID)

	// Waiting for second player to join so we'll be ready
	// Maybe better approach is the server will notify itself on game_started?
	for !gameSession.ready() {
		// TODO timeout
	}

	// Game is ready, handle movement logic
	so.Emit("game_started", gameSession.game.Board.Pieces())
	so.On("move", server.handleMove(gameID, gameSession, player))
}

func (server *SocketIOServer) handleMove(gameID string, gameSession *IOGameSession, player *Player) interface{} {
	game := gameSession.game
	return func(data string) {
		var position Position
		err := json.Unmarshal([]byte(data), &position)
		if err != nil {
			player.socket.Emit("error", "Invalid syntax!")
			return
		}

		result, err := game.Move(&Move{position.X, position.Y, player.piece})
		if result == Ok {
			gameSession.boardChanged(gameID)
			player.socket.Emit("message", fmt.Sprintf("Thank you mr, %s", player.name))
			log.Debugf(game.Board.String(false))
			log.Debugf("[%s] Player %s moved %s\n", gameID, player.piece)
		} else {
			player.socket.Emit("error", err.Error())
		}
	}
}
