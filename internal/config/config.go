package config

import (
	"github.com/go-kit/log/level"
	"github.com/spf13/viper"

	kitlog "github.com/go-kit/log"
)

type Config struct {
	ListenAddress string    `mapstructure:"listen-address"`
	AutofaqUrl    string    `mapstructure:"autofaq-url"`
	LogLevel      string    `mapstructure:"log-level"`
	Services      []Service `mapstructure:"services"`
	Severity      string    `mapstructure:"severity"`
}

type Service struct {
	Id      string   `mapstructure:"id"`
	Widgets []string `mapstructure:"widgets"`
}

func InitConfig(configPath *string, logger kitlog.Logger) (*Config, error) {
	var config *Config

	// Get env
	viper.AutomaticEnv()

	// Read config file
	if *configPath != "" {
		viper.SetConfigFile(*configPath)
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				level.Error(logger).Log("msg", err.Error())
			} else {
				level.Error(logger).Log("msg", err.Error())
			}
		}
	}
	err := viper.Unmarshal(&config)

	// Set config values from env
	viper.BindEnv("listen-address", "LISTEN_ADDRESS")
	viper.BindEnv("log-level", "LOG_LEVEL")
	viper.BindEnv("autofaq-url", "AUTOFAQ_URL")

	return config, err
}
