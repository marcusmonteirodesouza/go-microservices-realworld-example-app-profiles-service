package utils

import (
	"os"

	"github.com/marcusmonteirodesouza/go-microservices-realworld-example-app-client/pkg/client"
)

var TestPrefix = "TEST_"

func NewApiClient() client.Client {
	return client.NewClientWithBaseUrl(os.Getenv("API_BASE_URL"))
}
