package api

// Message represents a chat message
type Message struct {
	Role    string `json:"role"` // "user", "assistant", "system"
	Content string `json:"content"`
}

// ChatCompletionRequest is the request payload for chat completions
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
}

// Choice represents a completion choice in the response
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage contains token usage statistics
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionResponse is the response from a chat completion request
type ChatCompletionResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// ModelPricing contains pricing information for a model
type ModelPricing struct {
	Prompt     string `json:"prompt"`
	Completion string `json:"completion"`
}

// Architecture contains architectural information about a model
type Architecture struct {
	Modality     string `json:"modality"`
	Tokenizer    string `json:"tokenizer"`
	InstructType string `json:"instruct_type,omitempty"`
}

// Model represents an available LLM model
type Model struct {
	ID            string       `json:"id"`
	Name          string       `json:"name"`
	Created       int64        `json:"created"`
	ContextLength int          `json:"context_length"`
	Pricing       ModelPricing `json:"pricing"`
	Architecture  Architecture `json:"architecture"`
	Description   string       `json:"description,omitempty"`
	TopProvider   interface{}  `json:"top_provider,omitempty"` // Can be string or object
}

// ModelsResponse is the response from the models list endpoint
type ModelsResponse struct {
	Data []Model `json:"data"`
}

// APIError represents an error from the OpenRouter API
type APIError struct {
	StatusCode int
	Message    string
	Type       string
}

func (e *APIError) Error() string {
	if e.Type != "" {
		return e.Type + ": " + e.Message
	}
	return e.Message
}
