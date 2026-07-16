package operator

import (
	"net/http"
	"time"
)

type Client struct {
	client  *http.Client
}

func New() *Client {
	return &Client{
		client:  &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}