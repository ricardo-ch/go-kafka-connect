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
	"errors"
	"os"
	"strings"

	"github.com/ricardo-ch/go-kafka-connect/lib/connectors"
	"github.com/spf13/cobra"
)

type updateCmdConfig struct {
	file         string
	configString string
	connector    string
}

var update updateCmdConfig

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Updater a connector",
	RunE:  RunEUpdate,
}

//RunEUpdate ...
func RunEUpdate(cmd *cobra.Command, args []string) error {
	req := connectors.CreateConnectorRequest{}

	var err error
	req.Config, err = getUpdateCmdConfig(cmd)
	if err != nil {
		return err
	}

	req.Name = update.connector
	resp, err := getClient().UpdateConnector(req, sync)
	if err != nil {
		return err
	}

	return printResponse(resp)
}

func getUpdateCmdConfig(cmd *cobra.Command) (map[string]interface{}, error) {
	config := map[string]interface{}{}

	if cmd.Flag("file").Changed {
		fileReader, err := os.Open(update.file)
		if err != nil {
			return config, err
		}

		err = json.NewDecoder(fileReader).Decode(&config)
		if err != nil {
			return config, err
		}

	} else if cmd.Flag("string").Changed {
		err := json.NewDecoder(strings.NewReader(update.configString)).Decode(&config)
		if err != nil {
			return config, err
		}
	} else {
		return config, errors.New("neither file nor string was supplied")
	}
	return config, nil
}

func init() {
	RootCmd.AddCommand(updateCmd)

	updateCmd.PersistentFlags().StringVarP(&update.file, "file", "f", "", "path to the config file")
	updateCmd.MarkFlagFilename("file")
	updateCmd.PersistentFlags().StringVarP(&update.configString, "string", "s", "", "JSON configuration string")
	updateCmd.PersistentFlags().StringVarP(&update.connector, "connector", "n", "", "name of the target connector")
	updateCmd.MarkFlagRequired("connector")
	updateCmd.PersistentFlags().BoolVarP(&sync, "sync", "y", false, "execute synchronously")
}
