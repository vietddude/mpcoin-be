package config

import "github.com/caarlos0/env/v10"

type Config struct {
	Port  string `env:"PORT" envDefault:"5001"`
	DB    DBConfig
	Redis RedisConfig
	Eth   EthConfig
	CORS  struct {
		AllowOrigins     []string `envconfig:"CORS_ALLOW_ORIGINS" default:"*"`
		AllowCredentials bool     `envconfig:"CORS_ALLOW_CREDENTIALS" default:"true"`
		Debug            bool     `envconfig:"CORS_DEBUG" default:"false"`
	}
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
