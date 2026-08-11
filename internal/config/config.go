package config

// Config holds all application configuration.
type Config struct {
	ServerPort    int
	DBConnString  string
	DBUsername    string
	DBPassword    string
	LogLevel      string
}
