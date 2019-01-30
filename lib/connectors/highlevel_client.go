package connectors

import (
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

// HighLevelClient support all function of kafka-connect API + some more features
type HighLevelClient interface {
	// kafka-connect api
	GetAll() (GetAllConnectorsResponse, error)
	GetConnector(req ConnectorRequest) (ConnectorResponse, error)
	CreateConnector(req CreateConnectorRequest, sync bool) (ConnectorResponse, error)
	UpdateConnector(req CreateConnectorRequest, sync bool) (ConnectorResponse, error)
	DeleteConnector(req ConnectorRequest, sync bool) (EmptyResponse, error)
	GetConnectorConfig(req ConnectorRequest) (GetConnectorConfigResponse, error)
	GetConnectorStatus(req ConnectorRequest) (GetConnectorStatusResponse, error)
	RestartConnector(req ConnectorRequest) (EmptyResponse, error)
	PauseConnector(req ConnectorRequest, sync bool) (EmptyResponse, error)
	ResumeConnector(req ConnectorRequest, sync bool) (EmptyResponse, error)
	GetAllTasks(req ConnectorRequest) (GetAllTasksResponse, error)
	GetTaskStatus(req TaskRequest) (TaskStatusResponse, error)
	RestartTask(req TaskRequest) (EmptyResponse, error)

	// custom features, mostly composition of previous ones
	IsUpToDate(connector string, config map[string]interface{}) (bool, error)
	DeployConnector(req CreateConnectorRequest) (err error)
	DeployMultipleConnector(connectors []CreateConnectorRequest) (err error)
	SetInsecureSSL()
	SetDebug()
	SetParallelism(value int)
}

type highLevelClient struct {
	client             BaseClient
	maxParallelRequest int
}

//NewClient generates a new client
func NewClient(url string) HighLevelClient {
	return &highLevelClient{client: newBaseClient(url), maxParallelRequest: 3}
}

//Set the limit of parallel call to kafka-connect server
//Default to 3
func (c *highLevelClient) SetParallelism(value int) {
	c.maxParallelRequest = value
}

func (c *highLevelClient) SetInsecureSSL() {
	c.client.SetInsecureSSL()
}

func (c *highLevelClient) SetDebug() {
	c.client.SetDebug()
}

//GetAll gets the list of all active connectors
func (c *highLevelClient) GetAll() (GetAllConnectorsResponse, error) {
	return c.client.GetAll()
}

//GetConnector return information on specific connector
func (c *highLevelClient) GetConnector(req ConnectorRequest) (ConnectorResponse, error) {
	return c.client.GetConnector(req)
}

//CreateConnector create connector using specified config and name
func (c *highLevelClient) CreateConnector(req CreateConnectorRequest, sync bool) (ConnectorResponse, error) {
	result, err := c.client.CreateConnector(req)
	if err != nil {
		return result, err
	}

	if sync {
		if !tryUntil(
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
func (c *highLevelClient) UpdateConnector(req CreateConnectorRequest, sync bool) (ConnectorResponse, error) {
	result, err := c.client.UpdateConnector(req)
	if err != nil {
		return result, err
	}

	if sync {
		if !tryUntil(
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
func (c *highLevelClient) DeleteConnector(req ConnectorRequest, sync bool) (EmptyResponse, error) {
	result, err := c.client.DeleteConnector(req)
	if err != nil {
		return result, err
	}

	if sync {
		if !tryUntil(
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
func (c *highLevelClient) GetConnectorConfig(req ConnectorRequest) (GetConnectorConfigResponse, error) {
	return c.client.GetConnectorConfig(req)
}

//GetConnectorStatus return current status of connector
func (c *highLevelClient) GetConnectorStatus(req ConnectorRequest) (GetConnectorStatusResponse, error) {
	return c.client.GetConnectorStatus(req)
}

//RestartConnector restart connector
func (c *highLevelClient) RestartConnector(req ConnectorRequest) (EmptyResponse, error) {
	return c.client.RestartConnector(req)
}

//PauseConnector pause a running connector
//asynchronous operation
func (c *highLevelClient) PauseConnector(req ConnectorRequest, sync bool) (EmptyResponse, error) {
	result, err := c.client.PauseConnector(req)
	if err != nil {
		return result, err
	}

	if sync {
		if !tryUntil(
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
func (c *highLevelClient) ResumeConnector(req ConnectorRequest, sync bool) (EmptyResponse, error) {
	result, err := c.client.ResumeConnector(req)
	if err != nil {
		return result, err
	}

	if sync {
		if !tryUntil(
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
func (c *highLevelClient) IsUpToDate(connector string, config map[string]interface{}) (bool, error) {
	// copy the map to safely interact with it
	// we are going to need to add connector name to be able to exact match
	copyConfig := make(map[string]interface{}, len(config))
	for key, value := range config {
		copyConfig[key] = value
	}

	copyConfig["name"] = connector

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

	if len(configResp.Config) != len(copyConfig) {
		return false, nil
	}
	for key, value := range configResp.Config {
		if convertConfigValueToString(copyConfig[key]) != convertConfigValueToString(value) {
			return false, nil
		}
	}
	return true, nil
}

// Because trying to compare the same field on 2 different config may return false negative if one is encoded as a string and not the other
func convertConfigValueToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// tryUntil repeats exec until it return true or timeout is reached
// tryUntil itself return true if `exec` has return true (success), false if timeout (failure)
func tryUntil(exec func() bool, limit time.Duration) bool {
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
func (c *highLevelClient) DeployConnector(req CreateConnectorRequest) (err error) {
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

func (c *highLevelClient) DeployMultipleConnector(connectors []CreateConnectorRequest) (err error) {
	errSync := new(sync.Mutex)
	// Channel is used only to limit number of parallel request
	throttleCh := make(chan interface{}, c.maxParallelRequest)

	for _, connector := range connectors {
		throttleCh <- struct{}{}
		go func(req CreateConnectorRequest) {
			defer func() { <-throttleCh }()
			newErr := c.DeployConnector(req)
			if newErr != nil {
				errSync.Lock()
				defer errSync.Unlock()
				err = multierror.Append(err, errors.Wrapf(newErr, "error while deploying: %v", req.Name))
			}
		}(connector)
	}

	// wait for the end
	for i := 0; i < c.maxParallelRequest; i++ {
		throttleCh <- struct{}{}
	}

	return err
}

// --------------- tasks ---------------------

//GetAllTasks return list of running task
func (c *highLevelClient) GetAllTasks(req ConnectorRequest) (GetAllTasksResponse, error) {
	return c.client.GetAllTasks(req)
}

//GetTaskStatus return current status of task
func (c *highLevelClient) GetTaskStatus(req TaskRequest) (TaskStatusResponse, error) {
	return c.client.GetTaskStatus(req)
}

//RestartTask try to restart task
func (c *highLevelClient) RestartTask(req TaskRequest) (EmptyResponse, error) {
	return c.client.RestartTask(req)
}
