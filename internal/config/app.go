package config

import (
	"database/sql"

	"github.com/spf13/viper"
)

type BootstrapConfig struct {
	DB *sql.DB
	Config
	Viper *viper.Viper
}

type Config struct {
	WebPort string
	Host    string
	JWTKey  string
}

func NewBoostrapConfig(DB *sql.DB, viper *viper.Viper) *BootstrapConfig {
	return &BootstrapConfig{
		Viper: viper,
		DB:    DB,
		Config: Config{
			WebPort: viper.GetString("PORT"),
			Host:    viper.GetString("HOST"),
			JWTKey:  viper.GetString("JWT_KEY"),
		},
	}
}
