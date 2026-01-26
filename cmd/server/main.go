package main

import (
	"log"
	"os"

	"wrzapi/internal/server"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := server.New()
	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := srv.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}