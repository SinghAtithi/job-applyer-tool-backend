package config

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	publicPackage "example.com/pkg/dotEnvPackage"
	"gopkg.in/yaml.v3"
)

// ApplicationConfig represents the main application configuration
type ApplicationConfig struct {
	Port        string `yaml:"port" default:"8080"`
	Environment string `yaml:"environment" default:"production"`
	Logging     struct {
		Level string `yaml:"level" default:"info"`
	} `yaml:"logging"`
}

// AgentClient holds Groq API configuration loaded from environment variables and YAML config
// It merges secure values from .env and YAML, prioritizing environment variables for secrets.
type AgentClient struct {
	APIKey  string // API key for Groq (from env)
	BaseURL string // Base URL for Groq API (from YAML or default)
	Model   string // Model name (from YAML or default)
}

// NewApplicationConfig creates a new application configuration
func NewApplicationConfig() (*ApplicationConfig, error) {
	config := &ApplicationConfig{}
	configPath := getConfigPath()

	if err := readYamlFile(configPath, config); err != nil {
		if os.IsNotExist(err) {
			// Use defaults if file not found
			config.Port = "8080"
			config.Environment = "production"
			config.Logging.Level = "info"
		} else {
			// Propagate parse or IO errors so callers can decide
			return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
	}
	return config, nil
}

// getConfigPath determines the configuration file path based on environment
func getConfigPath() string {
	// Load .env file if exists
	_ = publicPackage.LoadDotEnvFile(".env")

	// Check for explicit config path
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		return configPath
	}

	// Use config mode or default to production
	configMode := os.Getenv("CONFIG_MODE")
	if configMode == "development" {
		return "configs/development.yaml"
	}

	return "configs/production.yaml"
}

// readYamlFile reads and unmarshals a YAML file
func readYamlFile(filePath string, target interface{}) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read YAML file: %w", err)
	}

	if err := yaml.Unmarshal(data, target); err != nil {
		return fmt.Errorf("failed to unmarshal YAML: %w", err)
	}
	return nil
}

// GetPort returns the configured port
func (c *ApplicationConfig) GetPort() string {
	return c.Port
}

// LoadAgentClient loads Groq configuration securely by merging YAML config and environment variables.
// Environment variables take precedence for secrets like API keys.
func LoadAgentClient() (*AgentClient, error) {
	models := GetAllAIModelNames()
	modelName := ""
	if len(models) > 0 {
		modelName = models[rand.Intn(len(models))]
	} else {
		// safe default if no models registered
		modelName = "gpt-default"
	}

	// Default configuration
	cfg := struct {
		BaseURL string `yaml:"agent_base_url"`
		Model   string `yaml:"agent_model"`
	}{
		BaseURL: "https://api.groq.com/openai/v1/chat/completions",
		Model:   modelName,
	}

	// Try to read from config file and propagate any non-not-found errors
	if err := readYamlFile(getConfigPath(), &cfg); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read agent config: %w", err)
		}
		// file not found -> continue with defaults
	}

	// Load API key from environment (required)
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GROQ_API_KEY not set in environment")
	}

	return &AgentClient{
		APIKey:  apiKey,
		BaseURL: cfg.BaseURL,
		Model:   cfg.Model,
	}, nil
}
