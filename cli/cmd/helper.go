package cmd

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"

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
	if basicAuthUsername != "" && basicAuthPassword != "" {
		client.SetBasicAuth(basicAuthUsername, basicAuthPassword)
	}
	if len(SSLClientCertificate) > 0 && len(SSLClientPrivateKey) > 0 {
		cert, err := tls.LoadX509KeyPair(SSLClientCertificate, SSLClientPrivateKey)
		if err != nil {
			log.Fatalf("client: loadkeys: %s", err)
		} else {
			client.SetClientCertificates(cert)
		}
	}
	return client
}
