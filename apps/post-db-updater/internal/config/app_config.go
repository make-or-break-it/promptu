package config

import (
	env "github.com/caarlos0/env/v6"
)

var globalAppConfig *appConfig

type appConfig struct {
	MongoDbConnParams string `env:"PROMPTU_MONGODB_CONN_PARAMS" envDefault:"retryWrites=true&w=majority"`
	MongoDbUrl        string `env:"PROMPTU_MONGODB_URL,required"`
	KafkaBrokers      string `env:"KAFKA_BROKERS,required" envDefault:"149.248.217.129:9092"`
	KafkaVersion      string `env:"KAFKA_VERSION,required" envDefault:"2.8.1"`
	PostTopic         string `env:"POST_TOPIC,required" envDefault:"posts"`
}

func AppConfig() *appConfig {
	if globalAppConfig == nil {
		cfg := &appConfig{}

		if err := env.Parse(cfg); err != nil {
			panic(err)
		}

		globalAppConfig = cfg
	}

	return globalAppConfig
}
