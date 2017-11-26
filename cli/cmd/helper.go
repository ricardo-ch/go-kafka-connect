package cmd

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
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

func getLogger(verbose bool) *zap.Logger{
	config := zap.NewProductionConfig()
	if verbose {
		config.Level.SetLevel(zap.DebugLevel)
	}
	logger, err := config.Build()
	if err != nil {
		// what am i supposed to do if logger init fail?
		fmt.Println(err)
	}
	return logger
}

func getClient() connectors.Client{
	return connectors.NewClient(url, connectors.SetLogger(getLogger(verbose)))
}