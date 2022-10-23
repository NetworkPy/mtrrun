package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// AgentConfig configuration for agent
type AgentConfig struct {
	Host                 string `yaml:"host"`
	Timeout              int    `yaml:"timeout"`
	MaxIdleConns         int    `yaml:"maxIdleConns"`
	MaxRequestsPerMoment int    `yaml:"maxRequestsPerMoment"`
	ReportInterval       int    `yaml:"reportInterval"`
	PollInterval         int    `yaml:"pollInterval"`
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
