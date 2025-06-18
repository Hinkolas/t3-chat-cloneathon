package chat

// Request types
type Request struct {
	Model               string     `json:"model"`
	Temperature         float64    `json:"temperature"`
	MaxCompletionTokens int        `json:"max_completion_tokens"`
	TopP                float64    `json:"top_p"`
	Stream              bool       `json:"stream"`
	ReasoningEffort     int32      `json:"reasoning_effort"`
	Stop                any        `json:"stop,omitempty"`
	Messages            []*Message `json:"messages"`
	System              string     `json:"system"` // System prompt for the chat session
}

type Message struct {
	Role        string        `json:"role"`
	Content     string        `json:"content"`
	Reasoning   string        `json:"reasoning"`
	Attachments []*Attachment `json:"attachments"` // Image attachments as byte slices
}

type Attachment struct {
	MimeType string `json:"mime_type"` // MIME type of attachment
	Data     []byte `json:"data"`
}

type Options map[string]string
