package config

import (
	"github.com/spf13/viper"
)

// UserConfig stores all configuration of user application
type UserConfig struct {
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
	BindAddr   string `mapstructure:"BIND_ADDR"`
}

// LoadConfig parses configuration and stores the result in
// the value pointed to by config.
func LoadConfig(name string, path string, config interface{}) error {
	var err error

	viper.AddConfigPath(path)
	viper.SetConfigName(name)
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(config)
	return err
}

// LoadUserConfig TODO: add description
func LoadUserConfig(name string, path string) (UserConfig, error) {
	var cfg UserConfig

	err := LoadConfig(name, path, &cfg)
	if err != nil {
		return UserConfig{}, err
	}

	return cfg, nil
}
