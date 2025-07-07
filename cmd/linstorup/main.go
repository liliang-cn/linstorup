package main

import (
	"flag"
	"log"

	"github.com/liliang-cn/linstorup/internal/server"
)

func main() {
	port := flag.Int("port", 3374, "Port to run the web server on")
	flag.Parse()

	srv, err := server.NewServer(*port)
	if err != nil {
		log.Fatalf("Error creating server: %v", err)
	}
	srv.Start()
}
