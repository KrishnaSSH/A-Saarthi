package main

import (
	"log"
	"net/http"
	"os"

	"github.com/KrishnaSSH/a-saarthi/internal/server"
)

func main() {
	mux := server.New()
	// Serve the built frontend when present (local/self-hosted; Vercel serves it from its CDN).
	if _, err := os.Stat("frontend/dist"); err == nil {
		mux.Handle("/", http.FileServer(http.Dir("frontend/dist")))
	}
	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
