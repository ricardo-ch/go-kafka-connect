package cmd

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/winniehuang-ap/kafka-connect/v4/lib/connectors"
	"github.com/stretchr/testify/assert"
)

func Test_getCreateCmdConfig_Invalid_Path(t *testing.T) {
	configs, err := getConfigFromFolder("missingFolder")

	assert.NotNil(t, err)
	assert.Zero(t, configs)
}

func Test_getConfigFromFolder_OK(t *testing.T) {
	expected := []connectors.CreateConnectorRequest{
		{
			ConnectorRequest: connectors.ConnectorRequest{
				Name: "test-connector",
			},
			Config: map[string]interface{}{
				"config-entry": "some-value",
			},
		},
		{
			ConnectorRequest: connectors.ConnectorRequest{
				Name: "test-connector-2",
			},
			Config: map[string]interface{}{
				"config-entry-2": "some-value-2",
			},
		},
	}

	_ = os.Mkdir("test", os.ModePerm)
	defer os.RemoveAll("test")

	j1, _ := json.MarshalIndent(expected[0], "", "")
	j2, _ := json.MarshalIndent(expected[1], "", "")

	_ = ioutil.WriteFile("test/test1.json", j1, os.ModePerm)
	_ = ioutil.WriteFile("test/test2.json", j2, os.ModePerm)

	actual, err := getConfigFromFolder("test")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)
}

func Test_getConfigFromFolder_Invalid_JSON(t *testing.T) {
	unsupported := `{"some-key": "some-value"}`

	_ = os.Mkdir("test", os.ModePerm)
	defer os.RemoveAll("test")

	j1, _ := json.MarshalIndent(unsupported, "", "")

	_ = ioutil.WriteFile("test/test1.json", j1, os.ModePerm)

	actual, err := getConfigFromFolder("test")

	assert.NotNil(t, err)
	assert.Equal(t, []connectors.CreateConnectorRequest{}, actual)
}

func Test_getConfigFromFolder_Missing_Folder(t *testing.T) {
	configs, err := getConfigFromFolder("test")

	assert.NotNil(t, err)
	assert.Zero(t, configs)
}

func Test_getConfigFromFolder_Not_A_Folder(t *testing.T) {
	// Create a file instead of a folder
	_, _ = os.Create("test")
	defer os.Remove("test")

	configs, err := getConfigFromFolder("test")

	assert.NotNil(t, err)
	assert.Zero(t, configs)
}

func Test_getConfigFromFile_OK(t *testing.T) {
	expected := connectors.CreateConnectorRequest{
		ConnectorRequest: connectors.ConnectorRequest{
			Name: "test-connector",
		},
		Config: map[string]interface{}{
			"config-entry": "some-value",
		},
	}

	j, _ := json.MarshalIndent(expected, "", "")

	_ = ioutil.WriteFile("test.json", j, 0644)
	defer os.Remove("test.json")

	actual, err := getConfigFromFile("test.json")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual)

}

func Test_getConfigFromFile_Invalid_JSON(t *testing.T) {
	unsupported := `{"some-key": "some-value"}`

	j, _ := json.MarshalIndent(unsupported, "", "")

	_ = ioutil.WriteFile("test.json", j, os.ModePerm)
	defer os.Remove("test.json")

	actual, err := getConfigFromFile("test.json")

	assert.NotNil(t, err)
	assert.Zero(t, actual)
}
