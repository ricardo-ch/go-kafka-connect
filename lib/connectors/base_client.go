package connectors

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
	"strconv"
	"time"
)

// BaseClient implement the kafka-connect contract as a client
// handle retries on 409 response
type BaseClient interface {
	GetAll() (GetAllConnectorsResponse, error)
	GetConnector(req ConnectorRequest) (ConnectorResponse, error)
	CreateConnector(req CreateConnectorRequest) (ConnectorResponse, error)
	UpdateConnector(req CreateConnectorRequest) (ConnectorResponse, error)
	DeleteConnector(req ConnectorRequest) (EmptyResponse, error)
	GetConnectorConfig(req ConnectorRequest) (GetConnectorConfigResponse, error)
	GetConnectorStatus(req ConnectorRequest) (GetConnectorStatusResponse, error)
	RestartConnector(req ConnectorRequest) (EmptyResponse, error)
	PauseConnector(req ConnectorRequest) (EmptyResponse, error)
	ResumeConnector(req ConnectorRequest) (EmptyResponse, error)
	GetAllTasks(req ConnectorRequest) (GetAllTasksResponse, error)
	GetTaskStatus(req TaskRequest) (TaskStatusResponse, error)
	RestartTask(req TaskRequest) (EmptyResponse, error)

	SetInsecureSSL()
	SetDebug()
}

type baseClient struct {
	restClient *resty.Client
}

func (c *baseClient) SetInsecureSSL() {
	c.restClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
}

func (c *baseClient) SetDebug() {
	c.restClient.SetDebug(true)
}

//ErrorResponse is generic error returned by kafka connect
type ErrorResponse struct {
	ErrorCode int    `json:"error_code,omitempty"`
	Message   string `json:"message,omitempty"`
}

func (err ErrorResponse) Error() string {
	return fmt.Sprintf("error code: %d , message: %s", err.ErrorCode, err.Message)
}

func newBaseClient(url string) BaseClient {
	restClient := resty.New().
		SetError(ErrorResponse{}).
		SetHostURL(url).
		SetHeader("Accept", "application/json").
		SetRetryCount(5).
		SetRetryWaitTime(500 * time.Millisecond).
		SetRetryMaxWaitTime(5 * time.Second).
		SetTimeout(10 * time.Second).
		AddRetryCondition(func(resp *resty.Response) (bool, error) {
			return resp.StatusCode() == 409, nil
		})

	return &baseClient{restClient: restClient}
}

// ------------- Connectors ------------

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
func (c *baseClient) GetAll() (GetAllConnectorsResponse, error) {
	result := GetAllConnectorsResponse{}
	var connectors []string

	resp, err := c.restClient.NewRequest().
		SetResult(&connectors).
		Get("connectors")

	if err != nil {
		return GetAllConnectorsResponse{}, err
	}
	if resp.StatusCode() >= 400 {
		return GetAllConnectorsResponse{}, errors.Errorf("Get all connector : %v", resp.String())
	}

	result.Code = resp.StatusCode()
	result.Connectors = connectors

	return result, nil
}

//GetConnector return information on specific connector
func (c *baseClient) GetConnector(req ConnectorRequest) (ConnectorResponse, error) {
	result := ConnectorResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Get("connectors/{name}")
	if err != nil {
		return ConnectorResponse{}, err
	}

	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return ConnectorResponse{}, errors.Errorf("Get connector : %v", resp.String())
	}

	result.Code = resp.StatusCode()
	return result, nil
}

//CreateConnector create connector using specified config and name
func (c *baseClient) CreateConnector(req CreateConnectorRequest) (ConnectorResponse, error) {
	result := ConnectorResponse{}

	resp, err := c.restClient.NewRequest().
		SetBody(req).
		SetResult(&result).
		Post("connectors")
	if err != nil {
		return ConnectorResponse{}, err
	}
	if resp.StatusCode() >= 400 {
		return ConnectorResponse{}, errors.Errorf("Create connector : %v", resp.String())
	}

	result.Code = resp.StatusCode()

	return result, nil
}

//UpdateConnector update a connector config
func (c *baseClient) UpdateConnector(req CreateConnectorRequest) (ConnectorResponse, error) {
	result := ConnectorResponse{}

	resp, err := c.restClient.NewRequest().
		SetPathParams(map[string]string{"name": req.Name}).
		SetBody(req.Config).
		SetResult(&result).
		Put("connectors/{name}/config")
	if err != nil {
		return ConnectorResponse{}, err
	}
	if resp.StatusCode() >= 400 {
		return ConnectorResponse{}, errors.Errorf("Update connector : %v", resp.String())
	}

	result.Code = resp.StatusCode()

	return result, nil
}

//DeleteConnector delete a connector
func (c *baseClient) DeleteConnector(req ConnectorRequest) (EmptyResponse, error) {
	result := EmptyResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Delete("connectors/{name}")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.StatusCode() >= 400 {
		return EmptyResponse{}, errors.Errorf("Delete connector : %v", resp.String())
	}

	result.Code = resp.StatusCode()

	return result, nil
}

////GetConnectorConfig return config of a connector
func (c *baseClient) GetConnectorConfig(req ConnectorRequest) (GetConnectorConfigResponse, error) {
	result := GetConnectorConfigResponse{}
	var config map[string]interface{}

	resp, err := c.restClient.NewRequest().
		SetResult(&config).
		SetPathParams(map[string]string{"name": req.Name}).
		Get("connectors/{name}/config")
	if err != nil {
		return GetConnectorConfigResponse{}, err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return GetConnectorConfigResponse{}, errors.Errorf("Get connector config : %v", resp.String())
	}

	result.Code = resp.StatusCode()
	result.Config = config
	return result, nil
}

//GetConnectorStatus return current status of connector
func (c *baseClient) GetConnectorStatus(req ConnectorRequest) (GetConnectorStatusResponse, error) {
	result := GetConnectorStatusResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Get("connectors/{name}/status")
	if err != nil {
		return GetConnectorStatusResponse{}, err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return GetConnectorStatusResponse{}, errors.Errorf("Get connector status : %v", resp.String())
	}

	result.Code = resp.StatusCode()
	return result, nil
}

//RestartConnector restart connector
func (c *baseClient) RestartConnector(req ConnectorRequest) (EmptyResponse, error) {
	result := EmptyResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Post("connectors/{name}/restart")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.StatusCode() >= 400 {
		return EmptyResponse{}, errors.Errorf("Restart connector : %v", resp.String())
	}

	result.Code = resp.StatusCode()
	return result, nil
}

//PauseConnector pause a running connector
//asynchronous operation
func (c *baseClient) PauseConnector(req ConnectorRequest) (EmptyResponse, error) {
	result := EmptyResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Put("connectors/{name}/pause")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.StatusCode() >= 400 {
		return EmptyResponse{}, errors.Errorf("Pause connector : %v", resp.String())
	}

	result.Code = resp.StatusCode()

	return result, nil
}

//ResumeConnector resume a paused connector
//asynchronous operation
func (c *baseClient) ResumeConnector(req ConnectorRequest) (EmptyResponse, error) {
	result := EmptyResponse{}

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Name}).
		Put("connectors/{name}/resume")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.StatusCode() >= 400 {
		return EmptyResponse{}, errors.Errorf("Resume connector : %v", resp.String())
	}

	result.Code = resp.StatusCode()

	return result, nil
}

// ----------- Tasks ---------

//TaskRequest is generic request when interacting with task endpoint
type TaskRequest struct {
	Connector string
	TaskID    int
}

//GetAllTasksResponse is response to get all tasks of a specific endpoint
type GetAllTasksResponse struct {
	Code  int
	Tasks []TaskDetails
}

//TaskDetails is detail of a specific task on a specific endpoint
type TaskDetails struct {
	ID     TaskID                 `json:"id"`
	Config map[string]interface{} `json:"config"`
}

//TaskID identify a task and its connector
type TaskID struct {
	Connector string `json:"connector"`
	TaskID    int    `json:"task"`
}

//TaskStatusResponse is response returned by get task status endpoint
type TaskStatusResponse struct {
	Code   int
	Status TaskStatus
}

//TaskStatus define task status
type TaskStatus struct {
	ID       int    `json:"id"`
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
	Trace    string `json:"trace,omitempty"`
}

//GetAllTasks return list of running task
func (c *baseClient) GetAllTasks(req ConnectorRequest) (GetAllTasksResponse, error) {
	var result GetAllTasksResponse

	resp, err := c.restClient.NewRequest().
		SetResult(&result.Tasks).
		SetPathParams(map[string]string{"name": req.Name}).
		Get("connectors/{name}/tasks")
	if err != nil {
		return GetAllTasksResponse{}, err
	}
	if resp.StatusCode() >= 400 {
		return GetAllTasksResponse{}, errors.Errorf("Get all tasks : %v", resp.String())
	}

	result.Code = resp.StatusCode()
	return result, nil
}

//GetTaskStatus return current status of task
func (c *baseClient) GetTaskStatus(req TaskRequest) (TaskStatusResponse, error) {
	var result TaskStatusResponse

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Connector, "task_id": strconv.Itoa(req.TaskID)}).
		Get("connectors/{name}/tasks/{task_id}/status")
	if err != nil {
		return TaskStatusResponse{}, err
	}
	if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
		return TaskStatusResponse{}, errors.Errorf("Get task status : %v", resp.String())
	}

	result.Code = resp.StatusCode()

	return result, nil
}

//RestartTask try to restart task
func (c *baseClient) RestartTask(req TaskRequest) (EmptyResponse, error) {
	var result EmptyResponse

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Connector, "task_id": strconv.Itoa(req.TaskID)}).
		Post("connectors/{name}/tasks/{task_id}/restart")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.StatusCode() >= 400 {
		return EmptyResponse{}, errors.Errorf("Restart task : %v", resp.String())
	}

	result.Code = resp.StatusCode()

	return result, nil
}
