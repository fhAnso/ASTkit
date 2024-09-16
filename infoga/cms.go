package infoga

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	astkitClient "github.com/fhAnso/ASTkit/client"
)

var acceptedCMS = map[string][]string{
	"WordPress":        {"wp-content", "wp-includes"},
	"Joomla":           {"Joomla!"},
	"Drupal":           {"Drupal"},
	"Magento":          {"Magento"},
	"Shopify":          {"Shopify"},
	"Blogger":          {"blogspot"},
	"Wix":              {"wix"},
	"Squarespace":      {"squarespace"},
	"TYPO3":            {"typo3"},
	"Concrete5":        {"concrete5"},
	"PrestaShop":       {"prestashop"},
	"OpenCart":         {"catalog"},
	"Ghost":            {"ghost"},
	"ExpressionEngine": {"expressionEngine"},
	"Craft CMS":        {"craft"},
	"SilverStripe":     {"silverstripe"},
	"DotNetNuke":       {"dnn"},
	"Weebly":           {"weebly"},
}

func ASTkitDetectCMS(client *astkitClient.ASTkitClient) (string, error) {
	response, err := http.Get(client.URL)
	if err != nil {
		return "", fmt.Errorf("failed to sent GET: %s", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("an error occured while reading the response body: %s", err)
	}
	for cmsName, indicators := range acceptedCMS {
		for idx := 0; idx < len(indicators); idx++ {
			if strings.Contains(string(body), indicators[idx]) {
				return cmsName, nil
			}
		}
	}
	return "unknown", nil
}
