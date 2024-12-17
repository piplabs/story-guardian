package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// AppConfig represents the application's configuration.
type AppConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
}

// NewAppConfig initializes a new AppConfig instance.
func NewAppConfig() (*AppConfig, error) {
	// Set environment variable prefix for configuration
	viper.SetEnvPrefix("cipherowl")
	viper.AutomaticEnv()

	clientID := viper.GetString("client_id")
	clientSecret := viper.GetString("client_secret")

	// Validate configuration inputs
	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("both CLIENT_ID and CLIENT_SECRET environment variables are required")
	}

	return &AppConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}
