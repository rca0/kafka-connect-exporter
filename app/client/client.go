package client

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// TODO:
// ADD KAFKA AUTHS (SASL, TLS)

type Client interface {
	Get(path string) (*http.Response, error)
}

type client struct {
	client  http.Client
	baseUrl string
}

func NewClient(baseUrl string) Client {
	return &client{
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		baseUrl: baseUrl,
	}
}

func (c *client) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", c.baseUrl+path, nil)
	if err != nil {
		logrus.Errorf("failed creating request: %v", err)
		return nil, err
	}
	return c.client.Do(req)
}
