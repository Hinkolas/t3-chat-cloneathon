package chat

// Request types
type Request struct {
	Model               string    `json:"model"`
	Temperature         float64   `json:"temperature,omitempty"`
	MaxCompletionTokens int       `json:"max_completion_tokens,omitempty"`
	TopP                float64   `json:"top_p,omitempty"`
	Stream              bool      `json:"stream"`
	ReasoningEffort     int32     `json:"reasoning_effort,omitempty"`
	Stop                any       `json:"stop,omitempty"`
	Messages            []Message `json:"messages"`
}

type Message struct {
	Role      string `json:"role"`
	Content   string `json:"content"`
	Reasoning string `json:"reasoning"`
}

type Options struct {
	AnthropicAPIKey string `db:"anthropic_api_key" json:"anthropic_api_key"`
	OpenAIAPIKey    string `db:"openai_api_key" json:"openai_api_key"`
	GeminiAPIKey    string `db:"gemini_api_key" json:"gemini_api_key"`
	OllamaBaseURL   string `db:"ollama_base_url" json:"ollama_base_url"`
	SystemPrompt    string `json:"system,omitempty"` // System prompt for the chat.
}
