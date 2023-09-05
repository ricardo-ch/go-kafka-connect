package examples_go_client

import (
	kafkaConnectClient "github.com/ricardo-ch/go-kafka-connect/v3/lib/connectors"

	"github.com/mitchellh/mapstructure"
)

type Config struct {
	ConnectorClass      string `mapstructure:"connector.class"`
	Topics              string `mapstructure:"topics"`
	TasksMax            string `mapstructure:"tasks.max"`
	ConnectionUrl       string `mapstructure:"connection.url"`
	ConnectionUserName  string `mapstructure:"connection.user"`
	ConnectionPassword  string `mapstructure:"connection.password"`
	InsertMode          string `mapstructure:"insert.mode"`
	DeleteEnabled       bool   `mapstructure:"delete.enabled"`
	PrimaryKeyMode      string `mapstructure:"primary.key.mode"`
	SchemaEvolutionMode string `mapstructure:"schema.evolution.mode"`
	PrimaryKeyFields    string `mapstructure:"primary.key.fields"`
	DatabaseTimezone    string `mapstructure:"database.time.zone"`
}

// This function is an example of Kafka Sink Connect Debezium for JDBC. The configuration for each type of Kafka Connect Sink varies from the plugin it implements.
// Configuration reference: https://debezium.io/documentation/reference/stable/connectors/jdbc.html.
func CreateDebeziumSink(req Config) (kafkaConnectClient.ConnectorResponse, error) {
	r := Config{
		ConnectorClass:      "io.debezium.connector.jdbc.JdbcSinkConnector",
		TasksMax:            "1",
		ConnectionUrl:       "jdbc:mysql://127.0.0.1:3306/test",
		ConnectionUserName:  "root",
		ConnectionPassword:  "testpassword",
		Topics:              "pto.dev.alert_rules",
		PrimaryKeyMode:      "record_key",
		PrimaryKeyFields:    "id",
		InsertMode:          "upsert",
		DeleteEnabled:       true,
		SchemaEvolutionMode: "basic",
		DatabaseTimezone:    "UTC",
	}
	out := make(map[string]interface{})

	err := mapstructure.Decode(r, &out)
	if err != nil {
		panic(err)
	}
	client := kafkaConnectClient.NewClient("http://127.0.0.1:31076")
	re := kafkaConnectClient.CreateConnectorRequest{
		ConnectorRequest: kafkaConnectClient.ConnectorRequest{
			Name: "kafka-connect-sink-name",
		},
		Config: out,
	}

	resp, err := client.CreateConnector(re, false)
	if err != nil {
		return kafkaConnectClient.ConnectorResponse{}, err
	}

	return resp, nil
}
