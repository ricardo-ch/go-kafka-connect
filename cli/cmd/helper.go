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

func getClient() connectors.HighLevelClient {
	client := connectors.NewClient(url)
	if verbose {
		client.SetDebug()
	}
	if SSLInsecure {
		client.SetInsecureSSL()
	}
	if len(SSLClientCertificate) + len(SSLClientPrivateKey) > 0 {
		client.SetClientCertificates(SSLClientCertificate, SSLClientPrivateKey)
	}
	return client
}
