package connectors

import (
	"crypto/tls"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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
	SetClientCertificates(certs ...tls.Certificate)
	SetParallelism(value int)
	SetBasicAuth(username string, password string)
	SetHeader(name string, value string)
	SetPauseBeforeDeploy(pauseBeforeDeploy bool)
	SetLogFormatter(formatter logrus.Formatter)
}

type highLevelClient struct {
	client             BaseClient
	maxParallelRequest int
	pauseBeforeDeploy  bool
	logger             *logrus.Entry
}

//NewClient generates a new client
func NewClient(url string) HighLevelClient {
	return &highLevelClient{
		client:             newBaseClient(url),
		maxParallelRequest: 3,
		logger:             logrus.WithField("component", "high-level-client"),
	}
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

func (c *highLevelClient) SetLogFormatter(formatter logrus.Formatter) {
	c.logger.Logger.Formatter = formatter
}

func (c *highLevelClient) SetClientCertificates(certs ...tls.Certificate) {
	c.client.SetClientCertificates(certs...)
}

func (c *highLevelClient) SetBasicAuth(username string, password string) {
	c.client.SetBasicAuth(username, password)
}

func (c *highLevelClient) SetPauseBeforeDeploy(pauseBeforeDeploy bool) {
	c.pauseBeforeDeploy = pauseBeforeDeploy
}

func (c *highLevelClient) SetHeader(name string, value string) {
	c.client.SetHeader(name, value)
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

	mu := &sync.Mutex{}
	run := true

	defer func() {
		mu.Lock()
		defer mu.Unlock()
		run = false
	}()

	success := make(chan bool)
	go func() {
		for {
			mu.Lock()
			if !run {
				break
			}
			mu.Unlock()

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
	logger := c.logger.WithField("connector", req.Name)
	logger.Info("Connector deployment starting...")

	existingConnector, err := c.GetConnector(ConnectorRequest{Name: req.Name})
	if err != nil {
		return err
	}

	if existingConnector.Code != 404 {
		logger.Info("Connector already exists")

		var upToDate bool
		upToDate, err = c.IsUpToDate(req.Name, req.Config)
		if err != nil {
			return err
		}

		// Connector is already up to date, stop there and return ok
		if upToDate {
			logger.Info("Connector is up to date, skipping update")
			return nil
		}

		if !c.pauseBeforeDeploy {
			logger.Info("Connector pause before deploy skipped per configuration")
			goto updateConnector
		}

		existingConnectorStatus, err := c.GetConnectorStatus(ConnectorRequest{Name: req.Name})
		if err != nil {
			return err
		}

		if existingConnectorStatus.ConnectorStatus["state"] == "RUNNING" {
			logger = logger.WithField("pause", strconv.FormatBool(true))

			logger.Info("Connector status is RUNNING, pausing it\n")
			_, err = c.PauseConnector(ConnectorRequest{Name: req.Name}, true)
			logger.WithError(err).Info("Connector status is now PAUSED")
			if err != nil {
				return err
			}

			defer func() {
				logger.Info("Connector status was RUNNING, resuming it")
				_, err = c.ResumeConnector(ConnectorRequest{Name: req.Name}, true)
				logger.WithError(err).Info("Connector is now RUNNING")
			}()
		} else {
			logger.Infof("Connector status '%s' is NOT 'RUNNING', connector will NOT be paused\n", existingConnectorStatus.ConnectorStatus["state"])
		}
	}

updateConnector:

	logger.Info("Connector update starting...")
	_, err = c.UpdateConnector(req, true)
	logger.WithError(err).Info("Connector is now updated")

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
