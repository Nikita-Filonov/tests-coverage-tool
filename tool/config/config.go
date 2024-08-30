package config

import (
	"fmt"
	"log"

	"tests-coverage-tool/tool/utils"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	Service          GRPCClientConfig `json:"service"`
	HTMLReportDir    string           `env:"TESTS_COVERAGE_HTML_REPORT_DIR"`
	InputResultsDir  string           `env:"TESTS_COVERAGE_INPUT_RESULTS_DIR"`
	OutputResultsDir string           `env:"TESTS_COVERAGE_OUTPUT_RESULTS_DIR"`
}

type GRPCClientConfig struct {
	Host string `json:"host" env:"TESTS_COVERAGE_TOOL_SERVICE_HOST"`
	Port int    `json:"port" env:"TESTS_COVERAGE_TOOL_SERVICE_PORT"`
}

func (c GRPCClientConfig) GetURL() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func NewConfig() (Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func (c *Config) SaveResults() error {
	log.Printf("Starting to save config into results folder")

	if c.OutputResultsDir == "" {
		log.Println("Env variable 'TESTS_COVERAGE_INPUT_RESULTS_DIR' empty, skipping")
		return nil
	}

	if err := utils.SaveJSONFile(c, c.OutputResultsDir, StateConfigJSON); err != nil {
		return err
	}

	return nil
}
