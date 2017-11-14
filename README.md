# go-kafka-connect
Go library providing bindings for the Kafka connect API

# Setup environment
Required:
 - Go 1.9
 - Docker (for testing purpose only)

run `make install` to install dependencies


# Testing
For now, only integration test are available.
run `docker-compose up` then wait patiently until it boots and run `make test-integration`

Right now, integration test only run locally