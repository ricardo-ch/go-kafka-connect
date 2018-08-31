package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ricardo-ch/go-kafka-connect/lib/connectors"
)

func printResponse(response interface{}) error {
	out, err := json.MarshalIndent(response, "", "    ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func getClient() *connectors.Client {
	client := connectors.NewClient(url)
	if verbose {
		client = client.WithDebug()
	}
	if SSLInsecure {
		client = client.WithInsecureSSL()
	}
	return client
}
