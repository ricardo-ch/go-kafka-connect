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
	"github.com/ricardo-ch/go-kafka-connect/v3/lib/connectors"
	"github.com/spf13/cobra"
)

// resumeCmd represents the resume command
var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a connector",
	RunE:  RunEResume,
}

// RunEResume ...
func RunEResume(cmd *cobra.Command, args []string) error {
	req := connectors.ConnectorRequest{
		Name: connector,
	}
	resp, err := getClient().ResumeConnector(req, sync)
	if err != nil {
		return err
	}
	return printResponse(resp)
}

func init() {
	RootCmd.AddCommand(resumeCmd)

	resumeCmd.PersistentFlags().BoolVarP(&sync, "sync", "y", false, "execute synchronously")
	resumeCmd.PersistentFlags().StringVarP(&connector, "connector", "n", "", "name of the target connector")
}
