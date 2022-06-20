# go-kafka-connect

Go project containing two different sub-projects: a kafka-connect client library and a CLI to use it.

## Kafka-connect client library

This library is to be used as an abstraction layer for the kafka-connect REST API.
It currently implements the following API calls:

- Create a connector
- Update a connector
- Delete a connector
- Pause a connector
- Resume a connector
- Restart a connector
- Get a connector's details (overview, configuration, status or tasks list)

It also contains two 'bonus' features:

- Do synchronously: All calls to the REST API trigger an asynchronous function on kafka-connect.
  This feature lets the library check regularly if the action has taken effect on kafka-connect's side,
  and considers the request as completed only when the consequences of the command can be verified.
  It allows users of this library to use its functions in a synchronous way.
- Deploy connector, a function used to deploy a connector, or replace an existing one gracefully.
  This function checks if the target connector exists. If it exists, it will then be paused before being updated.
  Before being updating it check the current config. If it match the deployment's config, nothing will be done.
  The new connector is then deployed, and resumed. This function is always synchronous.

## Installation

Install binary using go:

```bash
go install github.com/heetch/go-kafka-connect/cmd/kc-cli@latest
kc-cli --help
```

NB: this assume you have go installed and `$GOBIN` inside your `$PATH`.

## Example of command

- Get help

```bash
./kc-cli --help
```

- Deploy a connector which config is stored in a json:

```bash
./kc-cli deploy -u http://kafka-connect.local -f my-connector-config.json
```

- Deploy a bunch of connector in parallel and wait for the end:

```bash
jobs=''

./kc-cli deploy -u http://kafka-connect.local -f my-connector-config.json & jobs="$jobs $!"
./kc-cli deploy -u http://kafka-connect.local -f my-connector-2-config.json & jobs="$jobs $!"

status=0
for job in $jobs; do
  wait $job
  if [ $? != '0' ]; then $status=1; fi
done
if [ $status != 0 ]; then exit $status; fi
```

- Get connector status

```bash
./kc-cli -u http://kafka-connect.local get --status -n my-connector
```

## Setup environment for development

Required:

- Go 1.17
- Docker (for testing purpose only)

run `go get -u github.com/ricardo-ch/go-kafka-connect`
then inside repo run: `make install` to install dependencies

## Testing

For now, only integration test are available.
run `make test-integration`

Right now, integration tests only run locally

note: you may also run `make rundep`, if you just want to run a kafka-connect cluster in background for manual testing.
