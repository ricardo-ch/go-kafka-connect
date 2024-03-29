//go:build !unit

package connectors

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testFile    = "/etc/kafka-connect/kafka-connect.properties"
	hostConnect = "http://localhost:8083"
)

func TestHealthz(t *testing.T) {
	resp, err := http.Get(hostConnect)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCreateConnector(t *testing.T) {
	client := NewClient(hostConnect)
	resp, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-create-connector"},
			Config: map[string]interface{}{
				"connector.class": "FileStreamSource",
				"file":            testFile,
				"topic":           "connect-test",
			},
		},
		true,
	)

	assert.Nil(t, err)
	assert.Equal(t, 201, resp.Code)
}

func TestErrorCode(t *testing.T) {
	client := NewClient(hostConnect)
	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "not-a-valid-connector"},
			Config: map[string]interface{}{
				"connector.class": "not a valid connector class",
				"file":            testFile,
				"topic":           "connect-test",
			},
		},
		true,
	)

	assert.Error(t, err)
}

func TestGetConnector(t *testing.T) {
	client := NewClient(hostConnect)
	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-get-connector"},
			Config: map[string]interface{}{
				"connector.class": "FileStreamSource",
				"file":            testFile,
				"topic":           "connect-test",
			},
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	resp, err := client.GetConnector(ConnectorRequest{
		Name: "test-get-connector",
	})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "test-get-connector", resp.Name)
}

func TestGetAllConnectors(t *testing.T) {
	client := NewClient(hostConnect)
	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-get-all-connectors"},
			Config: map[string]interface{}{
				"connector.class": "FileStreamSource",
				"file":            testFile,
				"topic":           "connect-test",
			},
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	resp, err := client.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Connectors, "test-get-all-connectors")
}

func TestUpdateConnector(t *testing.T) {
	name := "test-update-connectors"
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"tasks.max":       "1",
		"file":            testFile,
		"topic":           "connect-test",
		"test":            "success",
	}

	client := NewClient(hostConnect)
	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: name},
			Config:           config,
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	config["test"] = "success"
	resp, err := client.UpdateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: name},
			Config:           config,
		},
		true,
	)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Config["test"])
}

func TestUpdateConnector_NoCreate(t *testing.T) {
	name := "test-update-connectors-nocreate"
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"tasks.max":       "1",
		"file":            testFile,
		"topic":           "connect-test",
		"test":            "success",
	}

	client := NewClient(hostConnect)
	resp, err := client.UpdateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: name},
			Config:           config,
		},
		true,
	)

	assert.Nil(t, err)
	assert.Equal(t, "success", resp.Config["test"])

	// use IsUpToDate to check sync worked (force get actual config for server rather than what was returned on update call)
	isUpToDate, err := client.IsUpToDate(name, config)
	assert.Nil(t, err)
	assert.True(t, isUpToDate)
}

func TestDeleteConnector(t *testing.T) {
	client := NewClient(hostConnect)
	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-delete-connectors"},
			Config: map[string]interface{}{
				"connector.class": "FileStreamSource",
				"file":            testFile,
				"topic":           "connect-test",
			},
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	resp, err := client.DeleteConnector(ConnectorRequest{Name: "test-delete-connectors"}, true)

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.Code)

	respget, err := client.GetConnector(ConnectorRequest{Name: "test-delete-connectors"})

	assert.Equal(t, 404, respget.Code)
}

func TestGetConnectorConfig(t *testing.T) {
	client := NewClient(hostConnect)
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"tasks.max":       "1",
		"file":            testFile,
		"topic":           "connect-test",
	}

	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-get-connector-config"},
			Config:           config,
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	resp, err := client.GetConnectorConfig(ConnectorRequest{Name: "test-get-connector-config"})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)

	config["name"] = "test-get-connector-config"
	assert.Equal(t, config, resp.Config)
}

func TestIsUpToDate(t *testing.T) {
	client := NewClient(hostConnect)
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"tasks.max":       "1",
		"file":            testFile,
		"topic":           "connect-test",
	}

	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-uptodate-connector-config"},
			Config:           config,
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	uptodate, err := client.IsUpToDate("test-uptodate-connector-config", config)
	assert.Nil(t, err)
	assert.True(t, uptodate)

	config["tasks.max"] = 1
	uptodate, err = client.IsUpToDate("test-uptodate-connector-config", config)
	assert.Nil(t, err)
	assert.True(t, uptodate)

	config["newparameter"] = "test"
	uptodate, err = client.IsUpToDate("test-uptodate-connector-config", config)
	assert.Nil(t, err)
	assert.False(t, uptodate)

}

func TestGetConnectorStatus(t *testing.T) {
	client := NewClient(hostConnect)
	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-get-connector-status"},
			Config: map[string]interface{}{
				"connector.class": "FileStreamSource",
				"file":            testFile,
				"topic":           "connect-test",
			},
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	resp, err := client.GetConnectorStatus(ConnectorRequest{Name: "test-get-connector-status"})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)

	assert.Equal(t, "RUNNING", resp.ConnectorStatus["state"])
}

func TestRestartConnector(t *testing.T) {
	client := NewClient(hostConnect)
	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-restart-connector"},
			Config: map[string]interface{}{
				"connector.class": "FileStreamSource",
				"file":            testFile,
				"topic":           "connect-test",
			},
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	resp, err := client.RestartConnector(ConnectorRequest{Name: "test-restart-connector"})

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.Code)
}

func TestPauseAndResumeConnector(t *testing.T) {
	client := NewClient(hostConnect)
	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-pause-and-resume-connector"},
			Config: map[string]interface{}{
				"connector.class": "FileStreamSource",
				"file":            testFile,
				"topic":           "connect-test",
			},
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creating test connector: %s", err.Error()))
		return
	}

	// First pause connector
	respPause, err := client.PauseConnector(ConnectorRequest{Name: "test-pause-and-resume-connector"}, true)
	assert.Nil(t, err)
	assert.Equal(t, 202, respPause.Code)

	statusResp, err := client.GetConnectorStatus(ConnectorRequest{Name: "test-pause-and-resume-connector"})
	assert.Nil(t, err)
	assert.Equal(t, 200, statusResp.Code)
	assert.Equal(t, "PAUSED", statusResp.ConnectorStatus["state"])

	// Then resume connector
	respResume, err := client.ResumeConnector(ConnectorRequest{Name: "test-pause-and-resume-connector"}, true)
	assert.Nil(t, err)
	assert.Equal(t, 202, respResume.Code)

	statusResp, err = client.GetConnectorStatus(ConnectorRequest{Name: "test-pause-and-resume-connector"})
	assert.Nil(t, err)
	assert.Equal(t, 200, statusResp.Code)
	assert.Equal(t, "RUNNING", statusResp.ConnectorStatus["state"])
}

func TestRestartTask(t *testing.T) {
	client := NewClient(hostConnect)
	_, err := client.CreateConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: "test-restart-task"},
			Config: map[string]interface{}{
				"connector.class": "FileStreamSource",
				"file":            testFile,
				"topic":           "connect-test",
			},
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creating test connector: %s", err.Error()))
		return
	}

	resp, err := client.RestartTask(TaskRequest{Connector: "test-restart-task", TaskID: 0})

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.Code)
}

func TestDeployConnector(t *testing.T) {
	name := "test-deploy-connectors"
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"file":            testFile,
		"topic":           "connect-test",
		"test":            "success",
	}

	client := NewClient(hostConnect)
	err := client.DeployConnector(
		CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{Name: name},
			Config:           config,
		},
	)

	assert.Nil(t, err)

	// use IsUpToDate to check sync worked (force get actual config for server rather than what was returned on update call)
	isUpToDate, err := client.IsUpToDate(name, config)
	assert.Nil(t, err)
	assert.True(t, isUpToDate)
}

func TestDeployMultipleConnectors(t *testing.T) {
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"file":            testFile,
		"topic":           "connect-test",
	}

	req := []CreateConnectorRequest{
		{
			ConnectorRequest: ConnectorRequest{Name: "test-deploy-multiple-connectors-1"},
			Config:           config,
		},
		{
			ConnectorRequest: ConnectorRequest{Name: "test-deploy-multiple-connectors-2"},
			Config:           config,
		},
		{
			ConnectorRequest: ConnectorRequest{Name: "test-deploy-multiple-connectors-3"},
			Config:           config,
		},
		{
			ConnectorRequest: ConnectorRequest{Name: "test-deploy-multiple-connectors-4"},
			Config:           config,
		},
	}

	client := NewClient(hostConnect)
	err := client.DeployMultipleConnector(req)

	assert.Nil(t, err)

	// use IsUpToDate to check sync worked (force get actual config for server rather than what was returned on update call)
	{
		isUpToDate, err := client.IsUpToDate("test-deploy-multiple-connectors-1", config)
		assert.Nil(t, err)
		assert.True(t, isUpToDate)
	}
	{
		isUpToDate, err := client.IsUpToDate("test-deploy-multiple-connectors-2", config)
		assert.Nil(t, err)
		assert.True(t, isUpToDate)
	}
	{
		isUpToDate, err := client.IsUpToDate("test-deploy-multiple-connectors-3", config)
		assert.Nil(t, err)
		assert.True(t, isUpToDate)
	}
	{
		isUpToDate, err := client.IsUpToDate("test-deploy-multiple-connectors-4", config)
		assert.Nil(t, err)
		assert.True(t, isUpToDate)
	}
}
