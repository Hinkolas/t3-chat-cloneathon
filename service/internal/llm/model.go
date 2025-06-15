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
	HasVision          bool `json:"has_vision,omitempty" mapstructure:"has_vision"`
	HasPDF             bool `json:"has_pdf,omitempty" mapstructure:"has_pdf"`
	HasReasoning       bool `json:"has_reasoning,omitempty" mapstructure:"has_reasoning"`
	HasEffortControl   bool `json:"has_effort_control,omitempty" mapstructure:"has_effort_control"`
	HasWebSearch       bool `json:"has_web_search,omitempty" mapstructure:"has_search"`
	HasFast            bool `json:"has_fast,omitempty" mapstructure:"has_fast"`
	HasImageGeneration bool `json:"has_image_generation,omitempty" mapstructure:"has_image_generation"`
}

type ModelFlags struct {
	IsPremium      bool `json:"is_premium,omitempty" mapstructure:"is_premium"`
	IsExperimental bool `json:"is_experimental,omitempty" mapstructure:"is_experimental"`
	IsKeyRequired  bool `json:"is_key_required,omitempty" mapstructure:"is_key_required"`
	IsFree         bool `json:"is_free,omitempty" mapstructure:"is_free"`
	IsNew          bool `json:"is_new,omitempty" mapstructure:"is_new"`
	IsRecommended  bool `json:"is_recommended,omitempty" mapstructure:"is_recommended"`
	IsOpenSource   bool `json:"is_open_source,omitempty" mapstructure:"is_open_source"`
}

type Model struct {
	Title       string        `json:"title" mapstructure:"title"`
	Description string        `json:"description" mapstructure:"description"`
	Icon        string        `json:"icon" mapstructure:"icon"`
	Name        string        `json:"name" mapstructure:"name"`
	Provider    ModelProvider `json:"provider" mapstructure:"provider"`
	Features    ModelFeatures `json:"features" mapstructure:"features"`
	Flags       ModelFlags    `json:"flags" mapstructure:"flags"`
}
