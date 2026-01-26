package handlers

import (
	_ "embed"
	"net/http"
	"os"
	"strings"
)

//go:embed openapi.yaml
var openapiYAML []byte

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>WRZ API Docs</title>
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <style>
    body { margin: 0; background: #f5f6f7; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      SwaggerUIBundle({
        url: '/openapi.yaml',
        dom_id: '#swagger-ui',
        presets: [SwaggerUIBundle.presets.apis],
        layout: 'BaseLayout'
      });
    };
  </script>
</body>
</html>
`

func OpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(applyServerURL(openapiYAML, r))
}

func Docs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerHTML))
}

func applyServerURL(src []byte, r *http.Request) []byte {
	serverURL := strings.TrimSpace(os.Getenv("SERVER_URL"))
	if serverURL == "" {
		serverURL = requestBaseURL(r)
	}
	return []byte(strings.ReplaceAll(string(src), "{{SERVERS_BLOCK}}", buildServersBlock(serverURL, r)))
}

func requestBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	if forwarded := r.Header.Get("X-Forwarded-Proto"); forwarded != "" {
		scheme = strings.Split(forwarded, ",")[0]
		scheme = strings.TrimSpace(scheme)
	}
	host := r.Host
	if forwardedHost := r.Header.Get("X-Forwarded-Host"); forwardedHost != "" {
		host = strings.Split(forwardedHost, ",")[0]
		host = strings.TrimSpace(host)
	}
	return scheme + "://" + host
}

func buildServersBlock(serverURL string, r *http.Request) string {
	requestURL := requestBaseURL(r)
	var urls []string
	if serverURL != "" {
		urls = append(urls, serverURL)
	}
	if requestURL != "" && requestURL != serverURL {
		urls = append(urls, requestURL)
	}
	if len(urls) == 0 {
		return "  - url: http://localhost:8080"
	}
	lines := make([]string, 0, len(urls))
	for _, u := range urls {
		lines = append(lines, "  - url: "+u)
	}
	return strings.Join(lines, "\n")
}
