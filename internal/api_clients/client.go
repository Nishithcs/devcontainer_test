package api_clients

import (
	"bytes"
	"clusterix-code/internal/auth"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	BaseURL    string
	httpClient *http.Client
	token      auth.TokenProvider
}

type ClientOption func(c *Client)

func NewClient(baseURL string, options ...ClientOption) *Client {
	c := &Client{
		BaseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range options {
		opt(c)
	}

	return c
}

func (c *Client) SetTokenProvider(provider auth.TokenProvider) {
	c.token = provider
}

func (c *Client) Request(ctx context.Context, method, path string, body interface{}, response interface{}, headers map[string]string) error {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	var req *http.Request
	var err error

	switch v := body.(type) {
	case *bytes.Buffer:
		req, err = http.NewRequestWithContext(ctx, method, url, v)
		if err != nil {
			return fmt.Errorf("failed to create request with multipart form data: %w", err)
		}
	default:
		var bodyReader io.Reader
		if body != nil {
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return fmt.Errorf("failed to marshal request body: %w", err)
			}
			bodyReader = bytes.NewReader(jsonBody)
			if headers != nil {
				if _, exists := headers["Content-Type"]; !exists {
					headers["Content-Type"] = "application/json"
				}
			}
		}
		req, err = http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Set authorization header if token provider exists
	if c.token != nil {
		token, err := c.token.GetToken(ctx)
		if err != nil {
			return fmt.Errorf("failed to get token: %w", err)
		}
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if response != nil {
		switch v := response.(type) {
		case *[]byte:
			data, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read binary response: %w", err)
			}
			*v = data
		default:
			if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
		}
	}

	return nil
}
