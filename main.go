package main

func main() {
	server := MakeSocketIOServer(9070)
	server.Start()
}
