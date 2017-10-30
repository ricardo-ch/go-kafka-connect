package connectors

//TaskRequest ...
type TaskRequest struct {
	Connector string
	TaskID    int
}

//GetAllTasksresponse ...
type GetAllTasksresponse struct {
	Code  int
	Tasks []Task
}

//TaskResponse ...
type TaskResponse struct {
	Code       int
	TaskStatus map[string]string
}

//Task ...
type Task struct {
	Task   map[string]string `json:"id"`
	Config map[string]string `json:"config"`
}

//GetAllTasks ...
func (c Client) GetAllTasks(req ConnectorRequest) GetAllTasksresponse {

	return GetAllTasksresponse{}
}

//GetTaskStatus ...
func (c Client) GetTaskStatus(req TaskRequest) {

}

//RestartTask ...
func (c Client) RestartTask(req TaskRequest) TaskResponse {

	return TaskResponse{}
}
