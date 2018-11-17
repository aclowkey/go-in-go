package main

func main() {
	server := &TcpServer{9600}
	server.start()
}
