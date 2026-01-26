package main

import (
	"flag"
	"log"
	"os"

	"wrzapi/internal/server"
)

func main() {
	var serverURL string
	var port string
	flag.StringVar(&serverURL, "server-url", "", "OpenAPI server URL (overrides SERVER_URL env)")
	flag.StringVar(&port, "port", "", "HTTP listen port (overrides PORT env)")
	flag.Parse()

	if serverURL != "" {
		_ = os.Setenv("SERVER_URL", serverURL)
	}

	if port == "" {
		port = os.Getenv("PORT")
	}
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
