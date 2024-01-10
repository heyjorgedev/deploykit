package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/jorgemurta/deploykit"
)

type CLIClient struct {
	HTTPClient *http.Client

	BaseURL string
}

func NewCliClient(baseUrl string) *CLIClient {
	return &CLIClient{
		HTTPClient: &http.Client{
			// Transport: &http2.Transport{
			// 	AllowHTTP: true,
			// 	DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			// 		return net.Dial(network, addr) // h2c-only right now
			// 	},
			// },
			Transport: http.DefaultTransport,
		},
		BaseURL: baseUrl,
	}
}

func (c *CLIClient) url(path string) (url.URL, error) {
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return url.URL{}, err
	}

	// Only allow http and https schemes.
	if u.Scheme != "http" && u.Scheme != "https" {
		return url.URL{}, fmt.Errorf("invalid scheme: %s", u.Scheme)
	}

	return url.URL{
		Scheme: u.Scheme,
		Host:   u.Host,
		Path:   path,
	}, nil

}

func (c *CLIClient) AppsCreate(ctx context.Context, a deploykit.App) (deploykit.App, error) {
	url, err := c.url("/apps")
	if err != nil {
		fmt.Println("a")
		return a, nil
	}

	b, err := json.Marshal(a)
	if err != nil {
		fmt.Println("b")
		return a, err
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("c")
		return a, err
	}
	req = req.WithContext(ctx)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		fmt.Println("d")
		return a, err
	}
	defer resp.Body.Close()

	var app deploykit.App
	json.NewDecoder(resp.Body).Decode(&app)

	return app, nil
}
