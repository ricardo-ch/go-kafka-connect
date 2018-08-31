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
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new connector",
	Long: `Deploy a new connector or replace the old version if it alrerady exists.
	This command is executes all its steps synchronously.
	flags:
		--url -u: url of the kafka-connect server
		--file -f: path to the config file
		--string -s: literal configuration string`,
	RunE: RunEDeploy,
}

func RunEDeploy(cmd *cobra.Command, args []string) error {
	config, err := getCreateCmdConfig(cmd)
	if err != nil {
		return err
	}

	return getClient().DeployConnector(config)
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to the config file")
	deployCmd.MarkFlagFilename("file")
	deployCmd.PersistentFlags().StringVarP(&configString, "string", "s", "", "JSON configuration string")
}
