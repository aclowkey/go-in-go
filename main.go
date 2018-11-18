package main

func main() {
	server := &SocketIOServer{9070}
	server.Start()
}
