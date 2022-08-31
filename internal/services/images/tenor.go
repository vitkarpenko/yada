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
			MP4 struct {
				URL string `json:"url"`
			} `json:"mp4"`
		} `json:"media_formats"`
	} `json:"results"`
}

func (s *Service) Random() (Body, error) {
	url := fmt.Sprintf("https://g.tenor.com/v2/search?q=%s&key=%s&limit=1", uuid.New().String(), s.tenorAPIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.New("couldn't fetch Tenor")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("couldn't read Tenor response body")
	}

	var data TenorResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, errors.New("couldn't unmarshal Tenor response")
	}

	if len(data.Results) != 1 {
		return nil, fmt.Errorf("too much images from tenor: %d", len(data.Results))
	}

	resp, err = http.Get(data.Results[0].MediaFormats.MP4.URL)
	if err != nil {
		return nil, errors.New("couldn't fetch image from Tenor URL")
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("couldn't read Tenor image body")
	}

	return body, nil
}
