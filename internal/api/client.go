package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles communication with the OpenRouter API
type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new OpenRouter API client
func NewClient(baseURL, apiKey string, timeout int) *Client {
	httpClient := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	return &Client{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		HTTPClient: httpClient,
	}
}

// SendChatCompletion sends a chat completion request to the API
func (c *Client) SendChatCompletion(req *ChatCompletionRequest) (*ChatCompletionResponse, error) {
	url := fmt.Sprintf("%s/chat/completions", c.BaseURL)

	// Marshal request to JSON
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("HTTP-Referer", "https://github.com/kdevrou/openrouter-cli")
	httpReq.Header.Set("X-Title", "OpenRouter CLI")

	// Send request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	// Unmarshal response
	var chatResp ChatCompletionResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &chatResp, nil
}

// ListModels fetches the list of available models
func (c *Client) ListModels() ([]Model, error) {
	url := fmt.Sprintf("%s/models", c.BaseURL)

	// Create HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIKey))
	req.Header.Set("HTTP-Referer", "https://github.com/kdevrou/openrouter-cli")
	req.Header.Set("X-Title", "OpenRouter CLI")

	// Send request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode >= 400 {
		return nil, parseAPIError(resp.StatusCode, respBody)
	}

	// Unmarshal response
	var modelsResp ModelsResponse
	if err := json.Unmarshal(respBody, &modelsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return modelsResp.Data, nil
}

// parseAPIError parses an error response from the API
func parseAPIError(statusCode int, body []byte) error {
	var errorResp map[string]interface{}
	if err := json.Unmarshal(body, &errorResp); err != nil {
		// If we can't parse as JSON, return a generic error
		return &APIError{
			StatusCode: statusCode,
			Message:    fmt.Sprintf("HTTP %d: %s", statusCode, string(body)),
		}
	}

	// Try to extract error message
	var message string
	if err, ok := errorResp["error"]; ok {
		if errMap, ok := err.(map[string]interface{}); ok {
			if msg, ok := errMap["message"]; ok {
				message = fmt.Sprintf("%v", msg)
			}
			if errType, ok := errMap["type"]; ok {
				return &APIError{
					StatusCode: statusCode,
					Message:    message,
					Type:       fmt.Sprintf("%v", errType),
				}
			}
		} else if msg, ok := err.(string); ok {
			message = msg
		}
	}

	if message == "" {
		message = "Unknown error"
	}

	return &APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}
