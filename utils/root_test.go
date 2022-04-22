package utils

import "os"

func makeTestClient() *Client {
	apiKey := os.Getenv("L27_TEST_KEY")
	apiUrl := os.Getenv("L27_TEST_API")

	return NewAPIClient(apiUrl, apiKey)
}