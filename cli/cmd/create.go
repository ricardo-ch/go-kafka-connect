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

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new connector",
	Long: `Create a connector using either a config file or a literal string
	flags:
		--url -u: url of the kafka-connect server
		--file -f: path to the config file
		--string -s: JSON configuration string
		--sync -y: execute synchronously
	`,
	RunE: RunECreate,
}

//RunECreate ...
func RunECreate(cmd *cobra.Command, args []string) error {
	config, err := getCreateCmdConfig(cmd)
	if err != nil {
		return err
	}

	resp, err := connectors.NewClient(url).CreateConnector(config, sync)
	if err != nil {
		return err
	}

	return printResponse(resp)
}

func getCreateCmdConfig(cmd *cobra.Command) (connectors.CreateConnectorRequest, error) {
	config := connectors.CreateConnectorRequest{}

	if cmd.Flag("file").Changed {
		fileReader, err := os.Open(file)
		if err != nil {
			return config, err
		}

		err = json.NewDecoder(fileReader).Decode(&config)
		if err != nil {
			return config, err
		}

	} else if cmd.Flag("string").Changed {
		err := json.NewDecoder(strings.NewReader(configString)).Decode(&config)
		if err != nil {
			return config, err
		}
	} else {
		return config, errors.New("neither file nor string was supplied")
	}
	return config, nil
}

func init() {
	RootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to the config file")
	createCmd.MarkFlagFilename("file")
	createCmd.PersistentFlags().StringVarP(&configString, "string", "s", "", "JSON configuration string")
	createCmd.PersistentFlags().BoolVarP(&sync, "sync", "y", false, "execute synchronously")

}
