package main

import (
	"flag"
	"fmt"
	"grep-server/internal/delivery"
	"grep-server/internal/service"
	"log"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "port to listen on")
	flag.Parse()

	srv := delivery.NewServer(service.NewService())
	for err := srv.Start(fmt.Sprintf(":%d", port)); err != nil && port < 65535; port++ {
		log.Printf("failed to start server on port %d: %v", port, err)
	}

	log.Printf("server started on port %d", port)
}
