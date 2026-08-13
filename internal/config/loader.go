package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	DefaultPort         = 8080
	DefaultLogLevel     = "INFO"
	DefaultMaxOpenConns = 25
	DefaultMaxIdleConns = 5
	// 0 means no limit
	DefaultConnMaxLifetime = 0
)

// Load reads configuration from config.yaml and environment variables.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.SetDefault("server_port", DefaultPort)
	viper.SetDefault("log_level", DefaultLogLevel)
	viper.SetDefault("db_max_open_conns", DefaultMaxOpenConns)
	viper.SetDefault("db_max_idle_conns", DefaultMaxIdleConns)
	viper.SetDefault("db_conn_max_lifetime", DefaultConnMaxLifetime)

	viper.BindEnv("server_port", "SERVER_PORT")
	viper.BindEnv("db_conn_string", "DB_CONN_STRING")
	viper.BindEnv("db_username", "DB_USERNAME")
	viper.BindEnv("db_password", "DB_PASSWORD")
	viper.BindEnv("db_max_open_conns", "DB_MAX_OPEN_CONNS")
	viper.BindEnv("db_max_idle_conns", "DB_MAX_IDLE_CONNS")
	viper.BindEnv("db_conn_max_lifetime", "DB_CONN_MAX_LIFETIME")
	viper.BindEnv("log_level", "LOG_LEVEL")

	if _, err := os.Stat("config.yaml"); err == nil {
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	return &Config{
		ServerPort:        viper.GetInt("server_port"),
		DBConnString:      viper.GetString("db_conn_string"),
		DBUsername:        viper.GetString("db_username"),
		DBPassword:        viper.GetString("db_password"),
		DBMaxOpenConns:    viper.GetInt("db_max_open_conns"),
		DBMaxIdleConns:    viper.GetInt("db_max_idle_conns"),
		DBConnMaxLifetime: viper.GetInt("db_conn_max_lifetime"),
		LogLevel:          viper.GetString("log_level"),
	}, nil
}
