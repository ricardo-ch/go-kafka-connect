package connectors

import (
	"strconv"
)

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
func (c Client) GetAllTasks(req ConnectorRequest) (GetAllTasksResponse, error) {
	var result GetAllTasksResponse

	resp, err := c.restClient.NewRequest().
		SetResult(&result.Tasks).
		SetPathParams(map[string]string{"name": req.Name}).
		Get("connectors/{name}/tasks")
	if err != nil {
		return GetAllTasksResponse{}, err
	}
	if resp.Error() != nil {
		return GetAllTasksResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()
	return result, nil
}

//GetTaskStatus return current status of task
func (c Client) GetTaskStatus(req TaskRequest) (TaskStatusResponse, error) {
	var result TaskStatusResponse

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Connector, "task_id": strconv.Itoa(req.TaskID)}).
		Get("connectors/{name}/tasks/{task_id}/status")
	if err != nil {
		return TaskStatusResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 404 && resp.StatusCode() != 409 {
		return TaskStatusResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()

	return result, nil
}

//RestartTask try to restart task
func (c Client) RestartTask(req TaskRequest) (EmptyResponse, error) {
	var result EmptyResponse

	resp, err := c.restClient.NewRequest().
		SetResult(&result).
		SetPathParams(map[string]string{"name": req.Connector, "task_id": strconv.Itoa(req.TaskID)}).
		Post("connectors/{name}/tasks/{task_id}/restart")
	if err != nil {
		return EmptyResponse{}, err
	}
	if resp.Error() != nil && resp.StatusCode() != 409 {
		return EmptyResponse{}, resp.Error().(*ErrorResponse)
	}

	result.Code = resp.StatusCode()

	return result, nil
}
