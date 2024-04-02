package config

import (
	"github.com/spf13/viper"
)

func NewViper() *viper.Viper {
	config := viper.New()

	config.SetConfigName(".env")
	config.SetConfigType("env")

	// Define possible paths to search for the .env file
	possiblePaths := []string{"./", "./../"}
	// Attempt to read the .env file from each possible path
	var err error
	for _, path := range possiblePaths {
		config.SetConfigFile(path + ".env")
		if err = config.ReadInConfig(); err == nil {
			return config
		}
	}

	if err != nil {
		config.AutomaticEnv()
		return config
	}

	panic("Fatal error: .env file not found in any specified path")
}
