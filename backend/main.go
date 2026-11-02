package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func newUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	b[6], b[8] = b[6]&0x0f|0x40, b[8]&0x3f|0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func lang(r *http.Request) string {
	if l := r.URL.Query().Get("lang"); l != "" {
		return l
	}
	return "en"
}

func main() {
	abdm := newABDM()
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]any{"ok": true, "abdmConfigured": abdm.configured()})
	})

	// Full key-value table of rejection reasons.
	mux.HandleFunc("GET /api/reasons", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, allReasons(lang(r)))
	})

	mux.HandleFunc("GET /api/reasons/{code}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, lookup(lang(r), r.PathValue("code")))
	})

	// Automatic location detection: browser sends coordinates, we return a structured address.
	mux.HandleFunc("GET /api/location", func(w http.ResponseWriter, r *http.Request) {
		lat, e1 := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		lon, e2 := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
		if e1 != nil || e2 != nil || lat < 6 || lat > 38 || lon < 68 || lon > 98 {
			writeJSON(w, 400, lookup(lang(r), "LOCATION_UNVERIFIED"))
			return
		}
		a, err := reverseGeocode(r.Context(), lat, lon)
		if err != nil {
			writeJSON(w, 502, lookup(lang(r), "LOCATION_UNVERIFIED"))
			return
		}
		writeJSON(w, 200, a)
	})

	// Pre-submission check: tells the user every problem up front, in plain words.
	mux.HandleFunc("POST /api/check", func(w http.ResponseWriter, r *http.Request) {
		var a Application
		if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
			writeJSON(w, 400, map[string]string{"error": "invalid JSON"})
			return
		}
		issues := []Reason{}
		for _, c := range check(a) {
			issues = append(issues, lookup(lang(r), c))
		}
		writeJSON(w, 200, map[string]any{"ok": len(issues) == 0, "issues": issues})
	})

	// ABHA lookup through the official ABDM gateway; failures are translated to readable reasons.
	mux.HandleFunc("POST /api/abha/search", func(w http.ResponseWriter, r *http.Request) {
		var in struct {
			ABHA string `json:"abha"`
		}
		if json.NewDecoder(r.Body).Decode(&in) != nil || len(in.ABHA) < 14 {
			writeJSON(w, 400, lookup(lang(r), "ABHA_NOT_FOUND"))
			return
		}
		status, body, err := abdm.Do(r.Context(), "POST", "/abha/api/v3/profile/account/abha/search", map[string]string{"abhaNumber": in.ABHA})
		if err == errNotConfigured {
			writeJSON(w, 503, map[string]any{"configured": false, "message": "Set ABDM_CLIENT_ID and ABDM_CLIENT_SECRET to talk to the official servers."})
			return
		}
		if err != nil || status >= 400 {
			writeJSON(w, 502, lookup(lang(r), mapUpstream(status, body)))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	})

	// Serve the built frontend if present.
	if _, err := os.Stat("../frontend/dist"); err == nil {
		mux.Handle("/", http.FileServer(http.Dir("../frontend/dist")))
	}

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
