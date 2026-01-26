package handlers

import (
	"net/http"
	"net/url"

	"wrzapi/internal/httpclient"
	"wrzapi/internal/pageinfo"
)

func PageInfo(w http.ResponseWriter, r *http.Request) {
	raw := r.URL.Query().Get("url")
	if raw == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing url"})
		return
	}

	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid url"})
		return
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "url must be http or https"})
		return
	}

	client := httpclient.New()
	body, finalURL, err := client.FetchHTML(r.Context(), parsed.String())
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]string{"error": err.Error()})
		return
	}

	meta := pageinfo.ParseHTML(body, finalURL)
	writeJSON(w, http.StatusOK, meta)
}