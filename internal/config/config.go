package config

import (
	"example.com/pkg/dotEnvPackage"
	"gopkg.in/yaml.v3"
	"log"
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
