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
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	url                  string
	connector            string
	filePath             string
	configString         string
	sync                 bool
	status               bool
	config               bool
	tasks                bool
	verbose              bool
	SSLInsecure          bool
	jsonLog              bool
	parallel             int
	pauseBeforeDeploy    bool
	expandEnv            bool
	SSLClientCertificate string
	SSLClientPrivateKey  string
	basicAuthUsername    string
	basicAuthPassword    string
	extraHeaders         HeadersFlag
)

var RootCmd = &cobra.Command{
	Use:   "kc-cli [command] [args]",
	Short: "CLI wrapper for kafka-connect API",
	Long: `This is a small tool to perform all available task on kafka-connect API via a CLI.
also contains two 'bonus' features:
- deploy connectors
- synchronous operations
`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&url, "url", "u", "http://localhost:8083", "kafka connect URL")
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, `/!\ very verbose`)
	RootCmd.PersistentFlags().BoolVarP(&SSLInsecure, "insecure-skip-verify", "i", false, `skip SSL/TLS verification`)
	RootCmd.PersistentFlags().StringVarP(&SSLClientCertificate, "ssl-client-certificate", "C", "", `path to client certificate, must contain PEM encoded data`)
	RootCmd.PersistentFlags().StringVarP(&SSLClientPrivateKey, "ssl-client-key", "K", "", `path to client private key`)
	RootCmd.PersistentFlags().StringVarP(&basicAuthUsername, "username", "U", "", `HTTP Basic Auth username`)
	RootCmd.PersistentFlags().StringVarP(&basicAuthPassword, "password", "P", "", `HTTP Basic Auth password`)
	RootCmd.PersistentFlags().VarP(&extraHeaders, "header", "H", "extra HTTP headers to attach to REST API requests")
}
