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
	"crypto/tls"
	"errors"
	"net/http"

	"github.com/ricardo-ch/go-kafka-connect/lib/connectors"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Retrieve information from kafka-connect",
	Long: `Get reads from the kafka-connect REST API.
	It can get the list of all deployed connectors, or details about a single one.
	flags:
		--url -u: url of the kafka-connect server
		--connector -n: name of the target connector
		--status -s: get the connector's status (requires -n)
		--config -c: get the connector's config (requires -n)
		--tasks -t: get the connector's tasks list (requires -n)`,
	RunE: handleCmd,
}

func handleCmd(cmd *cobra.Command, args []string) error {
	err := validateArgs()
	if err != nil {
		return err
	}

	switch true {
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
		return errors.New("Please specify the target connector's name")
	}
	if (status && config) || (status && tasks) || (config && tasks) {
		return errors.New("More than one action were provided")
	}

	return nil
}

func getConnector() error {

	client := connectors.NewClient(url)
	if insecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = client.WithHTTPClient(&http.Client{Transport: tr})
	}

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

	client := connectors.NewClient(url)
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

	client := connectors.NewClient(url)
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

	client := connectors.NewClient(url)
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
