package connectors

import (
	"log"
	"fmt"
	"net/http"
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
	Name   string            `json:"name"`
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

	statusCode, err := c.Request(http.MethodGet, "connectors", nil, &connectors)
	if err != nil {
		log.Fatal(err)
	}

	gar.Code = statusCode
	gar.Connectors = connectors

	return gar
}

//Get ...
func (c Client) Get(req ConnectorRequest) ConnectorResponse {
	var cr ConnectorResponse

	statusCode, err := c.Request(http.MethodGet, fmt.Sprintf("connectors/%s", req.Name), nil, &cr)
	if err != nil {
		log.Fatal(err)
	}

	cr.Code = statusCode
	return cr
}

//Create ...
func (c Client) Create(req CreateRequest) ConnectorResponse {
	resp := ConnectorResponse{}

	statusCode, err := c.Request(http.MethodPost, "connectors", req, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode

	return resp
}

//Update ...
func (c Client) Update(req UpdateRequest) ConnectorResponse {
	sr := ConnectorResponse{}

	statusCode, err := c.Request(http.MethodPut, fmt.Sprintf("connectors/%s/config", req.Name), req.Config, &sr)
	if err != nil {
		log.Fatal(err)
	}

	sr.Code = statusCode
	return sr
}

//Delete ...
func (c Client) Delete(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request(http.MethodDelete, fmt.Sprintf("connectors/%s", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

//GetConfig ...
func (c Client) GetConfig(req ConnectorRequest) ConfigResponse {
	var cr ConfigResponse
	var config map[string]string

	statusCode, err := c.Request(http.MethodGet, fmt.Sprintf("connectors/%s/config", req.Name), nil, &config)
	if err != nil {
		log.Fatal(err)
	}

	cr.Code = statusCode
	cr.Config = config
	return cr
}

//GetStatus ...
func (c Client) GetStatus(req ConnectorRequest) StatusResponse {
	var sr StatusResponse

	statusCode, err := c.Request(http.MethodGet, fmt.Sprintf("connectors/%s/status", req.Name), nil, &sr)
	if err != nil {
		log.Fatal(err)
	}

	sr.Code = statusCode
	return sr
}

//Restart ...
func (c Client) Restart(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request(http.MethodPost, fmt.Sprintf("connectors/%s/restart", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

//Pause ...
func (c Client) Pause(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request(http.MethodPut, fmt.Sprintf("connectors/%s/pause", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

//Resume ...
func (c Client) Resume(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request(http.MethodPut, fmt.Sprintf("connectors/%s/resume", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

