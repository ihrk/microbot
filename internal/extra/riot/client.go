package riot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

const riotApiDomain = ".api.riotgames.com"

var regions = map[string]string{
	"euw":  "euw1",
	"ru":   "ru",
	"na":   "na1",
	"br":   "br1",
	"eune": "eun1",
	"jp":   "jp1",
	"kr":   "kr",
	"tr":   "tr1",
	"oce":  "oc1",
	"la1":  "la1",
	"la2":  "la2",
}

var ErrRegionNotFound = errors.New("region not found")

type Client struct {
	host   string
	apiKey string
}

func NewClient(region, apiKey string) (*Client, error) {
	regionDomain, ok := regions[region]
	if !ok {
		return nil, ErrRegionNotFound
	}

	return &Client{
		host:   regionDomain + riotApiDomain,
		apiKey: apiKey,
	}, nil
}

func (c *Client) doRequest(path string, v interface{}) error {
	var u url.URL

	u.Scheme = "https"
	u.Host = c.host
	u.Path = path

	q := u.Query()
	q.Set("api_key", c.apiKey)
	u.RawQuery = q.Encode()

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		var e APIError
		_ = json.NewDecoder(resp.Body).Decode(&e)

		return &e
	}

	return json.NewDecoder(resp.Body).Decode(v)
}

type APIError struct {
	Status Status `json:"status"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("riot API call ended with status: %d, msg: '%s'", e.Status.StatusCode, e.Status.Message)
}

type Status struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}
