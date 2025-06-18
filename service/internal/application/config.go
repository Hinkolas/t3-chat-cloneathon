package application

import (
	"fmt"

	"github.com/Hinkolas/t3-chat-cloneathon/service/internal/llm"
	"github.com/spf13/viper"
)

func LoadConfig(cfgFile string) (*Config, error) {

	v := viper.NewWithOptions(viper.KeyDelimiter("|"))

	// Define default config values
	v.SetDefault("server|host", ":3141")
	v.SetDefault("logging|log_file_path", "data/app.log")
	v.SetDefault("logging|log_format", "text")
	v.SetDefault("logging|log_level", "debug")

	// Tell viper where to look for the config file
	v.SetConfigFile(cfgFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("failed to read config: %v\n", err)
		return nil, err
	}

	// Unmarshal the config into application.Config
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		fmt.Printf("failed to decode config: %v\n", err)
		return nil, err
	}

	return &cfg, nil

}

type Config struct {
	Server  ServerConfig         `mapstructure:"server" yaml:"server"`
	Logging LoggingConfig        `mapstructure:"logging" yaml:"logging"`
	Users   []UserConfig         `mapstructure:"users" yaml:"users"`
	Models  map[string]llm.Model `mapstructure:"models" yaml:"models"`
}

type ServerConfig struct {
	Host string `mapstructure:"host" yaml:"host"` // Hostname of the application
	// ReadTimeout  int    `mapstructure:"read_timeout" yaml:"read_timeout` // Time a request must take at most in seconds
	// WriteTimeout int  `mapstructure:"write_timeout" yaml:"write_timeout`  // Time a response must take at most in seconds
}

type LoggingConfig struct {
	LogFilePath string `mapstructure:"log_file_path" yaml:"log_file_path"` // Path of the log file used for structured logging
	LogLevel    string `mapstructure:"log_level" yaml:"log_level"`         // Set the logging level ("debug", "info", "warn", "error") (default "info")
	LogFormat   string `mapstructure:"log_format" yaml:"log_format"`       // Format of the structured logger (text, json) (default: json)
	LogVerbose  bool   // Output log messages to stdout in addition to the log file
}

type UserConfig struct {
	Username string `mapstructure:"username" yaml:"username"`
	Email    string `mapstructure:"email" yaml:"email"`
	Password string `mapstructure:"password" yaml:"password"`
}
