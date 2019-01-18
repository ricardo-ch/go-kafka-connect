package connectors

import (
	"crypto/tls"
	"encoding/json"
	"errors"
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
		OnAfterResponse(func(c *resty.Client, res *resty.Response) error {
			// The default error handling given by `SetRESTMode` is a bit weak. This is the override

			if res.StatusCode() >= 400 && res.StatusCode() != 404 {
				restErr := ErrorResponse{}
				decodeErr := json.Unmarshal(res.Body(), &restErr)
				if decodeErr != nil {
					return restErr
				}
				return errors.New(fmt.Sprintf("Error while decoding body while error: %v", res.Body()))
			}
			return nil
		}).
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
