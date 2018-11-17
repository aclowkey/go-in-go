package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GameSession has the game, and the players
type GameSession struct {
	game      *Game
	player1id *int
	player2id *int
	public    bool
}

func (session *GameSession) isReady() bool {
	return session.player2id != nil
}

// Position defines a place on the board
type Position struct {
	X int `json:"x" binding:"required"`
	Y int `json:"y" binding:"required"`
}

// HTTPServer is a Go-in-go server in HTTP
type HTTPServer struct {
	port  int
	games map[int]*GameSession
}

// MakeHTTPServer creates the server
func MakeHTTPServer(port int) *HTTPServer {
	return &HTTPServer{
		port,
		make(map[int]*GameSession),
	}
}

func (server *HTTPServer) start() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"health": "OK",
		})
	})
	r.POST("/game", func(c *gin.Context) {
		gameID := 5
		game := CreateGame(9, 4.5)
		sessionID := rand.Intn(10000)
		gameSession := &GameSession{
			game,
			&sessionID,
			nil,
			false,
		}
		server.games[gameID] = gameSession
		c.Header("gameID", strconv.Itoa(gameID))
		c.Header("sessionID", strconv.Itoa(sessionID))
		c.JSON(200, gin.H{
			"gameID": gameID,
		})
	})
	r.GET("/game/:id", func(c *gin.Context) {

		gameIDParam := c.Param("id")
		gameID, err := strconv.Atoi(gameIDParam)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid game ID",
			})
			return
		}
		gameSession, ok := server.games[gameID]
		if !ok {
			c.JSON(404, gin.H{
				"message": fmt.Sprintf("Game %d not found", gameID),
			})
			return
		}
		if !gameSession.public {
			sessionIDHeader := c.GetHeader("sessionID")
			sessionID, err := strconv.Atoi(sessionIDHeader)
			if err != nil {
				// Either expired, or not authorized, or invalid
				c.JSON(404, gin.H{
					"message": "Not allowed to access the game",
				})
				return
			}
			if !(sessionID == *gameSession.player1id || sessionID == *gameSession.player2id) {
				c.JSON(404, gin.H{
					"message": "Not allowed to access the game",
				})
				return
			}
		}

		c.JSON(200, gin.H{
			"turn": gameSession.game.turn,
		})
	})
	r.POST("/game/:id", func(c *gin.Context) {
		gameIDParam := c.Param("id")
		gameID, err := strconv.Atoi(gameIDParam)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid game ID",
			})
			return
		}
		gameSession, ok := server.games[gameID]
		if !ok {
			c.JSON(404, gin.H{
				"message": fmt.Sprintf("Game %d not found", gameID),
			})
			return
		}
		if gameSession.isReady() {
			c.JSON(400, gin.H{
				"message": "Game has already started",
			})
			return
		}
		sessionID := rand.Intn(10000)
		gameSession.player2id = &sessionID
		c.Header("sessionID", strconv.Itoa(sessionID))
		c.JSON(200, gin.H{
			"turn": gameSession.game.turn.String(),
		})
	})

	r.POST("/game/:id/move", func(c *gin.Context) {
		gameIDParam := c.Param("id")
		gameID, err := strconv.Atoi(gameIDParam)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid game ID",
			})
			return
		}
		gameSession, ok := server.games[gameID]
		if !ok {
			c.JSON(404, gin.H{
				"message": fmt.Sprintf("Game %d not found", gameID),
			})
			return
		}

		sessionIDHeader := c.GetHeader("sessionID")
		sessionID, err := strconv.Atoi(sessionIDHeader)
		if err != nil {
			// Either expired, or not authorized, or invalid
			c.JSON(404, gin.H{
				"message": "Not allowed to access the game",
			})
			return
		}
		if gameSession.game.turn == White {
			if sessionID != *gameSession.player1id {
				c.JSON(400, gin.H{
					"message": "It's white player turn",
				})
				return
			}
		} else {
			if sessionID != *gameSession.player2id {
				c.JSON(400, gin.H{
					"message": "It's black player turn",
				})
				return
			}
		}
		position := Position{}
		err = c.ShouldBindJSON(&position)
		if err != nil {
			c.JSON(400, gin.H{
				"message": "Invalid request: should have 'x' and 'y'",
			})
			return
		}
		move := &Move{position.X, position.Y, gameSession.game.turn}
		result := gameSession.game.Move(move)
		if result == Illegal {
			c.JSON(400, gin.H{
				"message": "Invalid move",
			})
			return
		}
		fmt.Fprintln(gin.DefaultWriter, gameSession.game.board.String(false))
	})

	r.Run(fmt.Sprintf(":%d", server.port))
}
