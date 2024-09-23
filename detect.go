package astkit

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ASTkitDetect struct {
	Response *http.Response
}

var AcceptedCMS = map[string][]string{
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

func ASTkitDetectCMS(astkitDetect *ASTkitDetect) (string, error) {
	defer astkitDetect.Response.Body.Close()
	body, err := io.ReadAll(astkitDetect.Response.Body)
	if err != nil {
		return "", fmt.Errorf("an error occured while reading the response body: %s", err)
	}
	for cmsName, indicators := range AcceptedCMS {
		for idx := 0; idx < len(indicators); idx++ {
			if strings.Contains(string(body), indicators[idx]) {
				return cmsName, nil
			}
		}
	}
	return "unknown", nil
}
