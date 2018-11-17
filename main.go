package main

func main() {
	server := MakeHTTPServer(9060)
	server.start()
}
