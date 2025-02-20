package config

import (
	"log"

	"github.com/caarlos0/env/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port        string `env:"PORT"`
	DB          DBConfig
	Redis       RedisConfig
	Eth         EthConfig
	OauthClient GoogleOAuthClient
	CORS        struct {
		AllowOrigins     []string `envconfig:"CORS_ALLOW_ORIGINS" default:"*"`
		AllowCredentials bool     `envconfig:"CORS_ALLOW_CREDENTIALS" default:"true"`
		Debug            bool     `envconfig:"CORS_DEBUG" default:"false"`
	}
}

func Load() (*Config, error) {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	// Self-validation
	// if cfg.OauthClient.ClientID == "" {
	// 	log.Fatal("OAUTH_CLIENT_ID is required")
	// }
	// if cfg.OauthClient.ClientSecret == "" {
	// 	log.Fatal("OAUTH_CLIENT_SECRET is required")
	// }
	// if cfg.OauthClient.RedirectURI == "" {
	// 	log.Fatal("OAUTH_REDIRECT_URI is required")
	// }
	return cfg, nil
}
