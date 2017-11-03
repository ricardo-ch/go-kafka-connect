package connectors

import (
	"log"
	"fmt"
	"net/http"
)

//ConnectorRequest is generic request used when interacting with connector endpoint
type ConnectorRequest struct {
	Name string  `json:"name"`
}

//CreateConnectorRequest is request used for creating connector
type CreateConnectorRequest struct {
	ConnectorRequest
	Config map[string]string `json:"config"`
}

//UpdateRequest is request used for updating connector
type UpdateRequest struct {
	ConnectorRequest
	Config map[string]string `json:"config"`
}

//GetAllConnectorsResponse is request used to get list of available connectors
type GetAllConnectorsResponse struct {
	Code       int
	Connectors []string
}

//ConnectorResponse is generic response when interacting with connector endpoint
type ConnectorResponse struct {
	Code   int
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
	Tasks  []TaskID          `json:"tasks"`
}

//GetConfigResponse is response returned by GetConfig endpoint
type GetConfigResponse struct {
	Code   int
	Config map[string]string
}

//GetStatusResponse is response returned by GetStatus endpoint
type GetStatusResponse struct {
	Code            int
	Name            string            `json:"name"`
	ConnectorStatus map[string]string `json:"connector"`
	TasksStatus     []TaskStatus      `json:"tasks"`
}

//EmptyResponse is response returned by multiple endpoint when only StatusCode matter
type EmptyResponse struct {
	Code int
}

//GetAll gets the list of all active connectors
func (c Client) GetAll() GetAllConnectorsResponse {
	var gar GetAllConnectorsResponse
	var connectors []string

	statusCode, err := c.Request(http.MethodGet, "connectors", nil, &connectors)
	if err != nil {
		log.Fatal(err)
	}

	gar.Code = statusCode
	gar.Connectors = connectors

	return gar
}

//GetConnector return information on specific connector
func (c Client) GetConnector(req ConnectorRequest) ConnectorResponse {
	var cr ConnectorResponse

	statusCode, err := c.Request(http.MethodGet, fmt.Sprintf("connectors/%s", req.Name), nil, &cr)
	if err != nil {
		log.Fatal(err)
	}

	cr.Code = statusCode
	return cr
}

//CreateConnector create connector using specified config and name
func (c Client) CreateConnector(req CreateConnectorRequest) ConnectorResponse {
	resp := ConnectorResponse{}

	statusCode, err := c.Request(http.MethodPost, "connectors", req, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode

	return resp
}

//UpdateConnector update a connector config
func (c Client) UpdateConnector(req UpdateRequest) ConnectorResponse {
	sr := ConnectorResponse{}

	statusCode, err := c.Request(http.MethodPut, fmt.Sprintf("connectors/%s/config", req.Name), req.Config, &sr)
	if err != nil {
		log.Fatal(err)
	}

	sr.Code = statusCode
	return sr
}

//DeleteConnector delete a connector
func (c Client) DeleteConnector(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request(http.MethodDelete, fmt.Sprintf("connectors/%s", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

//GetConnectorConfig return config of a connector
func (c Client) GetConnectorConfig(req ConnectorRequest) GetConfigResponse {
	var cr GetConfigResponse
	var config map[string]string

	statusCode, err := c.Request(http.MethodGet, fmt.Sprintf("connectors/%s/config", req.Name), nil, &config)
	if err != nil {
		log.Fatal(err)
	}

	cr.Code = statusCode
	cr.Config = config
	return cr
}

//GetConnectorStatus return current status of connector
func (c Client) GetConnectorStatus(req ConnectorRequest) GetStatusResponse {
	var sr GetStatusResponse

	statusCode, err := c.Request(http.MethodGet, fmt.Sprintf("connectors/%s/status", req.Name), nil, &sr)
	if err != nil {
		log.Fatal(err)
	}

	sr.Code = statusCode
	return sr
}

//RestartConnector restart connector
func (c Client) RestartConnector(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request(http.MethodPost, fmt.Sprintf("connectors/%s/restart", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

//PauseConnector pause a running connector
func (c Client) PauseConnector(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request(http.MethodPut, fmt.Sprintf("connectors/%s/pause", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

//ResumeConnector resume a paused connector
func (c Client) ResumeConnector(req ConnectorRequest) EmptyResponse {
	resp := EmptyResponse{}

	statusCode, err := c.Request(http.MethodPut, fmt.Sprintf("connectors/%s/resume", req.Name), nil, &resp)
	if err != nil {
		log.Fatal(err)
	}

	resp.Code = statusCode
	return resp
}

