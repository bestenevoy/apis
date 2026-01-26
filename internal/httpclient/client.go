package httpclient

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"
)

const maxBytes = 2 << 20 // 2MB

var ErrFetch = errors.New("failed to fetch page")

type Client struct {
	hc *http.Client
}

func New() *Client {
	return &Client{hc: &http.Client{Timeout: 12 * time.Second}}
}

func (c *Client) FetchHTML(ctx context.Context, rawURL string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, rawURL, ErrFetch
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Notelook/1.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, rawURL, ErrFetch
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, rawURL, ErrFetch
	}

	limited := io.LimitReader(resp.Body, maxBytes)
	body, err := io.ReadAll(limited)
	if err != nil {
		return nil, rawURL, ErrFetch
	}

	return body, resp.Request.URL.String(), nil
}
