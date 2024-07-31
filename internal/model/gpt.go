package model

type ModelType string

var (
	KimiType     = ModelType("kimi")
	DeepSeekType = ModelType("deepSeek")
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatGptResp struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

type RequestModel struct {
	Messages       []*Message        `json:"messages,omitempty"`
	Model          string            `json:"model,omitempty"`
	Stream         bool              `json:"stream"`
	Temperature    float64           `json:"temperature"`
	ResponseFormat map[string]string `json:"response_format"`
	MaxTokens      int               `json:"max_tokens,omitempty"`
}

type GPT struct {
	ChatGPTPlatform ModelType `mapstructure:"chat_gpt_platform" json:"chat_gpt_platform" yaml:"chat_gpt_platform"`
	Kimi            Kimi      `mapstructure:"kimi" json:"kimi" yaml:"kimi"`
	DeepSeek        Kimi      `mapstructure:"deep_seek" json:"deep_seek" yaml:"deep_seek"`
	Temperature     float64   `mapstructure:"temperature" json:"temperature" yaml:"temperature"`
	Prompt          string    `mapstructure:"prompt" json:"prompt" yaml:"prompt"`
}

type Kimi struct {
	Model          string            `mapstructure:"model" json:"model" yaml:"model"`
	ApiKey         string            `mapstructure:"api_key" json:"api_key" yaml:"api_key"`
	BaseURL        string            `mapstructure:"base_url" json:"base_url" yaml:"base_url"`
	ResponseFormat map[string]string `json:"response_format"`
}
