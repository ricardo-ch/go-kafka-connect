package connectors

import (
	"net/http"
	"strconv"
	"encoding/json"
	"net/url"
	"bytes"
)

//Client represents the kafka connect access configuration
type Client struct {
	Host   string
	Port   int
	Secure bool
}

//NewClient generates a new client
func NewClient(h string, p int, s bool) Client {
	return Client{Host: h, Port: p, Secure: s}
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

	endPointUrl, err := url.Parse(protocol + "://" + c.Host + ":" + strconv.Itoa(c.Port) + endpoint)
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

	req, err := http.NewRequest(method, endPointUrl.String(), bytes.NewReader(buf.Bytes()))
	if err != nil {
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	if result != nil {
		err = json.NewDecoder(res.Body).Decode(result)
		if err != nil {
			return res.StatusCode, err
		}
	}

	return res.StatusCode, nil
}
