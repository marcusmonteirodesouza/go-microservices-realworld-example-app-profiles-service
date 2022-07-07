package utils

import (
	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-client/pkg/client"
)

var TestPrefix = "TEST_"

func NewApiClient() client.Client {
	return client.NewClientWithBaseUrl("http://localhost:8081")
}
