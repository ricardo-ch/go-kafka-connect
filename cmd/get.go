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
	"errors"

	"github.com/heetch/go-kafka-connect/v4/pkg/connectors"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve information from kafka-connect",
	Long: `Get reads from the kafka-connect REST API.
	It can get the list of all deployed connectors, or details about a single one.`,
	RunE: handleCmd,
}

func handleCmd(cmd *cobra.Command, args []string) error {
	err := validateArgs()
	if err != nil {
		return err
	}

	switch {
	case config:
		return getConfig()
	case status:
		return getStatus()
	case tasks:
		return getTasks()
	default:
		return getConnector()
	}

}

func validateArgs() error {
	if connector == "" {
		return errors.New("please specify the target connector's name")
	}
	if (status && config) || (status && tasks) || (config && tasks) {
		return errors.New("more than one action were provided")
	}

	return nil
}

func getConnector() error {

	client := getClient()
	req := connectors.ConnectorRequest{
		Name: connector,
	}

	resp, err := client.GetConnector(req)
	if err != nil {
		return err
	}

	return printResponse(resp)
}

func getConfig() error {

	client := getClient()
	req := connectors.ConnectorRequest{
		Name: connector,
	}

	resp, err := client.GetConnectorConfig(req)
	if err != nil {
		return err
	}

	return printResponse(resp)
}

func getStatus() error {

	client := getClient()
	req := connectors.ConnectorRequest{
		Name: connector,
	}

	resp, err := client.GetConnectorStatus(req)
	if err != nil {
		return err
	}

	return printResponse(resp)
}

func getTasks() error {

	client := getClient()
	req := connectors.ConnectorRequest{
		Name: connector,
	}

	resp, err := client.GetAllTasks(req)
	if err != nil {
		return err
	}

	return printResponse(resp)
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().StringVarP(&connector, "connector", "n", "", "name of the target's connector")
	getCmd.PersistentFlags().BoolVarP(&status, "status", "s", false, "get the connector's status")
	getCmd.PersistentFlags().BoolVarP(&config, "config", "c", false, "get the connector's config")
	getCmd.PersistentFlags().BoolVarP(&tasks, "tasks", "t", false, "get the connector's tasks list")
}
