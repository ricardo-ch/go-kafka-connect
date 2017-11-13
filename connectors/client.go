package connectors

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"fmt"
)

//Client represents the kafka connect access configuration
type Client struct {
	Host   string
	Port   int
	Secure bool
}

//ErrorResponse is generic error returned by kafka connect
type ErrorResponse struct {
	ErrorCode int    `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
}

func (err ErrorResponse) Error() string{
	return fmt.Sprintf("error code: %d , message: %s", err.ErrorCode, err.Message)
}

//NewClient generates a new client
func NewClient(host string, port int, secure bool) Client {
	return Client{Host: host, Port: port, Secure: secure}
}

//Request handles an HTTP Get request to the client
// execute request and return parsed body content in result var
// result need to be pointer to expected type
func (c Client) Request(method string, endpoint string, request interface{}, result interface{}) (int, error) {
	var protocol string

	if c.Secure {
		protocol = "https"
	} else {
		protocol = "http"
	}

	endPointURL, err := url.Parse(protocol + "://" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + endpoint)
	if err != nil {
		return 0, err
	}

	buf := bytes.Buffer{}
	if request != nil {
		err = json.NewEncoder(&buf).Encode(request)
		if err != nil {
			return 0, err
		}
	}

	req, err := http.NewRequest(method, endPointURL.String(), bytes.NewReader(buf.Bytes()))
	req.Header.Add("Content-Type", "application/json")
	if err != nil {
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	if result != nil && res.Body != nil && res.ContentLength>0 {
		err = json.NewDecoder(res.Body).Decode(result)
		if err != nil {
			return res.StatusCode, err
		}
	}

	return res.StatusCode, nil
}
