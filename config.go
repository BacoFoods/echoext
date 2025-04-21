package echoext

import (
	"fmt"
	"strings"
)

type ServerConfig struct {
	PathPrefix       string
	Host             string
	Port             int
	HealthcheckPath  string
	SkipPaths        []string
	SwaggerConfig    SwaggerConfig
	ExtraCORSHeaders []string
}

func (c *ServerConfig) escapePrefix() string {
	if c.PathPrefix == "" {
		return "/"
	}

	prefix := c.PathPrefix
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

func (c *ServerConfig) escapeHealthcheckSuffix() string {
	if c.HealthcheckPath == "" {
		return "/healthcheck"
	}

	suffix := c.HealthcheckPath
	// ensure leading /
	if suffix[0] != '/' {
		suffix = "/" + suffix
	}

	// remove trailing /
	if suffix[len(suffix)-1] == '/' {
		suffix = suffix[:len(suffix)-1]
	}

	// paths are case insensitive

	return strings.ToLower(suffix)
}

func (c *ServerConfig) escapeSkipPaths() []string {
	escapedPaths := []string{"/", c.healthcheckFullPath()}
	for _, path := range c.SkipPaths {
		// ensure leading /
		if path[0] != '/' {
			path = "/" + path
		}

		// remove trailing /
		if path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}

		// paths are case insensitive
		escapedPaths = append(escapedPaths, strings.ToLower(path))
	}

	return escapedPaths
}

func (c *ServerConfig) escapePort() int {
	if c.Port == 0 {
		return 8080
	}

	return c.Port
}

func (c *ServerConfig) escapeHost() string {
	if c.Host == "" {
		return "0.0.0.0:" + fmt.Sprint(c.escapePort())
	}

	host := c.Host
	// remove trailing /
	if host[len(host)-1] == '/' {
		host = host[:len(host)-1]
	}

	// hosts are case insensitive

	return strings.ToLower(host) + ":" + fmt.Sprint(c.escapePort())
}

func (c *ServerConfig) swaggerPath() string {
	return c.escapePrefix() + c.SwaggerConfig.escapePrefix()
}

func (c *ServerConfig) healthcheckFullPath() string {
	return c.escapePrefix() + c.escapeHealthcheckSuffix()
}
