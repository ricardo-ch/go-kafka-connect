package connectors

import (
	"io/ioutil"
	"net/http"
	"strconv"
)

//Client ...
type Client struct {
	Host string
	Port int
}

//NewClient ...
func NewClient(h string, p int) Client {
	return Client{Host: h, Port: p}
}

//HTTPGet handles an HTTP Get request to the client
func (c Client) HTTPGet(adress string) ([]byte, error) {

	res, err := http.Get("http://" + c.Host + ":" + strconv.Itoa(c.Port) + "/connectors")
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
