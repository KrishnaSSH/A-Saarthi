package handler

import (
	"net/http"

	"github.com/KrishnaSSH/a-saarthi/server"
)

var mux = server.New()

// Handler is the Vercel serverless entrypoint; vercel.json rewrites /api/* here.
func Handler(w http.ResponseWriter, r *http.Request) { mux.ServeHTTP(w, r) }
