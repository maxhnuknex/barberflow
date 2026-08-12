package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const baseURL = "https://api.groq.com/openai/v1"

type Client struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

func NewClient(
	apiKey string,
	model string,
	httpClient *http.Client,
) *Client {
	return &Client{
		apiKey:     apiKey,
		model:      model,
		httpClient: httpClient,
	}
}

type chatRequest struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func (c *Client) Complete(
	ctx context.Context,
	text string,
) (string, error) {

	body := chatRequest{
		Model: c.model,
		Messages: []message{
			{
				Role:    "user",
				Content: text,
			},
		},
	}

	data, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal groq request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		baseURL+"/chat/completions",
		bytes.NewReader(data),
	)
	if err != nil {
		return "", fmt.Errorf("create groq request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send groq request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("groq returned status: %s", resp.Status)
	}

	var result chatResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode groq response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("groq returned no choices")
	}

	return result.Choices[0].Message.Content, nil
}
