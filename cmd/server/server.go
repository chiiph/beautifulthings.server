package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"beautifulthings/server"
	"beautifulthings/store"
)

func main() {
	const addr = ":8080"

	log.Println("Starting server at", addr)

	var st store.ObjectStore

	log.Println("Store selected: %s", os.Getenv("STORE"))

	if os.Getenv("STORE") == "memory" {
		st = store.NewInMemoryServer()
	} else {
		st = store.NewGCS()
	}

	cancel, err := server.ServeRest(context.Background(), addr, st)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	log.Println("Ctrl+C...")
	cancel()
}
