package config

import (
	"github.com/spf13/viper"
)

// UserConfig stores all configuration of user application
type UserConfig struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
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
