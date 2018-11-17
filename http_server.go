package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HTTPServer is a Go-in-go server in HTTP
type HTTPServer struct {
	port  int
	games map[int]*Game
}

// MakeHTTPServer creates the server
func MakeHTTPServer(port int) *HTTPServer {
	return &HTTPServer{
		port,
		make(map[int]*Game),
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
		server.games[gameID] = CreateGame(9, 4.5)
		sessionID := rand.Intn(10000)
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
		game, ok := server.games[gameID]
		if !ok {
			c.JSON(404, gin.H{
				"message": fmt.Sprintf("Game %d not found", gameID),
			})
			return
		}
		c.JSON(200, gin.H{
			"board": game.board.data,
			"turn":  game.turn.String(),
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
		game, ok := server.games[gameID]
		if !ok {
			c.JSON(404, gin.H{
				"message": fmt.Sprintf("Game %d not found", gameID),
			})
			return
		}
		sessionID := rand.Intn(10000)
		c.Header("sessionID", strconv.Itoa(sessionID))
		c.JSON(200, gin.H{
			"board": game.board.data,
			"turn":  game.turn.String(),
		})
	})
	r.Run(fmt.Sprintf(":%d", server.port))
}
