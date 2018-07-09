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
	"net/http"

	"github.com/ricardo-ch/go-kafka-connect/lib/connectors"
	"github.com/spf13/cobra"
)

// pauseCmd represents the pause command
var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause a connector",
	Long: `Suspend a connector without deleting it.
	flags:
		--url -u: url of the kafka-connect server
		--connector -n: name of the target connector
		--sync -y: execute synchronously`,
	RunE: RunEPause,
}

//RunEPause ...
func RunEPause(cmd *cobra.Command, args []string) error {
	req := connectors.ConnectorRequest{
		Name: connector,
	}

	client := connectors.NewClient(url)
	if insecureSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = client.WithHTTPClient(&http.Client{Transport: tr})
	}

	resp, err := client.PauseConnector(req, sync)
	if err != nil {
		return err
	}
	return printResponse(resp)
}

func init() {
	RootCmd.AddCommand(pauseCmd)

	pauseCmd.PersistentFlags().BoolVarP(&sync, "sync", "y", false, "execute synchronously")
	pauseCmd.PersistentFlags().StringVarP(&connector, "connector", "n", "", "name of the target connector")
}
