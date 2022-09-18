package openai

// RequestBase will be used in the request body.
type RequestBase struct {
	Model            string  `json:"model" config:"openai.model"`
	Temperature      float64 `json:"temperature" config:"openai.temperature"`
	MaxTokens        int     `json:"max_tokens" config:"openai.maxtokens"`
	TopP             float64 `json:"top_p" config:"openai.topp"`
	FrequencyPenalty float64 `json:"frequency_penalty" config:"openai.frequencypenalty"`
	PresencePenalty  float64 `json:"presence_penalty" config:"openai.presencepenalty"`
}

type Config struct {
	// ApiKey is the OpenAI API key.
	ApiKey string

	// OpenAI request configuration
	RequestBase
}
