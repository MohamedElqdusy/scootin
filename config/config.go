package config

import (
	"github.com/kelseyhightower/envconfig"
)

type PostgreConfig struct {
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	PostgresDataBase string `envconfig:"POSTGRES_DATABASE"`
	PostgresHost     string `envconfig:"POSTGRES_HOST"`
	PostgresPort     string `envconfig:"POSTGRES_PORT"`
}

func IniatilizePostgreConfig() (*PostgreConfig, error) {
	var p PostgreConfig
	if err := envconfig.Process("", &p); err != nil {
		return nil, err
	}
	return &p, nil
}
