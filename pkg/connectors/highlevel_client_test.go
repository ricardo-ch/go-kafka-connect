package connectors

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"

	"github.com/sirupsen/logrus"
	logrusTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_IsUpToDate_Should_Be_True(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":   "test1",
		"param1": 2,
		"param2": "abc",
		"param3": "3",
	}

	configLocal := map[string]interface{}{
		"param1": 2,
		"param2": "abc",
		"param3": 3,
	}

	mockBaseClient := NewMockBaseClient(t)
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)

	client := newHighLevelClient(mockBaseClient, 0)

	isUpToDate, err := client.IsUpToDate("test1", configLocal)

	assert.NoError(t, err)
	assert.True(t, isUpToDate)
}

func Test_IsUpToDate_Should_Be_False(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":   "test1",
		"param1": 3,
		"param2": "abc",
		"param3": "3",
	}

	configLocal := map[string]interface{}{
		"param1": 2,
		"param2": "abc",
		"param3": 3,
	}

	mockBaseClient := NewMockBaseClient(t)
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)

	client := newHighLevelClient(mockBaseClient, 0)

	isUpToDate, err := client.IsUpToDate("test1", configLocal)

	assert.NoError(t, err)
	assert.False(t, isUpToDate)
}

func Test_tryUntil_When_Success(t *testing.T) {
	result := tryUntil(
		func() bool {
			return true
		},
		100*time.Millisecond,
	)
	assert.True(t, result)
}

func Test_tryUntil_When_Timeout(t *testing.T) {
	result := tryUntil(
		func() bool {
			time.Sleep(200 * time.Millisecond)
			return true
		},
		100*time.Millisecond,
	)
	assert.False(t, result)
}

func Test_DeployConnector_When_Already_Up_To_Date(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":   "test1",
		"param1": 2,
		"param2": "abc",
		"param3": "3",
	}
	configLocal := map[string]interface{}{
		"param1": 2,
		"param2": "abc",
		"param3": 3,
	}

	mockBaseClient := NewMockBaseClient(t)
	mockBaseClient.On("GetConnector", mock.Anything).
		Return(ConnectorResponse{Name: "test1", Config: configOnline}, nil)
	//TODO there shouldn't be a need to make both these call
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)
	// note we don't mock the update part because it should not be called

	client := newHighLevelClient(mockBaseClient, 0)
	err := client.DeployConnector(CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{"test1"},
		Config:           configLocal,
	})

	assert.NoError(t, err)
}

func Test_DeployConnector_Ok(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":   "test1",
		"param1": 2,
	}
	configLocal := map[string]interface{}{
		"param1": 3,
	}

	t.Run("When_Connector_Is_Running_And_Deploy_Expect_Pause", func(t *testing.T) {
		pause := true
		running := true
		callDeployConnector(t, pause, running, configOnline, configLocal)
	})

	t.Run("When_Connector_IsNot_Running_And_Deploy_Expect_Pause", func(t *testing.T) {
		pause := true
		running := false
		callDeployConnector(t, pause, running, configOnline, configLocal)
	})

	t.Run("When_Connector_Is_Running_And_Deploy_Skip_Pause", func(t *testing.T) {
		pause := false
		running := true // running true or false is irrelevant
		callDeployConnector(t, pause, running, configOnline, configLocal)
	})
}

func Test_DeployMultipleConnector_Ok(t *testing.T) {
	mockBaseClient := NewMockBaseClient(t)
	client := newHighLevelClient(mockBaseClient, 2)

	lock := &sync.Mutex{}
	received := map[string]interface{}{}

	connectorRequests := []CreateConnectorRequest{}
	for i := 1; i <= 5; i++ {
		connectorName := fmt.Sprintf("test%d", i)
		request := CreateConnectorRequest{
			ConnectorRequest: ConnectorRequest{connectorName},
		}
		mockBaseClient.On("GetConnector", request.ConnectorRequest).
			Return(ConnectorResponse{Name: connectorName}, nil).Once()
		mockBaseClient.On("GetConnectorConfig", request.ConnectorRequest).
			Return(GetConnectorConfigResponse{}, nil).Once()

		mockBaseClient.On("UpdateConnector", request).Return(ConnectorResponse{}, nil).Run(func(args mock.Arguments) {
			req := args.Get(0).(CreateConnectorRequest)
			lock.Lock()
			defer lock.Unlock()
			received[req.Name] = true
		}).Once()
		mockBaseClient.On("GetConnectorConfig", request.ConnectorRequest).
			Return(GetConnectorConfigResponse{Config: map[string]interface{}{"name": connectorName}}, nil).Once()

		connectorRequests = append(connectorRequests, request)
	}

	err := client.DeployMultipleConnector(connectorRequests)

	assert.Equal(t, map[string]interface{}{"test1": true, "test2": true, "test3": true, "test4": true, "test5": true}, received)
	assert.NoError(t, err)
}

func Test_DeployMultipleConnector_Error(t *testing.T) {
	mockBaseClient := NewMockBaseClient(t)
	client := newHighLevelClient(mockBaseClient, 2)

	mockBaseClient.On("GetConnector", mock.Anything).Times(5).Return(ConnectorResponse{}, errors.New("random error"))

	err := client.DeployMultipleConnector([]CreateConnectorRequest{
		{ConnectorRequest: ConnectorRequest{Name: "test1"}},
		{ConnectorRequest: ConnectorRequest{Name: "test2"}},
		{ConnectorRequest: ConnectorRequest{Name: "test3"}},
		{ConnectorRequest: ConnectorRequest{Name: "test4"}},
		{ConnectorRequest: ConnectorRequest{Name: "test5"}},
	})

	assert.Error(t, err)
}

func newHighLevelClient(mockBaseClient *MockBaseClient, maxParallelRequest int) *highLevelClient {
	logger, _ := logrusTest.NewNullLogger()

	return &highLevelClient{
		client:             mockBaseClient,
		maxParallelRequest: maxParallelRequest,
		logger:             logrus.NewEntry(logger),
	}
}

// callDeployConnector is a helper function to test DeployConnector
// when we ask to pause the connector we expect:
// - get connector and config to perform comparison
// - if running: status call + pause/resume
// - if NOT running: status call only
// - update the connector and compare config until they are identical
func callDeployConnector(t *testing.T, pauseBeforeDeploy, isRunning bool, configOnline, configLocal map[string]interface{}) {
	t.Helper()
	mockBaseClient := NewMockBaseClient(t)

	mockBaseClient.On("GetConnector", mock.Anything).
		Return(ConnectorResponse{Name: "test1", Config: configOnline}, nil)
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil).Once()

	if pauseBeforeDeploy {
		if isRunning {
			mockBaseClient.On("GetConnectorStatus", mock.Anything).
				Return(GetConnectorStatusResponse{EmptyResponse: EmptyResponse{Code: 200}, ConnectorStatus: map[string]string{"state": "RUNNING"}}, nil).Once()

			mockBaseClient.On("PauseConnector", mock.Anything, mock.Anything).
				Return(EmptyResponse{}, nil)
			mockBaseClient.On("GetConnectorStatus", mock.Anything).
				Return(GetConnectorStatusResponse{EmptyResponse: EmptyResponse{Code: 200}, ConnectorStatus: map[string]string{"state": "PAUSED"}}, nil).Once()

			mockBaseClient.On("ResumeConnector", mock.Anything, mock.Anything).
				Return(EmptyResponse{}, nil)
			mockBaseClient.On("GetConnectorStatus", mock.Anything).
				Return(GetConnectorStatusResponse{EmptyResponse: EmptyResponse{Code: 200}, ConnectorStatus: map[string]string{"state": "RUNNING"}}, nil).Once()
		} else {
			mockBaseClient.On("GetConnectorStatus", mock.Anything).
				Return(GetConnectorStatusResponse{EmptyResponse: EmptyResponse{Code: 200}, ConnectorStatus: map[string]string{"state": "PAUSED"}}, nil).Once()
		}
	}

	mockBaseClient.On("UpdateConnector", mock.Anything).
		Return(ConnectorResponse{}, nil)
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: map[string]interface{}{"name": "test1", "param1": 3}}, nil).Once()

	client := newHighLevelClient(mockBaseClient, 0)
	client.SetPauseBeforeDeploy(pauseBeforeDeploy)

	err := client.DeployConnector(CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{"test1"},
		Config:           configLocal,
	})

	assert.NoError(t, err)
}
