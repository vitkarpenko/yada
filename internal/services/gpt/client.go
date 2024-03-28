package gpt

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"github.com/rs/zerolog/log"
)

const (
	gptURL             = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
	liteModelMaxTokens = 2000
)

type gptRequestBody struct {
	CompletionOptions completionOptions `json:"completionOptions,omitempty"`
	Messages          messages          `json:"messages,omitempty"`
	ModelUri          string            `json:"modelUri,omitempty"`
}
type completionOptions struct {
	MaxTokens   int     `json:"maxTokens,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
}

type messages []message

type message struct {
	Role string `json:"role,omitempty"`
	Text string `json:"text,omitempty"`
}

type gptResponse struct {
	Result struct {
		Alternatives []alternative `json:"alternatives,omitempty"`
		ModelVersion string        `json:"modelVersion,omitempty"`
		Usage        modelUsage    `json:"usage,omitempty"`
	} `json:"result,omitempty"`
}

type alternative struct {
	Message respMessage `json:"message,omitempty"`
	Status  string      `json:"status,omitempty"`
}

type respMessage struct {
	Role string `json:"role,omitempty"`
	Text string `json:"text,omitempty"`
}

type modelUsage struct {
	CompletionTokens string `json:"completionTokens,omitempty"`
	InputTextTokens  string `json:"inputTextTokens,omitempty"`
	TotalTokens      string `json:"totalTokens,omitempty"`
}

type client struct {
	token     string
	catalogID string
	prompt    string

	httpClient *http.Client
}

func newClient(token, catalogID string) client {
	return client{
		token:     token,
		catalogID: catalogID,
		prompt:    "Ответь так, словно ты - душный, занудливый и въедливый член ордена Адептус Механикус из Warhammer 40k",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *client) reply(text string) (string, error) {
	reqBody := gptRequestBody{
		ModelUri: fmt.Sprintf("gpt://%s/yandexgpt-lite", c.catalogID),
		CompletionOptions: completionOptions{
			MaxTokens:   liteModelMaxTokens,
			Temperature: 0.6,
		},
		Messages: []message{
			{
				Role: "user",
				Text: text,
			},
			{
				Role: "system",
				Text: c.prompt,
			},
		},
	}
	log.Debug().Any("req", reqBody).Msg("yandex gpt req")

	marshalled, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshalling yandex gpt request body: %w", err)
	}

	r, err := http.NewRequest(http.MethodPost, gptURL, bytes.NewBuffer(marshalled))
	if err != nil {
		return "", fmt.Errorf("preparing yandex gpt request: %w", err)
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Authorization", "Api-Key "+c.token)

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return "", fmt.Errorf("sending yandex gpt request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading yandex gpt resp body: %w", err)
	}

	var data gptResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", fmt.Errorf("unmarshalling yandex gpt resp body: %w", err)
	}
	log.Debug().Any("resp", data).Msg("yandex gpt resp")

	if len(data.Result.Alternatives) == 0 {
		return "", errors.New("empty yandex gpt response")
	}

	return data.Result.Alternatives[0].Message.Text, nil
}
