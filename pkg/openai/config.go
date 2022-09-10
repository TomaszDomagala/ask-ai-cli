package openai

// RequestBase will be used in the request body.
type RequestBase struct {
	Model            string  `json:"model" name:"model" usage:"OpenAI model to use" value:"code-davinci-002"`
	Temperature      float64 `json:"temperature" name:"temperature" usage:"temperature" value:"0"`
	MaxTokens        int     `json:"max_tokens" name:"max-tokens" usage:"max tokens" value:"256"`
	TopP             float64 `json:"top_p" name:"top-p" usage:"top p" value:"1"`
	FrequencyPenalty float64 `json:"frequency_penalty" name:"frequency-penalty" usage:"frequency penalty" value:"0"`
	PresencePenalty  float64 `json:"presence_penalty" name:"presence-penalty" usage:"presence penalty" value:"0"`
}

type Config struct {
	// ApiKey is the OpenAI API key.
	ApiKey string `name:"apikey" usage:"OpenAI API key"`

	// OpenAI request configuration
	RequestBase `name:",squash"`
}
