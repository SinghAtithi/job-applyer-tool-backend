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
	devMode   bool // Not persisted, runtime flag
	debugMode bool // Set via --debug flag
}

// AgentClient holds Groq API configuration loaded from environment variables and YAML config
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
			config.Port = "8080"
			config.Environment = "production"
			config.Logging.Level = "info"
		} else {
			return nil, fmt.Errorf("failed to load config from %s: %w", configPath, err)
		}
	}
	return config, nil
}

// getConfigPath determines the configuration file path based on environment
func getConfigPath() string {
	_ = publicPackage.LoadDotEnvFile(".env")

	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		return configPath
	}

	configMode := os.Getenv("CONFIG_MODE")
	devFlag := os.Getenv("DEV")

	if configMode == "development" || devFlag == "true" {
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

// GetLogLevel returns the configured logging level string
func (c *ApplicationConfig) GetLogLevel() string {
	if c.Logging.Level == "" {
		return "info"
	}
	return c.Logging.Level
}

// IsDev returns true if running in development mode
func (c *ApplicationConfig) IsDev() bool {
	if os.Getenv("DEV") == "true" {
		return true
	}
	if os.Getenv("CONFIG_MODE") == "development" {
		return true
	}
	return c.devMode
}

// SetDevMode sets the development mode flag
func (c *ApplicationConfig) SetDevMode(dev bool) {
	c.devMode = dev
}

// IsDebug returns true if debug mode is enabled
func (c *ApplicationConfig) IsDebug() bool {
	return c.debugMode || c.Logging.Level == "debug"
}

// SetDebugMode sets the debug mode flag
func (c *ApplicationConfig) SetDebugMode(debug bool) {
	c.debugMode = debug
}

// LoadAgentClient loads Groq configuration securely.
// Environment variables take precedence for secrets like API keys.
func LoadAgentClient() (*AgentClient, error) {
	models := GetAllAIModelNames()
	modelName := ""
	if len(models) > 0 {
		modelName = models[rand.Intn(len(models))]
	} else {
		modelName = "gpt-default"
	}

	cfg := struct {
		BaseURL string `yaml:"agent_base_url"`
		Model   string `yaml:"agent_model"`
	}{
		BaseURL: "https://api.groq.com/openai/v1/chat/completions",
		Model:   modelName,
	}

	if err := readYamlFile(getConfigPath(), &cfg); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to read agent config: %w", err)
		}
	}

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
