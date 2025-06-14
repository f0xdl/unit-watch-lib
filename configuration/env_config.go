package configuration

import (
	"errors"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func LoadEnvConfig(cfg interface{}) error {
	if err := godotenv.Load(); err != nil {
		return errors.New("error loading .env file: " + err.Error())
	}
	if err := env.Parse(cfg); err != nil {
		return errors.New("error parsing .env file: " + err.Error())
	}
	return nil
}
