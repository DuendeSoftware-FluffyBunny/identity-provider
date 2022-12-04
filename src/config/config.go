package config

import (
	"log"
	"os"
)

type Config struct {
	KratosPublicEndpoint string
	KratosAdminEndpoint  string
	HydraPublicEndpoint  string
	HydraAdminEndpoint   string
	OAuth2ClientID       string
	OAuth2ClientSecret   string
	AppPort              string
}

func NewConfig() *Config {
	return &Config{
		KratosPublicEndpoint: mustGetEnv("KRATOS_PUBLIC_ENDPOINT_URL"),
		KratosAdminEndpoint:  mustGetEnv("KRATOS_ADMIN_ENDPOINT_URL"),
		HydraPublicEndpoint:  mustGetEnv("HYDRA_PUBLIC_ENDPOINT_URL"),
		HydraAdminEndpoint:   mustGetEnv("HYDRA_ADMIN_ENDPOINT_URL"),
		OAuth2ClientID:       mustGetEnv("OAUTH2_CLIENT_ID"),
		OAuth2ClientSecret:   mustGetEnv("OAUTH2_CLIENT_SECRET"),
		AppPort:              mustGetEnv("APP_PORT"),
	}
}

func mustGetEnv(key string) string {
	env := os.Getenv(key)
	if env == "" {
		log.Fatalf("could not find variable %q", key)
	}
	return env
}
