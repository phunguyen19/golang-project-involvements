package config_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/phunguyen19/golang-project-involvements/internal/config"
	"github.com/phunguyen19/golang-project-involvements/internal/logger"
)

func TestLoadConfig_Defaults(t *testing.T) {
	viper.Reset()
	os.Clearenv()

	// Create log that disable output for test
	log := logger.NewLogger()
	log.SetOutput(io.Discard)

	// Create a new FlagSet for this test
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	cfg, err := config.LoadConfig(log, flags, []string{})
	assert.NoError(t, err)
	assert.Equal(t, 1*time.Second, cfg.TickInterval)
	assert.Equal(t, "Hello, world!", cfg.Message)
	assert.Equal(t, "2112", cfg.MetricsPort)
	assert.Equal(t, "8080", cfg.HealthPort)
	assert.Equal(t, "", cfg.ConfigFile)
}

func TestLoadConfig_EnvOverrides(t *testing.T) {
	os.Setenv("RX9PN_TICK_INTERVAL", "2s")
	os.Setenv("RX9PN_MESSAGE", "Env Message")
	os.Setenv("RX9PN_METRICS_PORT", "9090")
	os.Setenv("RX9PN_HEALTH_PORT", "8181")
	defer os.Unsetenv("RX9PN_TICK_INTERVAL")
	defer os.Unsetenv("RX9PN_MESSAGE")
	defer os.Unsetenv("RX9PN_METRICS_PORT")
	defer os.Unsetenv("RX9PN_HEALTH_PORT")

	// Create log that disable output for test
	log := logger.NewLogger()
	log.SetOutput(io.Discard)

	// Create a new FlagSet for this test
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)

	cfg, err := config.LoadConfig(log, flags, []string{})
	assert.NoError(t, err)
	assert.Equal(t, 2*time.Second, cfg.TickInterval)
	assert.Equal(t, "Env Message", cfg.Message)
	assert.Equal(t, "9090", cfg.MetricsPort)
	assert.Equal(t, "8181", cfg.HealthPort)
	assert.Equal(t, "", cfg.ConfigFile)
}

func TestLoadConfig_ConfigFile(t *testing.T) {
	// Create a temporary config file
	configContent := `
tick_interval: 3s
message: "Config File Message"
metrics_port: "9191"
health_port: "8282"
`
	tmpFile, err := os.CreateTemp("", "config*.yaml")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(configContent))
	assert.NoError(t, err)
	tmpFile.Close()

	// Create log that disable output for test
	log := logger.NewLogger()
	log.SetOutput(io.Discard)

	// Create a new FlagSet for this test, including the config file argument
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	args := []string{"--config", tmpFile.Name()}

	cfg, err := config.LoadConfig(log, flags, args)
	assert.NoError(t, err)
	assert.Equal(t, 3*time.Second, cfg.TickInterval)
	assert.Equal(t, "Config File Message", cfg.Message)
	assert.Equal(t, "9191", cfg.MetricsPort)
	assert.Equal(t, "8282", cfg.HealthPort)
	assert.Equal(t, tmpFile.Name(), cfg.ConfigFile)
}
