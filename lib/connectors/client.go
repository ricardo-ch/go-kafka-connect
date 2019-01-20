package connectors

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/resty.v1"
	"time"
)

//Client represents the kafka connect access configuration
type Client struct {
	restClient *resty.Client
}

//ErrorResponse is generic error returned by kafka connect
type ErrorResponse struct {
	ErrorCode int    `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
}

func (err ErrorResponse) Error() string {
	return fmt.Sprintf("error code: %d , message: %s", err.ErrorCode, err.Message)
}

//NewClient generates a new client
func NewClient(url string) *Client {
	restClient := resty.New().
		SetError(ErrorResponse{}).
		SetHostURL(url).
		SetHeader("Accept", "application/json").
		SetRetryCount(3).
		SetTimeout(5 * time.Second)

	return &Client{restClient: restClient}
}

func (c *Client) WithInsecureSSL() *Client {
	return &Client{restClient: c.restClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})}
}

func (c *Client) WithDebug() *Client {
	return &Client{restClient: c.restClient.SetDebug(true)}
}
