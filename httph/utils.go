package httph

import (
	"fmt"
)

type HTTP int

const (
	Basic HTTP = iota
	Secure
)

func MakeUrl(http HTTP, subdomain string) string {
	var proto string
	switch http {
	case Basic:
		proto = "http://"
	case Secure:
		proto = "https://"
	}
	return fmt.Sprintf("%s%s", proto, subdomain)
}
