package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	// apiUrl is the URL of the OpenAI API.
	apiUrl = "https://api.openai.com/v1/completions"

	// queryPrefixSequence is a sequence of tokens that is used to prefix the query.
	// it is also used as stop sequence to terminate the completion.
	queryPrefixSequence = "query"
)

type Client struct {
	Config Config
}

// NewClient creates a new OpenAI client.
func NewClient(config Config) *Client {
	return &Client{
		Config: config,
	}
}

// Suggest suggests a command for a given query.
func (c *Client) Suggest(query string) (string, error) {
	prompt := suggestPrompt(query)
	return c.doRequest(prompt)
}

// Explain explains a command.
func (c *Client) Explain(command string) (string, error) {
	prompt := explainPrompt(command)
	return c.doRequest(prompt)
}

// doRequest performs a request to the OpenAI API.
func (c *Client) doRequest(prompt string) (string, error) {
	reqBody := requestBody{
		RequestBase: c.Config.RequestBase,
		Prompt:      prompt,
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
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close response body")
		}
	}(res.Body)

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}
	var completion responseBody
	if err = json.Unmarshal(resBody, &completion); err != nil {
		return "", fmt.Errorf("failed to unmarshal response (status code: %d): %w", res.StatusCode, err)
	}

	log.Debug().Msgf("response body: %s", resBody)
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d, body: %s", res.StatusCode, resBody)
	}

	if len(completion.Choices) == 0 {
		return "", fmt.Errorf("no completion found")
	}

	return strings.TrimSpace(completion.Choices[0].Text), nil
}

// requestBody is the request body of a completion request.
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
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param"`
		Code    string `json:"code"`
	} `json:"error"`
}

/*
suggestPrompt creates a prompt for a completion request.
Example of a prompt with a query "show current directory":

	query: create foo directory:
	answer: mkdir foo
	query: show current directory:
	answer:
*/
func suggestPrompt(query string) string {
	var builder strings.Builder
	// Prompt example:
	builder.WriteString("query: create foo directory\n")
	builder.WriteString("answer: mkdir foo\n")
	// Actual prompt:
	builder.WriteString("query: ")
	builder.WriteString(query)
	builder.WriteString("\n")
	builder.WriteString("answer: ")

	return builder.String()
}

// explainPrompt creates a prompt for an explanation request.
func explainPrompt(command string) string {
	var builder strings.Builder
	// Prompt example:
	builder.WriteString("query: cd $HOME\n")
	builder.WriteString("answer:\nChange the current directory to the home directory\n")
	//	Actual command:
	builder.WriteString("query: ")
	builder.WriteString(command)
	builder.WriteString("\n")
	builder.WriteString("answer:\n")

	return builder.String()
}
