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
	"github.com/ricardo-ch/go-kafka-connect/lib/connectors"
	"encoding/json"
	"os"
	"strings"
	"github.com/pkg/errors"
	"fmt"
)

var(
	file string
	configString string
	syncCreate bool
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `Create a connector using either config file or string
	flags:
		--url -u : url of kafka-connect server
		--file -f : use file to define config
		--string -s : use string to define config
	`,
	RunE: RunECreate,
}

func RunECreate(cmd *cobra.Command, args []string) error {
	config := connectors.CreateConnectorRequest{}

	if cmd.Flag("file").Changed {
		fileReader, err := os.Open(file)
		if err != nil {
			return err
		}

		err = json.NewDecoder(fileReader).Decode(&config)
		if err != nil {
			return err
		}

	} else if cmd.Flag("string").Changed{
		err := json.NewDecoder(strings.NewReader(configString)).Decode(&config)
		if err != nil {
			return err
		}
	} else {
		return errors.New("neither file nor string was supplied")
	}


	resp, err :=  connectors.NewClient(url).CreateConnector(config, syncCreate)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	createCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "path to config file")
	createCmd.PersistentFlags().StringVarP(&configString, "string", "s", "", "json encoded string of config")
	createCmd.PersistentFlags().BoolVarP(&syncCreate, "sync", "n", false, "wait for asynchronous operation to be done")

}
