package connectors

import (
	"encoding/json"
	"log"
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
	Name   string              `json:"name"`
	Config map[string]string   `json:"config"`
	Tasks  []map[string]string `json:"tasks"`
}

//ConfigResponse ...
type ConfigResponse struct {
	Code   int
	Config map[string]string `json:"config"`
}

//StatusResponse ...
type StatusResponse struct {
	Code            int
	Name            string              `json:"name"`
	ConnectorStatus map[string]string   `json:"connector"`
	TasksStatus     []map[string]string `json:"tasks"`
}

//EmptyResponse ...
type EmptyResponse struct {
	Code int
}

//GetAll ...
func (c Client) GetAll() GetAllResponse {
	var gar GetAllResponse
	var connectors []string

	res, err := c.HTTPGet("/connectors")
	err = json.Unmarshal(res, &connectors)
	if err != nil {
		log.Fatal(err)
	}

	gar.Code = 200
	gar.Connectors = connectors

	return gar
}

//Create ...
func (c Client) Create(req CreateRequest) ConnectorResponse {

	return ConnectorResponse{}
}

//Get ...
func (c Client) Get(req ConnectorRequest) ConnectorResponse {

	return ConnectorResponse{}
}

//GetConfig ...
func (c Client) GetConfig(req ConnectorRequest) ConfigResponse {

	return ConfigResponse{}
}

//Update ...
func (c Client) Update(req UpdateRequest) ConnectorResponse {

	return ConnectorResponse{}
}

//GetStatus ...
func (c Client) GetStatus(req ConnectorRequest) StatusResponse {

	return StatusResponse{}
}

//Restart ...
func (c Client) Restart(req ConnectorRequest) EmptyResponse {

	return EmptyResponse{}
}

//Pause ...
func (c Client) Pause(req ConnectorRequest) EmptyResponse {

	return EmptyResponse{}
}

//Resume ...
func (c Client) Resume(req ConnectorRequest) EmptyResponse {

	return EmptyResponse{}
}

//Delete ...
func (c Client) Delete(req ConnectorRequest) EmptyResponse {

	return EmptyResponse{}
}
