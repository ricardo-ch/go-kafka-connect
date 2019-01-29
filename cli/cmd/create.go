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
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/ricardo-ch/go-kafka-connect/lib/connectors"
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
	configs, err := getCreateCmdConfig(cmd)
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

func getCreateCmdConfig(cmd *cobra.Command) ([]connectors.CreateConnectorRequest, error) {
	var configs []connectors.CreateConnectorRequest

	if cmd.Flag("file").Changed {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return nil, errors.Wrapf(err, "error while trying to find file or folder: %v", file)
		}
		if fileInfo.IsDir() {
			configs, err = getConfigFromFolder(file)
			if err != nil {
				return nil, err
			}
		} else {
			config, err := getConfigFromFile(file)
			if err != nil {
				return nil, err
			}
			configs = append(configs, config)
		}

	} else if cmd.Flag("string").Changed {
		config := connectors.CreateConnectorRequest{}
		err := json.NewDecoder(strings.NewReader(configString)).Decode(&config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	} else {
		return nil, errors.New("neither file nor string was supplied")
	}
	return configs, nil
}

func getConfigFromFolder(folderPath string) ([]connectors.CreateConnectorRequest, error) {
	configs := []connectors.CreateConnectorRequest{}
	configFiles, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return configs, err
	}
	for _, fileInfo := range configFiles {
		if fileInfo.IsDir() {
			log.Printf("found unexpected subfolder in folder: %s. This command will not search through it.", file)
			continue
		}
		config, err := getConfigFromFile(path.Join(folderPath, fileInfo.Name()))
		if err != nil {
			log.Printf("found unexpected not config file in folder: %s", file)
		} else {
			configs = append(configs, config)
		}
	}
	return configs, nil
}

func getConfigFromFile(filePath string) (connectors.CreateConnectorRequest, error) {
	config := connectors.CreateConnectorRequest{}
	fileReader, err := os.Open(filePath)
	if err != nil {
		return config, err
	}

	err = json.NewDecoder(fileReader).Decode(&config)
	return config, err
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to the config file")
	createCmd.MarkFlagFilename("file")
	createCmd.PersistentFlags().StringVarP(&configString, "string", "s", "", "JSON configuration string")
	createCmd.PersistentFlags().BoolVarP(&sync, "sync", "y", false, "execute synchronously")

}
