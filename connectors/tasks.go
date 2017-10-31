package connectors

import (
	"encoding/json"
	"log"
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

//TaskResponse ...
type TaskResponse struct {
	Code       int
	TaskStatus map[string]string
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
func (c Client) GetTaskStatus(req TaskRequest) {

}

//RestartTask ...
func (c Client) RestartTask(req TaskRequest) TaskResponse {

	return TaskResponse{}
}
