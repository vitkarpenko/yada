package images

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type TenorResponse struct {
	Results []struct {
		ID           string `json:"id"`
		MediaFormats struct {
			TinyGIF struct {
				URL string `json:"url"`
			} `json:"tinygif"`
		} `json:"media_formats"`
	} `json:"results"`
}

func (s *Service) RandomGifURL() (string, error) {
	url := fmt.Sprintf("https://g.tenor.com/v2/search?q=%s&key=%s&limit=1", uuid.New().String(), s.tenorAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New("couldn't fetch Tenor")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("couldn't read Tenor response body")
	}

	var data TenorResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", errors.New("couldn't unmarshal Tenor response")
	}

	if len(data.Results) != 1 {
		return "", fmt.Errorf("too much images from tenor: %d", len(data.Results))
	}

	return data.Results[0].MediaFormats.TinyGIF.URL, nil
}
