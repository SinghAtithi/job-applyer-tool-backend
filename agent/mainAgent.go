package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"example.com/internal/config"
	"example.com/pkg/logger"
)

// ConfigAgentClient wraps the AI API client with configuration
type ConfigAgentClient struct {
	config     *config.AgentClient
	httpClient *http.Client
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model          string          `json:"model"`
	Messages       []Message       `json:"messages"`
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"`
}

// Message represents a single message in the conversation
type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ResponseFormat specifies the expected response format
type ResponseFormat struct {
	Type   string                 `json:"type"`
	Schema map[string]interface{} `json:"schema,omitempty"`
}

// ChatResponse represents the API response
type ChatResponse struct {
	Choices []Choice  `json:"choices"`
	Error   *APIError `json:"error,omitempty"`
}

// Choice represents a single choice in the response
type Choice struct {
	Message MessageContent `json:"message"`
}

// MessageContent holds the response content
type MessageContent struct {
	Content interface{} `json:"content"`
}

// APIError represents an error returned by the API
type APIError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Code    string `json:"code"`
}

// NewAgentClient creates a new AI agent client
func NewAgentClient(cfg *config.AgentClient) *ConfigAgentClient {
	return &ConfigAgentClient{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ChatCompletion sends a chat completion request and returns the response
func (c *ConfigAgentClient) ChatCompletion(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	if req.Model == "" {
		req.Model = c.config.Model
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.config.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	c.setHeaders(httpReq)

	logger.Debug("sending chat completion request to %s (model: %s)", c.config.BaseURL, req.Model)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if chatResp.Error != nil {
			return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, chatResp.Error.Message)
		}
		return nil, fmt.Errorf("HTTP error: %d %s", resp.StatusCode, resp.Status)
	}

	return &chatResp, nil
}

func (c *ConfigAgentClient) setHeaders(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Title", "Job-Applyer-Tool")
}

// logHTTPDetails logs full HTTP request/response details (only in debug mode)
func logHTTPDetails(req *http.Request, resp *http.Response) {
	logger.Debug("=== HTTP REQUEST ===")
	logger.Debug("Method: %s", req.Method)
	logger.Debug("URL: %s", req.URL.String())

	for k, v := range req.Header {
		logger.Debug("  Header %s: %v", k, v)
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
