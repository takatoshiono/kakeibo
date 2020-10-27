# Hello World Samples

## System Tests

### Pub/Sub

1. `export FUNCTIONS_TOPIC=example-topic`

1. `export GCP_PROJECT=your-project-name`

1. `gcloud functions deploy HelloPubSub --runtime=go113 --trigger-topic=$FUNCTIONS_TOPIC`

1. `go test -v ./hello_pubsub_system_test.go`
