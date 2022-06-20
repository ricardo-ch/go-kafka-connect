//go:build integration

package connectors_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/heetch/go-kafka-connect/v4/pkg/connectors"
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
	client := connectors.NewClient(hostConnect)

	resp, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: "test-create-connector" + uuid.New().String()},
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
	client := connectors.NewClient(hostConnect)
	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: "not-a-valid-connector" + uuid.New().String()},
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
	connectorName := "test-get-connector" + uuid.New().String()

	client := connectors.NewClient(hostConnect)
	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
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

	resp, err := client.GetConnector(connectors.ConnectorRequest{
		Name: connectorName,
	})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, connectorName, resp.Name)
}

func TestGetAllConnectors(t *testing.T) {
	connectorName := "test-get-all-connectors" + uuid.New().String()
	client := connectors.NewClient(hostConnect)
	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
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
	assert.Contains(t, resp.Connectors, connectorName)
}

func TestUpdateConnector(t *testing.T) {
	connectorName := "test-update-connectors" + uuid.New().String()
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"tasks.max":       "1",
		"file":            testFile,
		"topic":           "connect-test",
		"test":            "success",
	}

	client := connectors.NewClient(hostConnect)
	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
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
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
			Config:           config,
		},
		true,
	)

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Config["test"])
}

func TestUpdateConnector_NoCreate(t *testing.T) {
	connectorName := "test-update-connectors-nocreate" + uuid.New().String()
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"tasks.max":       "1",
		"file":            testFile,
		"topic":           "connect-test",
		"test":            "success",
	}

	client := connectors.NewClient(hostConnect)
	resp, err := client.UpdateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
			Config:           config,
		},
		true,
	)

	assert.Nil(t, err)
	assert.Equal(t, "success", resp.Config["test"])

	// use IsUpToDate to check sync worked (force get actual config for server rather than what was returned on update call)
	isUpToDate, err := client.IsUpToDate(connectorName, config)
	assert.Nil(t, err)
	assert.True(t, isUpToDate)
}

func TestDeleteConnector(t *testing.T) {
	connectorName := "test-delete-connectors" + uuid.New().String()
	client := connectors.NewClient(hostConnect)
	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
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

	resp, err := client.DeleteConnector(connectors.ConnectorRequest{Name: connectorName}, true)

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.Code)

	respget, err := client.GetConnector(connectors.ConnectorRequest{Name: connectorName})

	assert.Equal(t, 404, respget.Code)
}

func TestGetConnectorConfig(t *testing.T) {
	connectorName := "test-get-connector-config" + uuid.New().String()
	client := connectors.NewClient(hostConnect)
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"tasks.max":       "1",
		"file":            testFile,
		"topic":           "connect-test",
	}

	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
			Config:           config,
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	resp, err := client.GetConnectorConfig(connectors.ConnectorRequest{Name: connectorName})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)

	config["name"] = connectorName
	assert.Equal(t, config, resp.Config)
}

func TestIsUpToDate(t *testing.T) {
	connectorName := "test-uptodate-connector-config" + uuid.New().String()
	client := connectors.NewClient(hostConnect)
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"tasks.max":       "1",
		"file":            testFile,
		"topic":           "connect-test",
	}

	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
			Config:           config,
		},
		true,
	)
	if err != nil {
		assert.Fail(t, fmt.Sprintf("error while creaating test connector: %s", err.Error()))
		return
	}

	uptodate, err := client.IsUpToDate(connectorName, config)
	assert.Nil(t, err)
	assert.True(t, uptodate)

	config["tasks.max"] = 1
	uptodate, err = client.IsUpToDate(connectorName, config)
	assert.Nil(t, err)
	assert.True(t, uptodate)

	config["newparameter"] = "test"
	uptodate, err = client.IsUpToDate(connectorName, config)
	assert.Nil(t, err)
	assert.False(t, uptodate)

}

func TestGetConnectorStatus(t *testing.T) {
	connectorName := "test-get-connector-status" + uuid.New().String()
	client := connectors.NewClient(hostConnect)
	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
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

	resp, err := client.GetConnectorStatus(connectors.ConnectorRequest{Name: connectorName})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)

	assert.Equal(t, "RUNNING", resp.ConnectorStatus["state"])
}

func TestRestartConnector(t *testing.T) {
	connectorName := "test-restart-connector" + uuid.New().String()
	client := connectors.NewClient(hostConnect)
	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
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

	resp, err := client.RestartConnector(connectors.ConnectorRequest{Name: connectorName})

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.Code)
}

func TestPauseAndResumeConnector(t *testing.T) {
	connectorName := "test-pause-and-resume-connector" + uuid.New().String()
	client := connectors.NewClient(hostConnect)
	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
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
	respPause, err := client.PauseConnector(connectors.ConnectorRequest{Name: connectorName}, true)
	assert.Nil(t, err)
	assert.Equal(t, 202, respPause.Code)

	statusResp, err := client.GetConnectorStatus(connectors.ConnectorRequest{Name: connectorName})
	assert.Nil(t, err)
	assert.Equal(t, 200, statusResp.Code)
	assert.Equal(t, "PAUSED", statusResp.ConnectorStatus["state"])

	// Then resume connector
	respResume, err := client.ResumeConnector(connectors.ConnectorRequest{Name: connectorName}, true)
	assert.Nil(t, err)
	assert.Equal(t, 202, respResume.Code)

	statusResp, err = client.GetConnectorStatus(connectors.ConnectorRequest{Name: connectorName})
	assert.Nil(t, err)
	assert.Equal(t, 200, statusResp.Code)
	assert.Equal(t, "RUNNING", statusResp.ConnectorStatus["state"])
}

func TestRestartTask(t *testing.T) {
	connectorName := "test-restart-task" + uuid.New().String()
	client := connectors.NewClient(hostConnect)
	_, err := client.CreateConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
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

	resp, err := client.RestartTask(connectors.TaskRequest{Connector: connectorName, TaskID: 0})

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.Code)
}

func TestDeployConnector(t *testing.T) {
	connectorName := "test-deploy-connectors" + uuid.New().String()
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"file":            testFile,
		"topic":           "connect-test",
		"test":            "success",
	}

	client := connectors.NewClient(hostConnect)
	err := client.DeployConnector(
		connectors.CreateConnectorRequest{
			ConnectorRequest: connectors.ConnectorRequest{Name: connectorName},
			Config:           config,
		},
	)

	assert.Nil(t, err)

	// use IsUpToDate to check sync worked (force get actual config for server rather than what was returned on update call)
	isUpToDate, err := client.IsUpToDate(connectorName, config)
	assert.Nil(t, err)
	assert.True(t, isUpToDate)
}

func TestDeployMultipleConnectors(t *testing.T) {
	config := map[string]interface{}{
		"connector.class": "FileStreamSource",
		"file":            testFile,
		"topic":           "connect-test",
	}

	req := []connectors.CreateConnectorRequest{
		{
			ConnectorRequest: connectors.ConnectorRequest{Name: "test-deploy-multiple-connectors-1"},
			Config:           config,
		},
		{
			ConnectorRequest: connectors.ConnectorRequest{Name: "test-deploy-multiple-connectors-2"},
			Config:           config,
		},
		{
			ConnectorRequest: connectors.ConnectorRequest{Name: "test-deploy-multiple-connectors-3"},
			Config:           config,
		},
		{
			ConnectorRequest: connectors.ConnectorRequest{Name: "test-deploy-multiple-connectors-4"},
			Config:           config,
		},
	}

	client := connectors.NewClient(hostConnect)
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
