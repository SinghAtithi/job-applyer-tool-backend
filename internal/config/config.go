package config

import (
	"errors"
	"example.com/pkg/dotEnvPackage"
	"gopkg.in/yaml.v3"
	"log"
	"math/rand"
	"os"
	"path/filepath"
)

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

func NewApplicationConfig() *ApplicationConfig {
	config := &ApplicationConfig{}
	configPath := getConfigPath()

	readYamlFile(configPath, config)
	return config
}

func getConfigPath() string {

	err := publicPackage.LoadDotEnvFile(".env")
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		return os.Getenv("CONFIG_PATH")
	}

	// Default to relative config path from project root
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ Failed to get working directory: %v", err)
	}
	configMode := os.Getenv("CONFIG_MODE")
	if configMode == "development" {
		// Default to development config path
		return filepath.Join(workingDir, "configs", "development.yaml")
	}

	return filepath.Join(workingDir, "configs", "production.yaml")
}

func getConfigFilePath() string {
	// Check if CONFIG_PATH is set in environment variables
	if configPath := os.Getenv("CONFIG_PATH"); configPath != "" {
		return configPath
	}

	// Default to relative path from project root
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("❌ Failed to get working directory: %v", err)
	}
	return filepath.Join(workingDir, "configs", "production.yaml")
}

func readYamlFile(filePath string, target interface{}) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("❌ Failed to read YAML file: %v", err)
	}

	err = yaml.Unmarshal(data, target)
}

func (c *ApplicationConfig) GetPort() string {
	return c.Port
}

// LoadAgentClient loads Groq configuration securely by merging YAML config and environment variables.
// Environment variables take precedence for secrets like API keys.
func LoadAgentClient() (*AgentClient, error) {
	models := GetAllAIModelNames()
	// Load YAML config
	tmp := struct {
		BaseURL string `yaml:"agent.base_url"`
		Model   string `yaml:"agent.model"`
	}{
		BaseURL: "https://api.groq.com/openai/v1/chat/completions", // default
		Model:   models[rand.Intn(len(models))],                    // default
	}
	readYamlFile(getConfigFilePath(), &tmp)

	// Load API key from environment
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		return nil, errors.New("GROQ_API_KEY not set in environment")
	}

	return &AgentClient{
		APIKey:  apiKey,
		BaseURL: tmp.BaseURL,
		Model:   tmp.Model,
	}, nil
}
