package config

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"

	"github.com/spf13/viper"
)

var (
	locationConfigKeys = []string{
		"APP_ENV",
		"DB_DRIVER",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
		"BIND_ADDR_HTTP",
		"BIND_ADDR_GRPC",
		"HISTORY_ADDR",
	}
	historyConfigKeys = []string{
		"APP_ENV",
		"DB_DRIVER",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
		"BIND_ADDR_HTTP",
		"BIND_ADDR_GRPC",
		"LOCATION_ADDR",
	}
)

// LocationConfig stores all configuration of user application
type LocationConfig struct {
	AppEnv       string `mapstructure:"APP_ENV"`
	DBDriver     string `mapstructure:"DB_DRIVER" validate:"required"`
	DBHost       string `mapstructure:"DB_HOST" validate:"required"`
	DBPort       string `mapstructure:"DB_PORT" validate:"required"`
	DBUser       string `mapstructure:"DB_USER" validate:"required"`
	DBPassword   string `mapstructure:"DB_PASSWORD" validate:"required"`
	DBName       string `mapstructure:"DB_NAME" validate:"required"`
	DBSSLMode    string `mapstructure:"DB_SSLMODE" validate:"required"`
	BindAddrHTTP string `mapstructure:"BIND_ADDR_HTTP" validate:"required"`
	BindAddrGRPC string `mapstructure:"BIND_ADDR_GRPC" validate:"required"`
	HistoryAddr  string `mapstructure:"HISTORY_ADDR" validate:"required"`
}

// HistoryConfig stores all configuration of user application
type HistoryConfig struct {
	AppEnv       string `mapstructure:"APP_ENV"`
	DBDriver     string `mapstructure:"DB_DRIVER" validate:"required"`
	DBHost       string `mapstructure:"DB_HOST" validate:"required"`
	DBPort       string `mapstructure:"DB_PORT" validate:"required"`
	DBUser       string `mapstructure:"DB_USER" validate:"required"`
	DBPassword   string `mapstructure:"DB_PASSWORD" validate:"required"`
	DBName       string `mapstructure:"DB_NAME" validate:"required"`
	DBSSLMode    string `mapstructure:"DB_SSLMODE" validate:"required"`
	BindAddrHTTP string `mapstructure:"BIND_ADDR_HTTP" validate:"required"`
	BindAddrGRPC string `mapstructure:"BIND_ADDR_GRPC" validate:"required"`
	LocationAddr string `mapstructure:"LOCATION_ADDR" validate:"required"`
}

// LoadConfig parses configuration and stores the result in
// the value pointed to by config.
func LoadConfig(v *viper.Viper, name string, path string, config interface{}) error {
	var err error
	validate := validator.New()

	if path != "" {
		v.AddConfigPath(path)
	}

	if name != "" {
		v.SetConfigName(name)
	}
	v.SetConfigType("env")

	v.AutomaticEnv()

	v.SetDefault("AppEnv", "development")

	err = v.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("config file not found")
		} else {
			log.Printf("failed to load config: %v", err)
			return err
		}
	}

	err = v.Unmarshal(config)
	if err != nil {
		return err
	}

	err = validate.Struct(config)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	return nil
}

func bindEnv(v *viper.Viper, keys []string) error {
	var err error

	for _, k := range keys {
		err = v.BindEnv(k)
		if err != nil {
			return fmt.Errorf("failed to bind keys: %v", err)
		}
	}

	return nil
}

// LoadLocationConfig TODO: add description
func LoadLocationConfig(name string, path string) (LocationConfig, error) {
	var err error
	var cfg LocationConfig
	v := viper.New()

	err = bindEnv(v, locationConfigKeys)
	if err != nil {
		return LocationConfig{}, err
	}

	err = LoadConfig(v, name, path, &cfg)
	if err != nil {
		return LocationConfig{}, err
	}

	return cfg, nil
}

// LoadHistoryConfig TODO: add description
func LoadHistoryConfig(name string, path string) (HistoryConfig, error) {
	var err error
	var cfg HistoryConfig
	v := viper.New()

	err = bindEnv(v, historyConfigKeys)
	if err != nil {
		return HistoryConfig{}, err
	}

	err = LoadConfig(v, name, path, &cfg)
	if err != nil {
		return HistoryConfig{}, err
	}

	return cfg, nil
}
