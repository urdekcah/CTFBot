package ctftime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
	Limiter    *rate.Limiter
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 100,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		UserAgent: "CTFTimeClient/1.0",
		Limiter:   rate.NewLimiter(rate.Every(time.Second), 10),
	}
}

func (c *Client) request(ctx context.Context, method, path string, body io.Reader, target interface{}) error {
	if err := c.Limiter.Wait(ctx); err != nil {
		return fmt.Errorf("rate limit exceeded: %w", err)
	}

	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	if target != nil {
		decoder := json.NewDecoder(resp.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(target); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
	}

	return nil
}

func (c *Client) Get(ctx context.Context, path string, target interface{}) error {
	return c.request(ctx, http.MethodGet, path, nil, target)
}

func (c *Client) Post(ctx context.Context, path string, body io.Reader, target interface{}) error {
	return c.request(ctx, http.MethodPost, path, body, target)
}

func (c *Client) Put(ctx context.Context, path string, body io.Reader, target interface{}) error {
	return c.request(ctx, http.MethodPut, path, body, target)
}

func (c *Client) Delete(ctx context.Context, path string, target interface{}) error {
	return c.request(ctx, http.MethodDelete, path, nil, target)
}
