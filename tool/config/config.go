package config

import (
	"errors"
	"fmt"
	"log"

	"github.com/caarlos0/env/v8"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
)

var (
	ServiceKeysShouldNotContainEmptyValuesError   = errors.New("service keys should not contain empty values")
	DuplicateServiceKeysFoundInConfigurationError = errors.New("duplicate service keys found in configuration")
)

type ServiceKey string

type Service struct {
	Key        ServiceKey `json:"key" yaml:"key"`
	Name       string     `json:"name" yaml:"name"`
	Host       string     `json:"host" yaml:"host"`
	Tags       []string   `json:"tags" yaml:"tags,omitempty"`
	Repository string     `json:"repository" yaml:"repository"`
}

type Config struct {
	Services              []Service `json:"services" yaml:"services"`
	ConfigFile            string    `env:"TESTS_COVERAGE_CONFIG_FILE" json:"-"`
	ResultsDir            string    `env:"TESTS_COVERAGE_RESULTS_DIR" envDefault:"." json:"-" yaml:"resultsDir"`
	HistoryDir            string    `env:"TESTS_COVERAGE_HISTORY_DIR" envDefault:"." json:"-" yaml:"historyDir"`
	HistoryFile           string    `env:"TESTS_COVERAGE_HISTORY_FILE" envDefault:"coverage-history.json" json:"-" yaml:"historyFile"`
	HTMLReportDir         string    `env:"TESTS_COVERAGE_HTML_REPORT_DIR" envDefault:"." json:"-" yaml:"htmlReportDir"`
	JSONReportDir         string    `env:"TESTS_COVERAGE_JSON_REPORT_DIR" envDefault:"." json:"-" yaml:"jsonReportDir"`
	HTMLReportFile        string    `env:"TESTS_COVERAGE_HTML_REPORT_FILE" envDefault:"index.html" json:"-" yaml:"htmlReportFile"`
	JSONReportFile        string    `env:"TESTS_COVERAGE_JSON_REPORT_FILE" envDefault:"coverage-report.json" json:"-" yaml:"jsonReportFile"`
	HistoryRetentionLimit int       `env:"TESTS_COVERAGE_HISTORY_RETENTION_LIMIT" envDefault:"30" json:"-" yaml:"historyRetentionLimit"`
}

func NewConfig() (Config, error) {
	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	if cfg.ConfigFile != "" {
		cfgBytes, err := utils.ReadFile(cfg.ConfigFile)
		if err != nil {
			return Config{}, err
		}

		if err = yaml.Unmarshal(cfgBytes, &cfg); err != nil {
			return Config{}, err
		}
	}

	if err := cfg.validate(); err != nil {
		log.Fatalf("Error building config: %v", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	serviceKeys := lo.Map(c.Services, func(s Service, _ int) ServiceKey { return s.Key })

	if lo.Contains(serviceKeys, "") {
		return ServiceKeysShouldNotContainEmptyValuesError
	}

	if len(c.Services) != len(lo.Uniq(serviceKeys)) {
		return DuplicateServiceKeysFoundInConfigurationError
	}

	return nil
}

func (c *Config) GetResultsDir() string {
	return fmt.Sprintf("%s/coverage-results", c.ResultsDir)
}

func (c *Config) GetHistoryFile() string {
	return fmt.Sprintf("%s/%s", c.HistoryDir, c.HistoryFile)
}

func (c *Config) GetJSONReportFile() string {
	return fmt.Sprintf("%s/%s", c.JSONReportDir, c.JSONReportFile)
}

func (c *Config) PrintConfig() {
	log.Printf("Tests coverage config: %+v\n", c)
}
