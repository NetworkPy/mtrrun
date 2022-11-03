package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// AgentConfig configuration for agent
type AgentConfig struct {
	Host                 string        `yaml:"host"`
	Timeout              time.Duration `yaml:"timeout"`
	MaxIdleConns         int           `yaml:"maxIdleConns"`
	MaxRequestsPerMoment int           `yaml:"maxRequestsPerMoment"`
	ReportInterval       time.Duration `yaml:"reportInterval"`
	PollInterval         time.Duration `yaml:"pollInterval"`
}

// ReadAgentConfig read file with configuration and load it
func ReadAgentConfig(path string) (*AgentConfig, error) {
	c := &AgentConfig{}

	b, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(b, c)

	if err != nil {
		return nil, err
	}

	return c, nil
}
