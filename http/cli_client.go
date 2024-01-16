package http

import (
	"fmt"
	"net/http"
	"net/url"
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
