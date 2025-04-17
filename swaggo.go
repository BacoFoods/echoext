package echoext

import "strings"

type SwaggerConfig struct {
	Prefix string
}

func (c *SwaggerConfig) escapePrefix() string {
	if c.Prefix == "" {
		return "/docs"
	}

	prefix := c.Prefix
	// ensure leading /
	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	// remove trailing /
	if prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}

	// paths are case insensitive

	return strings.ToLower(prefix)
}
