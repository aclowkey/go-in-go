package main

func main() {
	game := CreateGame(9, 4.5)
	// game.board.Move(&Move{1, 2, White})
	// game.board.Move(&Move{1, 1, Black})

	// game.board.Move(&Move{0, 3, White})
	// game.board.Move(&Move{0, 0, Black})

	// game.board.Move(&Move{0, 1, White}) // This bored
	// game.board.Move(&Move{0, 2, Black})

	// game.board.Move(&Move{0, 1, White}) // This move is a Ko!
	// game.board.PrintHistory(false)
	// game.board.Move(&Move{0, 1, Black})

	// A Ko happens if history[0] history[2]
	// game.board.PrintHistory(true)
	game.Start()

}
