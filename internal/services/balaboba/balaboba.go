package balaboba

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Balaboba struct {
	headers    map[string]string
	apiURL     string
	httpClient http.Client
}

type Response struct {
	Text string `json:"text"`
}

func NewBalaboba() *Balaboba {
	rand.Seed(time.Now().UnixNano())

	return &Balaboba{
		headers: map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_4) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.1 Safari/605.1.15",
			"Origin":       "https://yandex.ru",
			"Referer":      "https://yandex.ru/",
		},
		apiURL:     "https://zeapi.yandex.net/lab/api/yalm/text3",
		httpClient: http.Client{Timeout: 20 * time.Second},
	}
}

func (b *Balaboba) GenerateText(message string) string {
	marshalledBody, err := marshalBody(message)
	if err != nil {
		log.Println("Couldn't marshal balaboba's request body", marshalledBody)
	}

	request := b.prepareRequest(marshalledBody)
	resp := b.doRequest(request)

	return readResponse(resp)
}

func readResponse(resp *http.Response) string {
	text, _ := io.ReadAll(resp.Body)
	var result Response
	_ = json.Unmarshal(text, &result)
	return result.Text
}

func (b *Balaboba) doRequest(req *http.Request) *http.Response {
	resp, err := b.httpClient.Do(req)
	if err != nil {
		log.Println("Error calling Balaboba", err)
	}
	return resp
}

func (b *Balaboba) prepareRequest(marshalledBody []byte) *http.Request {
	request, _ := http.NewRequest("POST", b.apiURL, bytes.NewReader(marshalledBody))
	for key, value := range b.headers {
		request.Header.Set(key, value)
	}
	return request
}

func marshalBody(message string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"filter": rand.Intn(10),
		"intro":  4,
		"query":  message,
	})
}
