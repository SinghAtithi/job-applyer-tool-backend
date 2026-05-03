package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"example.com/internal/config"
	"example.com/pkg/logger"
)

type ConfigAgentClient struct {
	config     *config.AgentClient
	httpClient *http.Client
}

type ChatRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type ResponseFormat struct {
	Type   string                 `json:"type"`
	Schema map[string]interface{} `json:"schema,omitempty"`
}

type ChatResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

type Choice struct {
	Message MessageContent `json:"message"`
}

type MessageContent struct {
	Content interface{} `json:"content"`
}

type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

var (
	sharedClient *ConfigAgentClient
	clientOnce   sync.Once
	clientErr    error
)

func GetSharedClient() (*ConfigAgentClient, error) {
	clientOnce.Do(func() {
		cfg, err := config.LoadAgentClient()
		if err != nil {
			clientErr = fmt.Errorf("failed to load agent client config: %w", err)
			return
		}
		sharedClient = NewAgentClient(cfg)
	})
	if clientErr != nil {
		return nil, clientErr
	}
	return sharedClient, nil
}

func NewAgentClient(cfg *config.AgentClient) *ConfigAgentClient {
	return &ConfigAgentClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

const maxRetries = 3

func isTransient(statusCode int) bool {
	return statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusInternalServerError ||
		statusCode == http.StatusBadGateway ||
		statusCode == http.StatusServiceUnavailable
}

func (c *ConfigAgentClient) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = c.config.Model
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			delay := time.Duration(1<<uint(attempt-1)) * time.Second
			logger.Debug("retrying chat completion (attempt %d/%d) after %v", attempt+1, maxRetries, delay)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, body, err := c.doRequest(ctx, jsonData)
		if err != nil {
			lastErr = err
			continue
		}

		if isTransient(resp.StatusCode) {
			lastErr = fmt.Errorf("transient API error: %d %s", resp.StatusCode, string(body))
			logger.Debug("transient API error (attempt %d/%d): %d", attempt+1, maxRetries, resp.StatusCode)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			var errResp ChatResponse
			if unmarshalErr := json.Unmarshal(body, &errResp); unmarshalErr == nil && errResp.Error != nil {
				return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error.Message)
			}
			return nil, fmt.Errorf("HTTP error: %d %s — body: %s", resp.StatusCode, resp.Status, truncate(body, 200))
		}

		var chatResp ChatResponse
		if err := json.Unmarshal(body, &chatResp); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return &chatResp, nil
	}

	return nil, fmt.Errorf("chat completion failed after %d retries: %w", maxRetries, lastErr)
}

func (c *ConfigAgentClient) doRequest(ctx context.Context, jsonData []byte) (*http.Response, []byte, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	c.setHeaders(httpReq)

	logger.Debug("sending chat completion request to %s (model: %s)", c.config.BaseURL, c.config.Model)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return resp, body, nil
}

func truncate(b []byte, n int) string {
	if len(b) > n {
		return string(b[:n]) + "..."
	}
	return string(b)
}

func (c *ConfigAgentClient) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Title", "Job-Applyer-Tool")
}

func logHTTPDetails(req *http.Request, resp *http.Response) {
	logger.Debug("=== HTTP REQUEST ===")
	logger.Debug("Method: %s", req.Method)
	logger.Debug("URL: %s", req.URL.String())

	for k, v := range req.Header {
		if k == "Authorization" {
			logger.Debug("  Header %s: [REDACTED]", k)
		} else {
			logger.Debug("  Header %s: %v", k, v)
		}
	}

	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(body))
		logger.Debug("Body: %s", body)
	}

	logger.Debug("=== HTTP RESPONSE ===")
	logger.Debug("Status: %s", resp.Status)

	for k, v := range resp.Header {
		logger.Debug("  Header %s: %v", k, v)
	}

	body, _ := io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewReader(body))
	logger.Debug("Body: %s", body)
}