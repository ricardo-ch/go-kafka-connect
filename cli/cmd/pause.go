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
	"github.com/ricardo-ch/go-kafka-connect/lib/connectors"
	"github.com/spf13/cobra"
)

type pauseCmdConfig struct {
	sync      bool
	connector string
}

var pause pauseCmdConfig

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: RunEPause,
}

//RunEPause ...
func RunEPause(cmd *cobra.Command, args []string) error {
	req := connectors.ConnectorRequest{
		Name: pause.connector,
	}
	resp, err := connectors.NewClient(url).DeleteConnector(req, delete.sync)
	if err != nil {
		return err
	}
	return printResponse(resp)
}

func init() {
	RootCmd.AddCommand(pauseCmd)

	createCmd.PersistentFlags().BoolVarP(&pause.sync, "sync", "y", false, "wait for asynchronous operation to be done")
	updateCmd.PersistentFlags().StringVarP(&pause.connector, "connector", "n", "", "name of connector to pause")
}
