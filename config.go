package main

import (
	"time"

	"github.com/juju/loggo"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	MaxAsyncMinutes int    `envconfig:"MAX_ASYNC_MINUTES" default:"60"`
	APIPort         int    `envconfig:"API_PORT" default:"8080"`
	LogLevel        string `envconfig:"LOG_LEVEL" default:"info"`
}

func getConfig() (*config, error) {
	ret := &config{}
	if err := envconfig.Process("", ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// getMaxAsyncDuration returns the maximum time spent on an aysnchronous provisioning or
// deprovisioning before giving up on polling for the last operation
func (c *config) getMaxAsyncDuration() time.Duration {
	return time.Duration(time.Duration(c.MaxAsyncMinutes) * time.Minute)
}

func (c config) logLevel() loggo.Level {
	switch c.LogLevel {
	case "trace":
		return loggo.TRACE
	case "debug":
		return loggo.DEBUG
	case "info":
		return loggo.INFO
	case "warning":
		return loggo.WARNING
	case "error":
		return loggo.ERROR
	case "critical":
		return loggo.CRITICAL
	default:
		return loggo.INFO
	}
}
