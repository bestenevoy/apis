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
	var navData string
	var navDev bool
	flag.StringVar(&serverURL, "server-url", "", "OpenAPI server URL (overrides SERVER_URL env)")
	flag.StringVar(&port, "port", "", "HTTP listen port (overrides PORT env)")
	flag.StringVar(&navData, "nav-data", "", "Nav data file path (overrides NAV_DATA env)")
	flag.BoolVar(&navDev, "nav-dev", false, "Serve nav frontend from disk for hot reload")
	flag.Parse()

	if serverURL != "" {
		_ = os.Setenv("SERVER_URL", serverURL)
	}

	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}

	if navData == "" {
		navData = os.Getenv("NAV_DATA")
		if navData == "" {
			navData = "data.json"
		}
	}

	srv, err := server.New(server.Config{
		NavDataPath: navData,
		NavDev:      navDev,
	})
	if err != nil {
		log.Fatal(err)
	}
	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := srv.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}
