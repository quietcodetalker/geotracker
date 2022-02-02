package config

import (
	"github.com/spf13/viper"
)

// LocationConfig stores all configuration of user application
type LocationConfig struct {
	AppEnv       string `mapstructure:"APP_ENV"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	DBSSLMode    string `mapstructure:"DB_SSLMODE"`
	BindAddrHTTP string `mapstructure:"BIND_ADDR_HTTP"`
	BindAddrGRPC string `mapstructure:"BIND_ADDR_GRPC"`
	HistoryAddr  string `mapstructure:"HISTORY_ADDR"`
}

// HistoryConfig stores all configuration of user application
type HistoryConfig struct {
	AppEnv       string `mapstructure:"APP_ENV"`
	DBDriver     string `mapstructure:"DB_DRIVER"`
	DBHost       string `mapstructure:"DB_HOST"`
	DBPort       string `mapstructure:"DB_PORT"`
	DBUser       string `mapstructure:"DB_USER"`
	DBPassword   string `mapstructure:"DB_PASSWORD"`
	DBName       string `mapstructure:"DB_NAME"`
	DBSSLMode    string `mapstructure:"DB_SSLMODE"`
	BindAddrHTTP string `mapstructure:"BIND_ADDR_HTTP"`
	BindAddrGRPC string `mapstructure:"BIND_ADDR_GRPC"`
	LocationAddr string `mapstructure:"LOCATION_ADDR"`
}

// LoadConfig parses configuration and stores the result in
// the value pointed to by config.
func LoadConfig(name string, path string, config interface{}) error {
	var err error

	if path != "" {
		viper.AddConfigPath(path)
	}

	if name != "" {
		viper.SetConfigName(name)
	}
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	viper.SetDefault("AppEnv", "production")

	err = viper.ReadInConfig()
	if err != nil {
		return err
	}

	err = viper.Unmarshal(config)
	return err
}

// LoadLocationConfig TODO: add description
func LoadLocationConfig(name string, path string) (LocationConfig, error) {
	var cfg LocationConfig

	err := LoadConfig(name, path, &cfg)
	if err != nil {
		return LocationConfig{}, err
	}

	return cfg, nil
}

// LoadHistoryConfig TODO: add description
func LoadHistoryConfig(name string, path string) (HistoryConfig, error) {
	var cfg HistoryConfig

	err := LoadConfig(name, path, &cfg)
	if err != nil {
		return HistoryConfig{}, err
	}

	return cfg, nil
}
