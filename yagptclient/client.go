package yagptclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type MessageT struct {
	Role string `json:"role"`
	Text string `json:"text"`
} 

type GenerationResponse struct {
	Result struct {
		Alternatives []struct{
			Message MessageT`json:"message"`
		} `json:"alternatives"`
	} `json:"result"`
}

type GenerationRequest struct {
	ModelURI string `json:"modelUri"`
	CompletionOptions struct {
		Stream bool `json:"stream"`
		Temperature float32 `json:"temperature"`
		MaxTokens string `json:"maxTokens"`
	} `json:"completionOptions"`
	Messages []MessageT `json:"messages"`
}

type Client interface {
	GenerateResponse(ctx context.Context, context string, message string) (string, error)
}

var _ Client = (*ClientImpl)(nil)

type ClientImpl struct {
	Token    string
	Endpoint string
	FolderID string
	LastReplyTime time.Time
}

var DefaultEndpoint = "https://llm.api.cloud.yandex.net"

func (c *ClientImpl) GenerateResponse(ctx context.Context, context string, message string) (string, error) {
	client := http.Client{}
	fullURL, err := url.JoinPath(c.Endpoint, "/foundationModels/v1/completion")
	if err != nil {
		return "", err
	}

	generationReq := GenerationRequest{
		ModelURI: fmt.Sprintf("gpt://%s/yandexgpt-lite", c.FolderID),
		Messages: []MessageT{
			{
				Role: "system",
				Text: context,
			},
			{
				Role: "user",
				Text: message,
			},
		},
	}
	generationReq.CompletionOptions.MaxTokens = "2000"
	generationReq.CompletionOptions.Temperature = 0.6
	reqJSON, err := json.Marshal(generationReq)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, bytes.NewReader(reqJSON))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Folder-ID", c.FolderID)
	resp, err := client.Do(req)

	for {
		if resp.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Second)
		} else {
			break
		}
		resp, err = client.Do(req)
	}
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to translate (%d): %s", resp.StatusCode, string(body))
	}

	var respParsed GenerationResponse

	if err := json.Unmarshal(body, &respParsed); err != nil {
		return "", err
	}

	if len(respParsed.Result.Alternatives) == 0 {
		return "", fmt.Errorf("zero alternatives")
	}

	return respParsed.Result.Alternatives[0].Message.Text, nil
}
