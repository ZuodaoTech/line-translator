package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var cfg *Config

type (
	Config struct {
		ConfigFilePath string
		DB             DatabaseConfig `yaml:"db"`
		Sys            Sys            `yaml:"sys"`
		Azure          Azure          `yaml:"azure"`
		Auth           Auth           `yaml:"auth"`
		Aws            Aws            `yaml:"aws"`
		OpenAI         OpenAI         `yaml:"openai"`
		Line           Line           `yaml:"line"`
	}

	DatabaseConfig struct {
		Driver string `yaml:"driver"`
		DSN    string `yaml:"dsn"`
	}

	Sys struct {
		Host       string `yaml:"host"`
		ApiBase    string `yaml:"api_base"`
		HashIDSalt string `yaml:"hash_id_salt"`
	}

	Azure struct {
		OpenAI struct {
			APIKey                string `yaml:"api_key"`
			Endpoint              string `yaml:"endpoint"`
			GptDeploymentID       string `yaml:"gpt_deployment_id"`
			EmbeddingDeploymentID string `yaml:"embedding_deployment_id"`
		} `yaml:"openai"`
		Speech struct {
			APIKey   string `yaml:"api_key"`
			Endpoint string `yaml:"endpoint"`
		} `yaml:"speech"`
	}

	Auth struct {
		JwtSecret string `yaml:"jwt_secret"`
	}

	Aws struct {
		Key    string `yaml:"key"`
		Secret string `yaml:"secret"`
		Region string `yaml:"region"`
	}

	OpenAI struct {
		APIKey string `yaml:"api_key"`
	}

	Line struct {
		ChannelID     string `yaml:"channel_id"`
		ChannelKey    string `yaml:"channel_key"`
		ChannelSecret string `yaml:"channel_secret"`
		JWTPublicKey  string `yaml:"jwt_public_key"`
		JWTPrivateKey string `yaml:"jwt_private_key"`
	}
)

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read YAML file: %v", err)
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("failed to unmarshal YAML: %v", err)
		return nil, err
	}

	cfg.ConfigFilePath = path
	return cfg, nil
}

func C() *Config {
	var err error
	if cfg == nil {
		cfg, err = LoadConfig("./config.yaml")
		if err != nil {
			log.Fatalf("failed to load config: %v", err)
			os.Exit(-1)
		}
	}
	return cfg
}

func (c *Config) ToYaml() (string, error) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
