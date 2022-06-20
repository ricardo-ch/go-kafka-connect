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
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a new connector",
	Long: `Deploy a new connector or replace the old version if it alrerady exists.
	This command is executes all its steps synchronously.`,
	RunE: RunEDeploy,
}

func RunEDeploy(cmd *cobra.Command, args []string) error {
	configs, err := getCreateCmdConfig(cmd, expandEnv)
	if err != nil {
		return err
	}

	client := getClient()
	client.SetParallelism(parallel)
	client.SetPauseBeforeDeploy(pauseBeforeDeploy)
	if jsonLog {
		client.SetLogFormatter(&logrus.JSONFormatter{})
	}

	return client.DeployMultipleConnector(configs)
}

func init() {
	RootCmd.AddCommand(deployCmd)

	deployCmd.PersistentFlags().StringVarP(&filePath, "path", "p", "", "path to the config file or folder")
	deployCmd.MarkFlagFilename("path")
	deployCmd.PersistentFlags().StringVarP(&configString, "string", "s", "", "JSON configuration string")
	deployCmd.PersistentFlags().IntVarP(&parallel, "parallel", "r", 3, "limit of parallel call to kafka-connect")
	deployCmd.PersistentFlags().BoolVar(&pauseBeforeDeploy, "pause", true, "pause connector before performing it's deployment")
	deployCmd.PersistentFlags().BoolVar(&expandEnv, "expand-env", false, "expand environment variables for each connector file using os.Expand, skip $Key/$Value, if value is not found in env, it will be returned as $YOUR_DEFINED_ENV_VAR")
	deployCmd.PersistentFlags().BoolVarP(&jsonLog, "json-log", "j", false, "Log in JSON format")
}
