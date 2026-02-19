package config

import (
	"github.com/spf13/viper"
)

// Config stores all application configuration
type Config struct {
	Environment string `mapstructure:"ENVIRONMENT"`
	Port        string `mapstructure:"PORT"`
	DBUrl       string `mapstructure:"DATABASE_URL"`
	RedisURL    string `mapstructure:"REDIS_URL"`
	JWTSecret   string `mapstructure:"JWT_SECRET"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv() // Override with environment variables if set

	err = viper.ReadInConfig()
	if err != nil {
		// If the config file is not found, we can ignore the error and rely on environment variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
	}

	err = viper.Unmarshal(&config)
	return
}
