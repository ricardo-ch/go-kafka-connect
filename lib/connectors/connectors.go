package connectors

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

//ConnectorRequest is generic request used when interacting with connector endpoint
type ConnectorRequest struct {
	Name string `json:"name"`
}

//EmptyResponse is response returned by multiple endpoint when only StatusCode matter
type EmptyResponse struct {
	Code int
	ErrorResponse
}

//CreateConnectorRequest is request used for creating connector
type CreateConnectorRequest struct {
	ConnectorRequest
	Config map[string]interface{} `json:"config"`
}

//GetAllConnectorsResponse is request used to get list of available connectors
type GetAllConnectorsResponse struct {
	EmptyResponse
	Connectors []string
}

//ConnectorResponse is generic response when interacting with connector endpoint
type ConnectorResponse struct {
	EmptyResponse
	Name   string                 `json:"name"`
	Config map[string]interface{} `json:"config"`
	Tasks  []TaskID               `json:"tasks"`
}

//GetConnectorConfigResponse is response returned by GetConfig endpoint
type GetConnectorConfigResponse struct {
	EmptyResponse
	Config map[string]interface{}
}

//GetConnectorStatusResponse is response returned by GetStatus endpoint
type GetConnectorStatusResponse struct {
	EmptyResponse
	Name            string            `json:"name"`
	ConnectorStatus map[string]string `json:"connector"`
	TasksStatus     []TaskStatus      `json:"tasks"`
}

//GetAll gets the list of all active connectors
func (c *Client) GetAll() (GetAllConnectorsResponse, error) {
	result := GetAllConnectorsResponse{}
	var connectors []string

	resp, err := c.restClient.NewRequest().
		SetResult(&connectors).
		Get("connectors")

	if err != nil {
		return GetAllConnectorsResponse{}, err
	}
	if resp.Error() != nil {
		return GetAllConnectorsResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()
	result.Connectors = connectors

	return result, nil
}

//GetConnector return information on specific connector
func (c Client) GetConnector(req ConnectorRequest) (ConnectorResponse, error) {
	result := ConnectorResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Get("connectors/{name}")
	if err != nil {
		return ConnectorResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 404 && resp.StatusCode() != 409 {
		return ConnectorResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()
	return result, nil
}

//CreateConnector create connector using specified config and name
func (c *Client) CreateConnector(req CreateConnectorRequest, sync bool) (ConnectorResponse, error) {
	result := ConnectorResponse{}

	resp, err := c.restClient.NewRequest().
		SetBody(req).
		SetResult(&result).
		Post("connectors")
	if err != nil {
		return ConnectorResponse{}, err
	}
	if resp.Error() != nil {
		return ConnectorResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()

	if sync {
		if !TryUntil(
			func() bool {
				resp, err := c.GetConnector(req.ConnectorRequest)
				return err == nil && resp.Code == 200
			},
			2*time.Minute,
		) {
			return result, errors.New("timeout on creating connector sync")
		}
	}

	return result, nil
}

//UpdateConnector update a connector config
func (c Client) UpdateConnector(req CreateConnectorRequest, sync bool) (ConnectorResponse, error) {
	result := ConnectorResponse{}

	resp, err := c.restClient.NewRequest().
		SetPathParams(map[string]string{"name": req.Name}).
		SetBody(req.Config).
		SetResult(&result).
		Put("connectors/{name}/config")
	if err != nil {
		return ConnectorResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 409 {
		return ConnectorResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()

	if sync {
		if !TryUntil(
			func() bool {
				upToDate, err := c.IsUpToDate(req.Name, req.Config)
				return err == nil && upToDate
			},
			2*time.Minute,
		) {
			return result, errors.New("timeout on creating connector sync")
		}
	}

	return result, nil
}

//DeleteConnector delete a connector
func (c Client) DeleteConnector(req ConnectorRequest, sync bool) (EmptyResponse, error) {
	result := EmptyResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Delete("connectors/{name}")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 409 {
		return EmptyResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()

	if sync {
		if !TryUntil(
			func() bool {
				r, e := c.GetConnector(req)
				return e == nil && r.Code == 404
			},
			2*time.Minute,
		) {
			return result, errors.New("timeout on deleting connector sync")
		}
	}

	return result, nil
}

////GetConnectorConfig return config of a connector
func (c Client) GetConnectorConfig(req ConnectorRequest) (GetConnectorConfigResponse, error) {
	result := GetConnectorConfigResponse{}
	var config map[string]interface{}

	resp, err := c.restClient.NewRequest().
		SetResult(&config).
		SetPathParams(map[string]string{"name": req.Name}).
		Get("connectors/{name}/config")
	if err != nil {
		return GetConnectorConfigResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 404 && resp.StatusCode() != 409 {
		return GetConnectorConfigResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()
	result.Config = config
	return result, nil
}

//GetConnectorStatus return current status of connector
func (c Client) GetConnectorStatus(req ConnectorRequest) (GetConnectorStatusResponse, error) {
	result := GetConnectorStatusResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Get("connectors/{name}/status")
	if err != nil {
		return GetConnectorStatusResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 404 && resp.StatusCode() != 409 {
		return GetConnectorStatusResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()
	return result, nil
}

//RestartConnector restart connector
func (c Client) RestartConnector(req ConnectorRequest) (EmptyResponse, error) {
	result := EmptyResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Post("connectors/{name}/restart")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 409 {
		return EmptyResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()
	return result, nil
}

//PauseConnector pause a running connector
//asynchronous operation
func (c Client) PauseConnector(req ConnectorRequest, sync bool) (EmptyResponse, error) {
	result := EmptyResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Put("connectors/{name}/pause")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 409 {
		return EmptyResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()

	if sync {
		if !TryUntil(
			func() bool {
				resp, err := c.GetConnectorStatus(req)
				return err == nil && resp.Code == 200 && resp.ConnectorStatus["state"] == "PAUSED"
			},
			2*time.Minute,
		) {
			return result, errors.New("timeout on pausing connector sync")
		}
	}
	return result, nil
}

//ResumeConnector resume a paused connector
//asynchronous operation
func (c Client) ResumeConnector(req ConnectorRequest, sync bool) (EmptyResponse, error) {
	result := EmptyResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Put("connectors/{name}/resume")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 409 {
		return EmptyResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()

	if sync {
		if !TryUntil(
			func() bool {
				resp, err := c.GetConnectorStatus(req)
				return err == nil && resp.Code == 200 && resp.ConnectorStatus["state"] == "RUNNING"
			},
			2*time.Minute,
		) {
			return result, errors.New("timeout on resuming connector sync")
		}
	}
	return result, nil
}

//IsUpToDate checks if the given configuration is different from the deployed one.
//Returns true if they are the same
func (c Client) IsUpToDate(connector string, config map[string]interface{}) (bool, error) {
	config["name"] = connector

	configResp, err := c.GetConnectorConfig(ConnectorRequest{Name: connector})
	if err != nil {
		return false, err
	}
	if configResp.Code == 404 {
		return false, nil
	}
	if configResp.Code >= 400 {
		return false, errors.New(fmt.Sprintf("status code: %d", configResp.Code))
	}

	if len(configResp.Config) != len(config) {
		return false, nil
	}
	for key, value := range configResp.Config {
		if convertConfigValueToString(config[key]) != convertConfigValueToString(value) {
			return false, nil
		}
	}
	return true, nil
}

// Because trying to compare the same field on 2 different config may return false negative if one is encoded as a string and not the other
func convertConfigValueToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	default:
		return ""
	}
}

// TryUntil repeats exec until it return true or timeout is reached
// TryUntil itself return true if `exec` has return true (success), false if timeout (failure)
func TryUntil(exec func() bool, limit time.Duration) bool {
	timeLimit := time.After(limit)

	run := true
	defer func() { run = false }()
	success := make(chan bool)
	go func() {
		for run {
			if exec() {
				success <- true
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	select {
	case <-timeLimit:
		return false
	case <-success:
		return true
	}
}

//DeployConnector checks if the configuration changed before deploying.
//It does nothing if it is the same
func (c Client) DeployConnector(req CreateConnectorRequest) (err error) {
	existingConnector, err := c.GetConnector(ConnectorRequest{Name: req.Name})
	if err != nil {
		return err
	}

	if existingConnector.Code != 404 {
		var upToDate bool
		upToDate, err = c.IsUpToDate(req.Name, req.Config)
		if err != nil {
			return err
		}
		// Connector is already up to date, stop there and return ok
		if upToDate {
			return nil
		}

		_, err = c.PauseConnector(ConnectorRequest{Name: req.Name}, true)
		if err != nil {
			return err
		}

		defer func() {
			_, err = c.ResumeConnector(ConnectorRequest{Name: req.Name}, true)
		}()
	}

	_, err = c.UpdateConnector(req, true)
	if err != nil {
		return err
	}

	return err
}
