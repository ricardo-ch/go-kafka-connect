package connectors

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

var endpoint = "/connectors"

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

//HTTPGet handles an HTTP Get request to the client
func (c Client) HTTPGet(adress string) ([]byte, error) {

	var protocol string

	if c.Secure {
		protocol = "https"
	} else {
		protocol = "http"
	}

	res, err := http.Get(protocol + "://" + c.Host + ":" + strconv.Itoa(c.Port) + endpoint + adress)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return []byte(body), nil
}
