// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/heetch/go-kafka-connect/v4/pkg/connectors"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new connector",
	RunE:  RunECreate,
}

//RunECreate ...
func RunECreate(cmd *cobra.Command, args []string) error {
	configs, err := getCreateCmdConfig(cmd, expandEnv)
	if err != nil {
		return err
	}

	//TODO was not expecting I would have to update CreateConnector when adding multiple file deployment feature
	// will have to add properly later
	for _, config := range configs {
		resp, err := getClient().CreateConnector(config, sync)
		printResponse(resp)
		if err != nil {
			return err
		}
	}
	return nil
}

func getCreateCmdConfig(cmd *cobra.Command, expandEnv bool) ([]connectors.CreateConnectorRequest, error) {
	var configs []connectors.CreateConnectorRequest

	if cmd.Flag("path").Changed {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			return nil, errors.Wrapf(err, "error while trying to find input or folder: %v", filePath)
		}
		if fileInfo.IsDir() {
			configs, err = getConfigFromFolder(filePath, expandEnv)
			if err != nil {
				return nil, err
			}
		} else {
			config, err := getConfigFromFile(filePath, expandEnv)
			if err != nil {
				return nil, err
			}
			configs = append(configs, config)
		}

	} else if cmd.Flag("string").Changed {
		config, err := getConfigFromString(configString, expandEnv)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	} else {
		return nil, errors.New("neither path nor string was supplied")
	}
	return configs, nil
}

func getConfigFromFolder(folderPath string, expandEnv bool) ([]connectors.CreateConnectorRequest, error) {
	configs := []connectors.CreateConnectorRequest{}
	configFiles, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return configs, err
	}
	for _, fileInfo := range configFiles {
		if fileInfo.IsDir() {
			log.Printf("found unexpected subfolder in folder: %s. This command will not search through it.", filePath)
			continue
		}
		config, err := getConfigFromFile(path.Join(folderPath, fileInfo.Name()), expandEnv)
		if err != nil {
			log.Printf("found unexpected not config file in folder: %s", filePath)
		} else {
			configs = append(configs, config)
		}
	}
	return configs, nil
}

func getConfigFromFile(filePath string, expandEnv bool) (connectors.CreateConnectorRequest, error) {
	config := connectors.CreateConnectorRequest{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return config, err
	}

	return getConfigFromString(string(data), expandEnv)
}

func getConfigFromString(configString string, expandEnv bool) (connectors.CreateConnectorRequest, error) {
	config := connectors.CreateConnectorRequest{}

	if expandEnv {
		// Kafka connect connectors configuration values can contains values looking like enviroment variables
		// In particular, '*.type' fields will be looking like:
		// - "transforms.TimestampConverter.type": "org.apache.kafka.connect.transforms.TimestampConverter$Value",
		// we don't want to expand environment variables for this kind of values
		// e.g. https://docs.confluent.io/platform/current/connect/transforms/timestampconverter.html
		// rather than replacing not found environment variables with empty string, we want to keep the original value
		configString = os.Expand(configString, func(key string) string {
			if key != "Value" && key != "Key" {
				if v, ok := os.LookupEnv(key); ok {
					return v
				}
			}
			return "$" + key
		})
	}

	err := json.Unmarshal([]byte(configString), &config)
	return config, err
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&filePath, "path", "p", "", "path to the config file")
	createCmd.MarkFlagFilename("path")
	createCmd.PersistentFlags().StringVarP(&configString, "string", "s", "", "JSON configuration string")
	createCmd.PersistentFlags().BoolVarP(&sync, "sync", "y", false, "execute synchronously")
	createCmd.PersistentFlags().BoolVar(&expandEnv, "expand-env", false, "expand environment variables for each connector")
}
