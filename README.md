# Scootin

### System requirements
`docker-compose`

### Run
`docker-compose up`

### Use
I have implemented a `service-client` to speak to the REST-API
so no need to use http client, just call the API in code as illustrated in `client/client_test.go`


### Test
you need to run `docker-compose up` first to run the integration tests and the simulation with
`go test ./...`


