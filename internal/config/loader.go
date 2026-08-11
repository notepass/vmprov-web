package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	DefaultPort    = 8080
	DefaultLogLevel = "INFO"
)

// Load reads configuration from config.yaml and environment variables.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("server_port", DefaultPort)
	viper.SetDefault("log_level", DefaultLogLevel)

	viper.BindEnv("server_port", "SERVER_PORT")
	viper.BindEnv("db_conn_string", "DB_CONN_STRING")
	viper.BindEnv("db_username", "DB_USERNAME")
	viper.BindEnv("db_password", "DB_PASSWORD")
	viper.BindEnv("log_level", "LOG_LEVEL")

	if _, err := os.Stat("config.yaml"); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &Config{
		ServerPort:   viper.GetInt("server_port"),
		DBConnString: viper.GetString("db_conn_string"),
		DBUsername:   viper.GetString("db_username"),
		DBPassword:   viper.GetString("db_password"),
		LogLevel:     viper.GetString("log_level"),
	}, nil
}
