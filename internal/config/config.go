package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

const (
	envPrefix = "AMAREH"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
	instance *Config
	once     sync.Once // Ensure thread-safe singleton initialization
)

// Config holds all configuration settings
type Config struct {
	Telegram struct {
		BotToken string `yaml:"bot_token" split_words:"true" required:"true"`
	} `yaml:"telegram"`
	Database struct {
		Host     string `yaml:"host" split_words:"true" required:"true" default:"localhost"`
		Port     int    `yaml:"port" split_words:"true" required:"true" validate:"max=65535,numeric"`
		Username string `yaml:"username" split_words:"true" required:"true"`
		Password string `yaml:"password" split_words:"true"`
		Name     string `yaml:"name" split_words:"true" required:"true"`
	} `yaml:"database"`
	LogDirectory string `yaml:"log_directory" split_words:"true" required:"true"`
}

// String writes for safe printing (make sensitive data redacted)
func (c *Config) String() string {
	return fmt.Sprintf("Telegram: [REDACTED], Database: %s, LogDirectory: %s", c.Database.Host, c.LogDirectory)
}

// LoadConfig loads configuration from yaml files and environment variables
func LoadConfig(path string) (*Config, error) {
	var loadErr error
	once.Do(func() {
		instance = new(Config)
		if err := loadFromFile(path, instance); err != nil {
			// Try environment variables if file loading fails
			if envErr := envconfig.Process(envPrefix, instance); envErr != nil {
				loadErr = fmt.Errorf("failed to load config from both file (%v) and environment (%v)", err, envErr)
				return
			}
		}

		if err := validate.Struct(instance); err != nil {
			loadErr = fmt.Errorf("config validation failed: %w", err)
			instance = nil // Reset instance on validation failure
		}
	})

	if loadErr != nil {
		return nil, loadErr
	}
	return instance, nil
}

// loadFromFile loads configuration from yaml file(s)
func loadFromFile(path string, cfg *Config) error {
	if path == "" {
		return errors.New("no config path provided")
	}

	stat, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat path: %w", err)
	}

	var yamlFile []byte
	if stat.IsDir() {
		yamlFile, err = loadFromDirectory(path)
	} else {
		yamlFile, err = loadFromSingleFile(path)
	}
	if err != nil {
		return err
	}

	if yamlFile == nil {
		return errors.New("no valid yaml configuration file found")
	}

	return yaml.Unmarshal(yamlFile, cfg)
}

// loadFromDirectory loads configuration from yaml files in directory
func loadFromDirectory(path string) ([]byte, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" && entry.Name() == "config.example.yml" {
			continue
		}

		return os.ReadFile(filepath.Join(path, entry.Name()))
	}

	return nil, errors.New("no yaml files found in directory")
}

// loadFromSingleFile loads configuration from a single yaml file
func loadFromSingleFile(path string) ([]byte, error) {
	ext := filepath.Ext(path)
	if ext != ".yaml" && ext != ".yml" {
		return nil, fmt.Errorf("file %s is not a yaml file", path)
	}
	return os.ReadFile(path)
}

// GetConfig returns the current configuration instance
func GetConfig() *Config {
	return instance
}
