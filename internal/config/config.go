package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Config holds the configuration for the application
type Config struct {
	TickInterval time.Duration
	Message      string
	MetricsPort  string
	HealthPort   string
	ConfigFile   string
}

// LoadConfig loads configuration using Viper
func LoadConfig(log *logrus.Logger, flags *pflag.FlagSet, args []string) (*Config, error) {
	// Define default values
	viper.SetDefault("tick_interval", "1s")
	viper.SetDefault("message", "Hello, world!")
	viper.SetDefault("metrics_port", "2112")
	viper.SetDefault("health_port", "8080")

	tickIntervalStr := viper.GetString("tick_interval")
	tickInterval, err := time.ParseDuration(tickIntervalStr)
	if err != nil {
		log.Printf("Invalid duration for tick_interval '%s', using default 5s", tickIntervalStr)
		tickInterval = 5 * time.Second
	}

	// Environment variables
	viper.SetEnvPrefix("RX9PN")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Bind command-line flags
	flags.String("config", "", "Path to configuration file")
	flags.Duration("tick_interval", tickInterval, "Tick interval for printing messages")
	flags.String("message", viper.GetString("message"), "Message to display")
	flags.String("metrics_port", viper.GetString("metrics_port"), "Port for metrics endpoint")
	flags.String("health_port", viper.GetString("health_port"), "Port for health check endpoint")

	// Parse the provided arguments
	if err := flags.Parse(args); err != nil {
		return nil, fmt.Errorf("error parsing flags: %v", err)
	}

	// Bind pflag to viper
	err = viper.BindPFlags(flags)
	if err != nil {
		return nil, fmt.Errorf("error binding flags: %v", err)
	}

	// Read config file if specified
	configFile := viper.GetString("config")
	fmt.Println(configFile)
	if configFile != "" {
		viper.SetConfigFile(configFile)
		if err := viper.ReadInConfig(); err != nil {
			return nil, fmt.Errorf("error reading config file: %v", err)
		}
	} else {
		// Try to read default config file if exists
		_ = viper.ReadInConfig()
	}

	// Populate Config struct
	cfg := &Config{
		TickInterval: viper.GetDuration("tick_interval"),
		Message:      viper.GetString("message"),
		MetricsPort:  viper.GetString("metrics_port"),
		HealthPort:   viper.GetString("health_port"),
		ConfigFile:   configFile,
	}

	return cfg, nil
}
