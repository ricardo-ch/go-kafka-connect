package connectors

import (
	"log"
	"fmt"
)

//CreateRequest ...
type CreateRequest struct {
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
}

//ConnectorRequest ...
type ConnectorRequest struct {
	Name string
}

//UpdateRequest ...
type UpdateRequest struct {
	Config map[string]string `json:"config"`
}

//GetAllResponse ...
type GetAllResponse struct {
	Code       int
	Connectors []string
}

//ConnectorResponse ...
type ConnectorResponse struct {
	Code   int
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
	Tasks  []TaskID          `json:"tasks"`
}

//ConfigResponse ...
type ConfigResponse struct {
	Code   int
	Config map[string]string
}

//StatusResponse ...
type StatusResponse struct {
	Code            int
	Name            string            `json:"name"`
	ConnectorStatus map[string]string `json:"connector"`
	TasksStatus     []TaskStatus      `json:"tasks"`
}

//EmptyResponse ...
type EmptyResponse struct {
	Code int
}

//GetAll gets the list of all active connectors
func (c Client) GetAll() GetAllResponse {
	var gar GetAllResponse
	var connectors []string

	statusCode, err := c.Request("GET", "connectors", nil, &connectors)
	if err != nil {
		log.Fatal(err)
	}

	gar.Code = statusCode
	gar.Connectors = connectors

	return gar
}

//Create ...
func (c Client) Create(req CreateRequest) ConnectorResponse {
	resp := ConnectorResponse{}

	statusCode, err := c.Request("POST", "connectors", req, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode

	return resp
}

//Get ...
func (c Client) Get(req ConnectorRequest) ConnectorResponse {
	var cr ConnectorResponse

	statusCode, err := c.Request("GET", fmt.Sprintf("connectors/%s", req.Name), nil, &cr)
	if err != nil {
		log.Fatal(err)
	}

	cr.Code = statusCode
	return cr
}

//GetConfig ...
func (c Client) GetConfig(req ConnectorRequest) ConfigResponse {
	var cr ConfigResponse
	var config map[string]string

	statusCode, err := c.Request("GET", fmt.Sprintf("connectors/%s/config", req.Name), nil, &config)
	if err != nil {
		log.Fatal(err)
	}

	cr.Code = statusCode
	cr.Config = config
	return cr
}

//Update ...
func (c Client) Update(req UpdateRequest) ConnectorResponse {

	return ConnectorResponse{}
}

//GetStatus ...
func (c Client) GetStatus(req ConnectorRequest) StatusResponse {
	var sr StatusResponse

	statusCode, err := c.Request("GET", fmt.Sprintf("connectors/%s/status", req.Name), nil, &sr)
	if err != nil {
		log.Fatal(err)
	}

	sr.Code = statusCode
	return sr
}

//Restart ...
func (c Client) Restart(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request("POST", fmt.Sprintf("connectors/%s/restart", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

//Pause ...
func (c Client) Pause(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request("PUT", fmt.Sprintf("connectors/%s/pause", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

//Resume ...
func (c Client) Resume(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request("PUT", fmt.Sprintf("connectors/%s/resume", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

//Delete ...
func (c Client) Delete(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request("DELETE", fmt.Sprintf("connectors/%s", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}
