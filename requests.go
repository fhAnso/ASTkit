package astkit

import (
	"fmt"
	"net/http"
)

func SendRequestHTTP(client *ASTkitClient, method string, url string) (*http.Response, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to setup %s request: %s", method, err)
	}
	response, err := client.HttpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("could not send %s to %s: %s", method, url, err)
	}
	return response, nil
}
