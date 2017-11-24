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

	"github.com/ricardo-ch/go-kafka-connect/lib/connectors"
	"github.com/spf13/cobra"
)

var connector string
var status, config, tasks bool

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

	resp, err := client.GetConnectorConfig(req)
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

	getCmd.PersistentFlags().StringVarP(&connector, "connector", "n", "", "connector's name")
	getCmd.PersistentFlags().BoolVarP(&status, "status", "s", false, "connector's status")
	getCmd.PersistentFlags().BoolVarP(&config, "config", "c", false, "connector's configuration")
	getCmd.PersistentFlags().BoolVarP(&tasks, "tasks", "t", false, "connector's tasks list")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
