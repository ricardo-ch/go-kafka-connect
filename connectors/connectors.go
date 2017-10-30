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
	Name   string            `json:"name"`
	Config map[string]string `json:"config"`
	Tasks  []Task            `json:"tasks"`
}

//Task ...
type Task struct {
	Connector string `json:"connector"`
	TaskID    int    `json:"task"`
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

//TaskStatus ...
type TaskStatus struct {
	ID       int    `json:"id"`
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
	Trace    string `json:"trace,omitempty"`
}

//EmptyResponse ...
type EmptyResponse struct {
	Code int
}

//GetAll gets the list of all active connectors
func (c Client) GetAll() GetAllResponse {
	var gar GetAllResponse
	var connectors []string

	res, err := c.HTTPGet("/")
	if err != nil {
		log.Fatal(err)
	}

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
	var cr ConnectorResponse

	res, err := c.HTTPGet("/" + req.Name)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(res, &cr)
	if err != nil {
		log.Fatal(err)
	}

	cr.Code = 200
	return cr
}

//GetConfig ...
func (c Client) GetConfig(req ConnectorRequest) ConfigResponse {
	var cr ConfigResponse
	var config map[string]string

	res, err := c.HTTPGet("/" + req.Name + "/config")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(res, &config)
	if err != nil {
		log.Fatal(err)
	}

	cr.Code = 200
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

	res, err := c.HTTPGet("/" + req.Name + "/status")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(res, &sr)
	if err != nil {
		log.Fatal(err)
	}

	sr.Code = 200
	return sr
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
