package connectors

import (
	"log"
	"fmt"
	"net/http"
)

//TaskRequest ...
type TaskRequest struct {
	Connector string
	TaskID    int
}

//GetAllTasksResponse ...
type GetAllTasksResponse struct {
	Code  int
	Tasks []TaskDetails
}

//TaskDetails ...
type TaskDetails struct {
	ID     TaskID            `json:"id"`
	Config map[string]string `json:"config"`
}

//TaskID ...
type TaskID struct {
	Connector string `json:"connector"`
	TaskID    int    `json:"task"`
}

//TaskStatusResponse ...
type TaskStatusResponse struct {
	Code   int
	Status TaskStatus
}

//TaskStatus ...
type TaskStatus struct {
	ID       int    `json:"id"`
	State    string `json:"state"`
	WorkerID string `json:"worker_id"`
	Trace    string `json:"trace,omitempty"`
}

//GetAllTasks ...
func (c Client) GetAllTasks(req ConnectorRequest) GetAllTasksResponse {
	var gatr GetAllTasksResponse
	var taskDetails []TaskDetails

	statusCode, err := c.Request(http.MethodGet, fmt.Sprintf("connectors/%s/tasks", req.Name), nil, &taskDetails)
	if err != nil {
		log.Fatal(err)
	}

	gatr.Code = statusCode
	gatr.Tasks = taskDetails
	return gatr
}

//GetTaskStatus ...
func (c Client) GetTaskStatus(req TaskRequest) TaskStatusResponse {
	var tsr TaskStatusResponse
	var ts TaskStatus

	statusCode, err := c.Request(http.MethodGet, fmt.Sprintf( "connectors/%s/tasks/%s/status", req.Connector, req.TaskID), nil, &ts)
	if err != nil {
		log.Fatal(err)
	}

	tsr.Code = statusCode
	tsr.Status = ts

	return tsr
}

//RestartTask ...
func (c Client) RestartTask(req TaskRequest) EmptyResponse {
	var er EmptyResponse

	statusCode, err := c.Request(http.MethodGet, fmt.Sprintf("connectors/%s/tasks/%s/restart", req.Connector, req.TaskID ), nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	er.Code = statusCode

	return er
}
