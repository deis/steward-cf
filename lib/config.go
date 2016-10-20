package lib

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	BrokerAccessScheme      string `envconfig:"BROKER_ACCESS_SCHEME" required:"true"`
	BrokerHost              string `envconfig:"BROKER_HOST" required:"true"`
	BrokerPort              int    `envconfig:"BROKER_PORT" required:"true"`
	BrokerUsername          string `envconfig:"BROKER_USERNAME" required:"true"`
	BrokerPassword          string `envconfig:"BROKER_PASSWORD" required:"true"`
	BrokerRequestTimeoutSec int    `envconfig:"BROKER_REQUEST_TIMEOUT_SEC" default:"5"`
}

func getConfig() (config, error) {
	ret := config{}
	err := envconfig.Process("", &ret)
	return ret, err
}

func (c *config) brokerRequestTimeout() time.Duration {
	return time.Duration(c.BrokerRequestTimeoutSec) * time.Second
}
