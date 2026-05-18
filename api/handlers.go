package api

import (
	"encoding/json"
	"net/http"

	"pc-dashboard/system"
)

func setCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

func InfoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setCacheHeaders(w)

	json.NewEncoder(w).Encode(system.SystemInfoData)
}

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	setCacheHeaders(w)

	json.NewEncoder(w).Encode(system.SystemStatsData)
}

func RegisterRoutes() {
	http.HandleFunc("/api/info", InfoHandler)
	http.HandleFunc("/api/stats", StatsHandler)

	fs := http.FileServer(http.Dir("./static"))

	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})
}
