package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
)

const (
	// apiUrl is the URL of the OpenAI API.
	apiUrl = "https://api.openai.com/v1/completions"

	// queryPrefixSequence is a sequence of tokens that is used to prefix the query.
	// it is also used as stop sequence to terminate the completion.
	queryPrefixSequence = "command-query"
)

var (
	errNoSuggestionFound = fmt.Errorf("no suggestion found")
)

//// CompletionConfig represents a base configuration for a completionRequest
//type CompletionConfig struct {
//	Model            string  `json:"model"`
//	Temperature      float64 `json:"temperature"`
//	MaxTokens        int     `json:"max_tokens"`
//	TopP             float64 `json:"top_p"`
//	FrequencyPenalty float64 `json:"frequency_penalty"`
//	PresencePenalty  float64 `json:"presence_penalty"`
//}

//// DefaultCompletionConfig is a default configuration for completion.
//var DefaultCompletionConfig = CompletionConfig{
//	Model:            "code-davinci-002",
//	Temperature:      0,
//	MaxTokens:        256,
//	TopP:             1,
//	FrequencyPenalty: 0,
//	PresencePenalty:  0,
//}

type Client struct {
	Config Config
}

func NewClient(config Config) *Client {
	return &Client{
		Config: config,
	}
}

type requestBody struct {
	RequestBase
	Prompt string   `json:"prompt"`
	Stop   []string `json:"stop"`
}

// responseBody is a body of completion request response
type responseBody struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Text         string      `json:"text"`
		Index        int         `json:"index"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

// errorResponse is a body of error response
type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

/*
prompt creates a prompt for a completion request.
Example of a prompt with a query "show current directory":

	command-query: create foo directory:
	mkdir foo
	command-query: show current directory:
*/
func prompt(query string) string {
	return fmt.Sprintf("%s\n%s\n%s\n", promptQuery("create foo directory"), "mkdir foo", promptQuery(query))
}

// promptQuery creates the query part of a prompt.
func promptQuery(query string) string {
	return fmt.Sprintf("%s: %s:", queryPrefixSequence, query)
}

func (c *Client) Suggest(query string) (string, error) {

	reqBody := requestBody{
		RequestBase: c.Config.RequestBase,
		Prompt:      prompt(query),
		Stop:        []string{queryPrefixSequence},
	}
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiUrl, bytes.NewReader(jsonReqBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Config.ApiKey))
	req.Header.Add("Content-Type", "application/json")

	for k, v := range req.Header {
		log.Debug().Strs(k, v).Msg("request header")
	}
	log.Debug().Str("body", string(jsonReqBody)).Msg("request body")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()
	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	log.Debug().Msgf("response body: %s", resBody)

	if res.StatusCode != http.StatusOK {
		var errRes errorResponse
		if err := json.Unmarshal(resBody, &errRes); err != nil {
			return "", fmt.Errorf("failed to unmarshal error response: %w", err)
		}

		return "", fmt.Errorf("unexpected status code: %d, body: %s", res.StatusCode, resBody)
	}

	var completion responseBody
	if err := json.Unmarshal(resBody, &completion); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no completion found")
	}

	return strings.TrimSpace(completion.Choices[0].Text), nil
}
