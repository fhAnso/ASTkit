package astkit

import (
	"fmt"
	"net/http"
)

func SendRequest(client *ASTkitClient, method string, url string) (*http.Response, error) {
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.HttpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send %s: %s", method, err)
	}
	return response, nil
}
