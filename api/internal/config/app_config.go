package config

import (
	env "github.com/caarlos0/env/v6"
)

var globalAppConfig *appConfig

type appConfig struct {
	MongoDbConnParams string `env:"PROMPTU_MONGODB_CONN_PARAMS" envDefault:"retryWrites=true&w=majority"`
	MongoDbUrl        string `env:"PROMPTU_MONGODB_URL,required"`
}

func AppConfig() *appConfig {
	if globalAppConfig == nil {
		cfg := &appConfig{}

		if err := env.Parse(&cfg); err != nil {
			panic(err)
		}

		globalAppConfig = cfg
	}

	return globalAppConfig
}
