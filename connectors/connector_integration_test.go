package connectors

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var testFile = "/etc/kafka-connect/kafka-connect.properties"

func TestHealtz(t *testing.T) {
	resp, err := http.Get("http://localhost:8083")
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestCreateConnector(t *testing.T) {
	client := NewClient("localhost", 8083 , false )
	resp, err := client.CreateConnector(	CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{ Name: "test-create-connector" },
		Config: map[string]string{
			"connector.class": "FileStreamSource",
			"tasks.max": "1",
			"file": testFile,
			"topic": "connect-test",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, 201, resp.Code)
}

func TestGetConnector(t *testing.T) {
	client := NewClient("localhost", 8083 , false )
	client.CreateConnector(	CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{ Name: "test-get-connector" },
		Config: map[string]string{
			"connector.class": "FileStreamSource",
			"tasks.max": "1",
			"file": testFile,
			"topic": "connect-test",
		},
	})

	resp, err := client.GetConnector(ConnectorRequest{
		Name:"test-get-connector",
	})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "test-get-connector", resp.Name)
}


func TestGetAllConnectors(t *testing.T) {
	client := NewClient("localhost", 8083 , false )
	client.CreateConnector(	CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{ Name: "test-get-all-connectors" },
		Config: map[string]string{
			"connector.class": "FileStreamSource",
			"tasks.max": "1",
			"file": testFile,
			"topic": "connect-test",
		},
	})

	resp, err := client.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Contains(t, resp.Connectors, "test-get-all-connectors")
}

func TestUpdateConnector(t *testing.T) {
	client := NewClient("localhost", 8083 , false )
	client.CreateConnector(	CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{ Name: "test-update-connectors" },
		Config: map[string]string{
			"connector.class": "FileStreamSource",
			"tasks.max": "1",
			"file": testFile,
			"topic": "connect-test",
		},
	})

	resp, err := client.UpdateConnector(UpdateConnectorRequest{
		ConnectorRequest: ConnectorRequest{Name:"test-update-connectors" },
		Config: map[string]string{
			"connector.class": "FileStreamSource",
			"tasks.max": "1",
			"file": testFile,
			"topic": "connect-test",
			"test": "success",
		},
	})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)
	assert.Equal(t, "success", resp.Config["test"])
}

func TestDeleteConnector(t *testing.T) {
	client := NewClient("localhost", 8083 , false )
	client.CreateConnector(	CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{ Name: "test-delete-connectors" },
		Config: map[string]string{
			"connector.class": "FileStreamSource",
			"tasks.max": "1",
			"file": testFile,
			"topic": "connect-test",
		},
	})

	resp, err := client.DeleteConnector(ConnectorRequest{Name: "test-delete-connectors"})

	assert.Nil(t, err)
	assert.Equal(t, 204, resp.Code)

	respget, err := client.GetConnector(ConnectorRequest{Name: "test-delete-connectors"})

	assert.Equal(t, 404, respget.Code)
}

func TestGetConnectorConfig(t *testing.T) {
	client := NewClient("localhost", 8083 , false )
	config := map[string]string{
		"connector.class": "FileStreamSource",
		"tasks.max": "1",
		"file": testFile,
		"topic": "connect-test",
	}

	client.CreateConnector(	CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{ Name: "test-get-connector-config" },
		Config: config,
	})

	resp, err := client.GetConnectorConfig(ConnectorRequest{Name: "test-get-connector-config"})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)

	config["name"]= "test-get-connector-config"
	assert.Equal(t, config, resp.Config)
}

func TestGetConnectorStatus(t *testing.T) {
	client := NewClient("localhost", 8083 , false )
	client.CreateConnector(	CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{ Name: "test-get-connector-status" },
		Config: map[string]string{
			"connector.class": "FileStreamSource",
			"tasks.max": "1",
			"file": testFile,
			"topic": "connect-test",
		},
	})

	resp, err := client.GetConnectorStatus(ConnectorRequest{Name: "test-get-connector-status"})

	assert.Nil(t, err)
	assert.Equal(t, 200, resp.Code)

	assert.Equal(t, "RUNNING", resp.ConnectorStatus["state"])
}

//These last 3 feature are officially asynchronous and shall not be tested as is

//func TestRestartConnector(t *testing.T) {
//	client := NewClient("localhost", 8083 , false )
//	client.CreateConnector(	CreateConnectorRequest{
//		ConnectorRequest: ConnectorRequest{ Name: "test-restart-connector" },
//		Config: map[string]string{
//			"connector.class": "FileStreamSource",
//			"tasks.max": "1",
//			"file": testFile,
//			"topic": "connect-test",
//		},
//	})
//
//	resp, err := client.RestartConnector(ConnectorRequest{Name: "test-restart-connector"})
//
//	assert.Nil(t, err)
//	assert.Equal(t, 204, resp.Code)
//}
//
//
//
//func TestPauseConnector(t *testing.T) {
//	client := NewClient("localhost", 8083 , false )
//	client.CreateConnector(	CreateConnectorRequest{
//		ConnectorRequest: ConnectorRequest{ Name: "test-pause-connector" },
//		Config: map[string]string{
//			"connector.class": "FileStreamSource",
//			"tasks.max": "1",
//			"file": testFile,
//			"topic": "connect-test",
//		},
//	})
//
//	resp, err := client.RestartConnector(ConnectorRequest{Name: "test-pause-connector"})
//
//	assert.Nil(t, err)
//	assert.Equal(t, 204, resp.Code)
//}
//
//func TestResumeConnector(t *testing.T) {
//	client := NewClient("localhost", 8083 , false )
//	client.CreateConnector(	CreateConnectorRequest{
//		ConnectorRequest: ConnectorRequest{ Name: "test-resume-connector" },
//		Config: map[string]string{
//			"connector.class": "FileStreamSource",
//			"tasks.max": "1",
//			"file": testFile,
//			"topic": "connect-test",
//		},
//	})
//
//	resp, err := client.RestartConnector(ConnectorRequest{Name: "test-resume-connector"})
//
//	assert.Nil(t, err)
//	assert.Equal(t, 204, resp.Code)
//}
//

