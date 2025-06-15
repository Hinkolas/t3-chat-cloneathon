package llm

import "fmt"

type ModelProvider string

const (
	OpenAI    ModelProvider = "openai"
	Anthropic ModelProvider = "anthropic"
	Gemini    ModelProvider = "gemini"
	Ollama    ModelProvider = "ollama"
)

func (p *ModelProvider) UnmarshalText(text []byte) error {
	s := string(text)
	switch s {
	case "openai":
		*p = OpenAI
	case "anthropic":
		*p = Anthropic
	case "ollama":
		*p = Ollama
	case "gemini":
		*p = Gemini
	default:
		return fmt.Errorf("unknown provider %q", s)
	}
	return nil
}

type ModelFeatures struct {
	HasVision          bool `json:"has_vision,omitempty"`
	HasPDF             bool `json:"has_pdf,omitempty"`
	HasReasoning       bool `json:"has_reasoning,omitempty"`
	HasEffortControl   bool `json:"has_effort_control,omitempty"`
	HasWebSearch       bool `json:"has_web_search,omitempty"`
	HasFast            bool `json:"has_fast,omitempty"`
	HasImageGeneration bool `json:"has_image_generation,omitempty"`
}
type ModelFlags struct {
	IsPremium      bool `json:"is_premium,omitempty"`
	IsExperimental bool `json:"is_experimental,omitempty"`
	IsKeyRequired  bool `json:"is_key_required,omitempty"`
	IsFree         bool `json:"is_free,omitempty"`
	IsNew          bool `json:"is_new,omitempty"`
	IsRecommended  bool `json:"is_recommended,omitempty"`
	IsOpenSource   bool `json:"is_open_source,omitempty"`
}
type Model struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	Name        string        `json:"name"`
	Provider    ModelProvider `json:"provider"`
	Features    ModelFeatures `json:"features"`
	Flags       ModelFlags    `json:"flags"`
}
