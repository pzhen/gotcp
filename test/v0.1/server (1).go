package main

import "gotcp/Gotcp"

func main() {
	srv := Gotcp.InitServer("Test V0.1", "127.0.0.1", 9888, "tcp")
	srv.Run()
}
