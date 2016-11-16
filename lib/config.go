package lib

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	BrokerRequestTimeoutSec int `envconfig:"BROKER_REQUEST_TIMEOUT_SEC" default:"5"`
}

func getConfig() (config, error) {
	ret := config{}
	err := envconfig.Process("", &ret)
	return ret, err
}

func (c *config) brokerRequestTimeout() time.Duration {
	return time.Duration(c.BrokerRequestTimeoutSec) * time.Second
}
