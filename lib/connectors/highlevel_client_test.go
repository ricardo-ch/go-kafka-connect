//go:build !integration

package connectors

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"

	"bou.ke/monkey"
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

	mockBaseClient := &MockBaseClient{}
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)

	client := &highLevelClient{client: mockBaseClient}

	isUpToDate, err := client.IsUpToDate("test1", configLocal)

	assert.NoError(t, err)
	assert.True(t, isUpToDate)
}

func Test_IsUpToDate_Ignores_Internal_Keys(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":                    "test1",
		"param1":                  2,
		"param2":                  "abc",
		"param3":                  "3",
		"__internal.offset.topic": "connect-offsets",
		"__internal.task.id":      "0",
	}

	configLocal := map[string]interface{}{
		"param1": 2,
		"param2": "abc",
		"param3": 3,
	}

	mockBaseClient := &MockBaseClient{}
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)

	client := &highLevelClient{client: mockBaseClient}
	client.SetDisableInternalConfig(true)

	isUpToDate, err := client.IsUpToDate("test1", configLocal)

	assert.NoError(t, err)
	assert.True(t, isUpToDate)
}

// By default "__internal." keys are NOT stripped: deployed-only internal keys make the
// config differ and report not up to date.
func Test_IsUpToDate_Keeps_Internal_Keys_By_Default(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":                    "test1",
		"param1":                  2,
		"param2":                  "abc",
		"param3":                  "3",
		"__internal.offset.topic": "connect-offsets",
		"__internal.task.id":      "0",
	}

	configLocal := map[string]interface{}{
		"param1": 2,
		"param2": "abc",
		"param3": 3,
	}

	mockBaseClient := &MockBaseClient{}
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)

	client := &highLevelClient{client: mockBaseClient}

	isUpToDate, err := client.IsUpToDate("test1", configLocal)

	assert.NoError(t, err)
	assert.False(t, isUpToDate)
}

// Internal keys present on BOTH sides with different values must still be ignored,
// proving the value is genuinely dropped rather than just count-balanced.
func Test_IsUpToDate_Ignores_Internal_Keys_With_Differing_Values(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":             "test1",
		"param1":           2,
		"param2":           "abc",
		"__internal.state": "deployed-value",
	}

	configLocal := map[string]interface{}{
		"param1":           2,
		"param2":           "abc",
		"__internal.state": "local-value",
	}

	mockBaseClient := &MockBaseClient{}
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)

	client := &highLevelClient{client: mockBaseClient}
	client.SetDisableInternalConfig(true)

	isUpToDate, err := client.IsUpToDate("test1", configLocal)

	assert.NoError(t, err)
	assert.True(t, isUpToDate)
}

// A mismatch on a NON-internal key must still be detected even when internal keys are
// present, proving stripping does not mask a real diff.
func Test_IsUpToDate_Detects_Real_Diff_When_Internal_Keys_Present(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":             "test1",
		"param1":           3,
		"param2":           "abc",
		"__internal.state": "deployed-value",
	}

	configLocal := map[string]interface{}{
		"param1": 2,
		"param2": "abc",
	}

	mockBaseClient := &MockBaseClient{}
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)

	client := &highLevelClient{client: mockBaseClient}
	client.SetDisableInternalConfig(true)

	isUpToDate, err := client.IsUpToDate("test1", configLocal)

	assert.NoError(t, err)
	assert.False(t, isUpToDate)
}

// Keys that merely contain "__internal." but are not prefixed by it must NOT be stripped.
func Test_IsUpToDate_Does_Not_Strip_Substring_Keys(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":                   "test1",
		"param1":                 2,
		"transforms.__internal.": "should-not-be-stripped",
	}

	configLocal := map[string]interface{}{
		"param1": 2,
	}

	mockBaseClient := &MockBaseClient{}
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)

	client := &highLevelClient{client: mockBaseClient}
	client.SetDisableInternalConfig(true)

	isUpToDate, err := client.IsUpToDate("test1", configLocal)

	assert.NoError(t, err)
	// The substring key survives, so the deployed config has an extra key the local one
	// lacks -> not up to date.
	assert.False(t, isUpToDate)
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

	mockBaseClient := &MockBaseClient{}
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)

	client := &highLevelClient{client: mockBaseClient}

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

	mockBaseClient := &MockBaseClient{}
	mockBaseClient.On("GetConnector", mock.Anything).
		Return(ConnectorResponse{Name: "test1", Config: configOnline}, nil)
	//TODO there shouldn't be a need to make both these call
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil)
	// note we don't mock the update part because it should not be called

	client := &highLevelClient{client: mockBaseClient}
	err := client.DeployConnector(CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{"test1"},
		Config:           configLocal,
	})

	assert.NoError(t, err)
	mockBaseClient.AssertExpectations(t)
}

func Test_DeployConnector_Ok(t *testing.T) {
	configOnline := map[string]interface{}{
		"name":   "test1",
		"param1": 2,
	}
	configLocal := map[string]interface{}{
		"param1": 3,
	}

	// expected steps:
	// - get connector info
	// - compare online config and local config
	// - pause online connector
	// - loop get connector status until it is paused
	// - update connector
	// - loop get connector config until it match deployed
	// - resume connector
	// - loop get connector status until it is running

	mockBaseClient := &MockBaseClient{}
	mockBaseClient.On("GetConnector", mock.Anything).
		Return(ConnectorResponse{Name: "test1", Config: configOnline}, nil)
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: configOnline}, nil).Once()
	mockBaseClient.On("UpdateConnector", mock.Anything).
		Return(ConnectorResponse{}, nil)
	mockBaseClient.On("GetConnectorConfig", mock.Anything).
		Return(GetConnectorConfigResponse{Config: map[string]interface{}{"name": "test1", "param1": 3}}, nil).Once()

	client := &highLevelClient{client: mockBaseClient}
	err := client.DeployConnector(CreateConnectorRequest{
		ConnectorRequest: ConnectorRequest{"test1"},
		Config:           configLocal,
	})

	assert.NoError(t, err)
	mockBaseClient.AssertExpectations(t)
}

func Test_DeployMultipleConnector_Ok(t *testing.T) {
	client := &highLevelClient{client: &MockBaseClient{}, maxParallelRequest: 2}

	lock := &sync.Mutex{}
	received := map[string]interface{}{}

	// Don't want to mock every baseClient call, so I am going the lazy way.
	patch := monkey.PatchInstanceMethod(reflect.TypeOf(client), "DeployConnector", func(_ *highLevelClient, req CreateConnectorRequest) (err error) {
		lock.Lock()
		defer lock.Unlock()
		received[req.Name] = true
		return nil
	})
	defer patch.Restore()

	err := client.DeployMultipleConnector([]CreateConnectorRequest{
		{ConnectorRequest: ConnectorRequest{Name: "test1"}},
		{ConnectorRequest: ConnectorRequest{Name: "test2"}},
		{ConnectorRequest: ConnectorRequest{Name: "test3"}},
		{ConnectorRequest: ConnectorRequest{Name: "test4"}},
		{ConnectorRequest: ConnectorRequest{Name: "test5"}},
	})

	assert.Equal(t, map[string]interface{}{"test1": true, "test2": true, "test3": true, "test4": true, "test5": true}, received)
	assert.NoError(t, err)
}

func Test_DeployMultipleConnector_Error(t *testing.T) {
	client := &highLevelClient{client: &MockBaseClient{}, maxParallelRequest: 2}

	// Don't want to mock every baseClient call, so I am going the lazy way.
	patch := monkey.PatchInstanceMethod(reflect.TypeOf(client), "DeployConnector", func(_ *highLevelClient, req CreateConnectorRequest) (err error) {
		return errors.New("random error")
	})
	defer patch.Restore()

	err := client.DeployMultipleConnector([]CreateConnectorRequest{
		{ConnectorRequest: ConnectorRequest{Name: "test1"}},
		{ConnectorRequest: ConnectorRequest{Name: "test2"}},
		{ConnectorRequest: ConnectorRequest{Name: "test3"}},
		{ConnectorRequest: ConnectorRequest{Name: "test4"}},
		{ConnectorRequest: ConnectorRequest{Name: "test5"}},
	})

	assert.Error(t, err)
}
