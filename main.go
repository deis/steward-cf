package main

import (
	"os"

	"github.com/deis/steward-cf/lib"
	"github.com/deis/steward-framework/runner"
	"github.com/juju/loggo"
)

func main() {
	// Default to INFO level logging until we load configuration and learn the desired log level
	logger.SetLogLevel(loggo.INFO)
	logger.Infof("steward-cf version %s starting", version)

	cfg, err := getConfig()
	if err != nil {
		logger.Criticalf("error getting config: %s", err)
		os.Exit(1)
	}
	logger.SetLogLevel(cfg.logLevel())

	cataloger, lifecycler, err := lib.GetComponents()
	if err != nil {
		logger.Criticalf("error getting components: %s", err)
		os.Exit(1)
	}

	if err = runner.Run(
		cfg.BrokerName,
		cfg.Namespaces,
		cataloger,
		lifecycler,
		cfg.getMaxAsyncDuration(),
		cfg.APIPort,
	); err != nil {
		logger.Criticalf("error running steward-framework: %s", err)
		os.Exit(1)
	}
}
