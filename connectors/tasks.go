package connectors

import (
	"encoding/json"
	"log"
	"strconv"
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

	res, err := c.HTTPGet("/" + req.Name + "/tasks")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(res, &taskDetails)
	if err != nil {
		log.Fatal(err)
	}

	gatr.Code = 200
	gatr.Tasks = taskDetails
	return gatr
}

//GetTaskStatus ...
func (c Client) GetTaskStatus(req TaskRequest) TaskStatusResponse {
	var tsr TaskStatusResponse
	var ts TaskStatus

	res, err := c.HTTPGet("/" + req.Connector + "/tasks/" + strconv.Itoa(req.TaskID) + "/status")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(res, &ts)
	if err != nil {
		log.Fatal(err)
	}

	tsr.Code = 200
	tsr.Status = ts

	return tsr
}

//RestartTask ...
func (c Client) RestartTask(req TaskRequest) EmptyResponse {
	var er EmptyResponse

	_, err := c.HTTPGet("/" + req.Connector + "/tasks/" + strconv.Itoa(req.TaskID) + "/restart")
	if err != nil {
		log.Fatal(err)
	}

	er.Code = 200

	return er
}
