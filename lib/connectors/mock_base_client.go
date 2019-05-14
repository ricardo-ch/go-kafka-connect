// Code generated by mockery v1.0.0. DO NOT EDIT.

// NOTE: run 'make update-mocks' from this project top folder to update this file and generate new ones.

package connectors

import mock "github.com/stretchr/testify/mock"

// MockBaseClient is an autogenerated mock type for the BaseClient type
type MockBaseClient struct {
	mock.Mock
}

// CreateConnector provides a mock function with given fields: req
func (_m *MockBaseClient) CreateConnector(req CreateConnectorRequest) (ConnectorResponse, error) {
	ret := _m.Called(req)

	var r0 ConnectorResponse
	if rf, ok := ret.Get(0).(func(CreateConnectorRequest) ConnectorResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(ConnectorResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(CreateConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteConnector provides a mock function with given fields: req
func (_m *MockBaseClient) DeleteConnector(req ConnectorRequest) (EmptyResponse, error) {
	ret := _m.Called(req)

	var r0 EmptyResponse
	if rf, ok := ret.Get(0).(func(ConnectorRequest) EmptyResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(EmptyResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *MockBaseClient) GetAll() (GetAllConnectorsResponse, error) {
	ret := _m.Called()

	var r0 GetAllConnectorsResponse
	if rf, ok := ret.Get(0).(func() GetAllConnectorsResponse); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(GetAllConnectorsResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllTasks provides a mock function with given fields: req
func (_m *MockBaseClient) GetAllTasks(req ConnectorRequest) (GetAllTasksResponse, error) {
	ret := _m.Called(req)

	var r0 GetAllTasksResponse
	if rf, ok := ret.Get(0).(func(ConnectorRequest) GetAllTasksResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(GetAllTasksResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConnector provides a mock function with given fields: req
func (_m *MockBaseClient) GetConnector(req ConnectorRequest) (ConnectorResponse, error) {
	ret := _m.Called(req)

	var r0 ConnectorResponse
	if rf, ok := ret.Get(0).(func(ConnectorRequest) ConnectorResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(ConnectorResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConnectorConfig provides a mock function with given fields: req
func (_m *MockBaseClient) GetConnectorConfig(req ConnectorRequest) (GetConnectorConfigResponse, error) {
	ret := _m.Called(req)

	var r0 GetConnectorConfigResponse
	if rf, ok := ret.Get(0).(func(ConnectorRequest) GetConnectorConfigResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(GetConnectorConfigResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetConnectorStatus provides a mock function with given fields: req
func (_m *MockBaseClient) GetConnectorStatus(req ConnectorRequest) (GetConnectorStatusResponse, error) {
	ret := _m.Called(req)

	var r0 GetConnectorStatusResponse
	if rf, ok := ret.Get(0).(func(ConnectorRequest) GetConnectorStatusResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(GetConnectorStatusResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTaskStatus provides a mock function with given fields: req
func (_m *MockBaseClient) GetTaskStatus(req TaskRequest) (TaskStatusResponse, error) {
	ret := _m.Called(req)

	var r0 TaskStatusResponse
	if rf, ok := ret.Get(0).(func(TaskRequest) TaskStatusResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(TaskStatusResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(TaskRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PauseConnector provides a mock function with given fields: req
func (_m *MockBaseClient) PauseConnector(req ConnectorRequest) (EmptyResponse, error) {
	ret := _m.Called(req)

	var r0 EmptyResponse
	if rf, ok := ret.Get(0).(func(ConnectorRequest) EmptyResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(EmptyResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RestartConnector provides a mock function with given fields: req
func (_m *MockBaseClient) RestartConnector(req ConnectorRequest) (EmptyResponse, error) {
	ret := _m.Called(req)

	var r0 EmptyResponse
	if rf, ok := ret.Get(0).(func(ConnectorRequest) EmptyResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(EmptyResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RestartTask provides a mock function with given fields: req
func (_m *MockBaseClient) RestartTask(req TaskRequest) (EmptyResponse, error) {
	ret := _m.Called(req)

	var r0 EmptyResponse
	if rf, ok := ret.Get(0).(func(TaskRequest) EmptyResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(EmptyResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(TaskRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResumeConnector provides a mock function with given fields: req
func (_m *MockBaseClient) ResumeConnector(req ConnectorRequest) (EmptyResponse, error) {
	ret := _m.Called(req)

	var r0 EmptyResponse
	if rf, ok := ret.Get(0).(func(ConnectorRequest) EmptyResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(EmptyResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetDebug provides a mock function with given fields:
func (_m *MockBaseClient) SetDebug() {
	_m.Called()
}

// SetInsecureSSL provides a mock function with given fields:
func (_m *MockBaseClient) SetInsecureSSL() {
	_m.Called()
}

// SetClientCertificates provides a mock function with given fields:
func (_m *MockBaseClient) SetClientCertificates(certFile string, keyFile string) {
	_m.Called()
}

// UpdateConnector provides a mock function with given fields: req
func (_m *MockBaseClient) UpdateConnector(req CreateConnectorRequest) (ConnectorResponse, error) {
	ret := _m.Called(req)

	var r0 ConnectorResponse
	if rf, ok := ret.Get(0).(func(CreateConnectorRequest) ConnectorResponse); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(ConnectorResponse)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(CreateConnectorRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
