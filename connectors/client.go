package connectors

//Client ...
type Client struct {
	Host string
	Port int
}

//NewClient ...
func NewClient(h string, p int) Client {
	return Client{Host: h, Port: p}
}
