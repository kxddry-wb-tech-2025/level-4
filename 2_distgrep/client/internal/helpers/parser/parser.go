package parser

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"client/internal/models"
)

// default ports per scheme
var defaultPorts = map[string]string{
	"http":  "80",
	"https": "443",
	"ftp":   "21",
}

// ParseAddress parses an address string into scheme, host, and port.
// If no port is provided but the scheme has a default, that port is used.
func ParseAddress(addr string, defaultScheme string) (*models.ParsedAddr, error) {
	if addr == "" {
		return nil, fmt.Errorf("address is empty")
	}

	// Has scheme?
	if strings.Contains(addr, "://") {
		u, err := url.Parse(addr)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %w", err)
		}

		host, port := splitHostPort(u.Host)

		// apply default if missing
		if port == "" {
			if def, ok := defaultPorts[u.Scheme]; ok {
				port = def
			} else {
				return nil, fmt.Errorf("no port specified and no default for scheme %q", u.Scheme)
			}
		}

		return &models.ParsedAddr{
			Raw:    addr,
			Scheme: u.Scheme,
			Host:   host,
			Port:   port,
		}, nil
	}

	// If no scheme, treat as host[:port]
	host, port := splitHostPort(addr)
	if port == "" {
		return nil, fmt.Errorf("address must include port: %s", addr)
	}

	return &models.ParsedAddr{
		Raw:    addr,
		Scheme: defaultScheme,
		Host:   host,
		Port:   port,
	}, nil
}

func splitHostPort(hp string) (string, string) {
	if strings.Contains(hp, ":") {
		h, p, err := net.SplitHostPort(hp)
		if err == nil {
			return h, p
		}
	}
	return hp, "" // no port
}
