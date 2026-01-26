package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func Time(w http.ResponseWriter, r *http.Request) {
	now := time.Now().UTC()
	resp := map[string]interface{}{
		"utc":   now.Format(time.RFC3339),
		"unix":  now.Unix(),
		"local": time.Now().Format(time.RFC3339),
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	_ = enc.Encode(v)
}