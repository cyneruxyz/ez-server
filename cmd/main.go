package main

import (
	"ex-server/internal/handler"
	"ex-server/internal/server"
	"log"
)

const serverPort = 8080

func main() {
	log.Fatal(server.Init(serverPort, handler.Init()).Run())
}