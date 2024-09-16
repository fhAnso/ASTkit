package astkit

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type ASTkitClient struct {
	URL        string
	Timeout    int // ms
	EnableTor  bool
	HttpClient *http.Client // Defined in ASTkitHttpClientInit()
}

func setupProxy(proxyClient *http.Client, proxyUrl string) error {
	parseProxyUrl, err := url.Parse(proxyUrl)
	if err != nil {
		return fmt.Errorf("failed to parse given proxy url: %s\n%s", proxyUrl, err)
	}
	proxyClient.Transport = &http.Transport{Proxy: http.ProxyURL(parseProxyUrl)}
	return nil
}

func (astClient *ASTkitClient) ASTkitHttpClientInit() error {
	astClient.HttpClient = &http.Client{
		Timeout: time.Duration(astClient.Timeout) * time.Millisecond,
	}
	if astClient.EnableTor {
		const proxyUrl = "socks5://127.0.0.1:9050"
		if err := setupProxy(astClient.HttpClient, proxyUrl); err != nil {
			return err
		}
	}
	return nil
}
