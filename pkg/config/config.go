package config

import (
	"flag"
	"github.com/mdshahjahanmiah/explore-go/logging"
)

type Config struct {
	LoggerConfig logging.LoggerConfig
}

func Load() (Config, error) {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	loggerConfig := logging.LoggerConfig{}
	fs.StringVar(&loggerConfig.CommandHandler, "logger.handler.type", "json", "handler type e.g json, otherwise default will be text type")
	fs.StringVar(&loggerConfig.LogLevel, "logger.log.level", "info", "log level wise logging with fatal log")

	config := Config{
		LoggerConfig: loggerConfig,
	}

	return config, nil
}
