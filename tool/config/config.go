package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v8"
	"gopkg.in/yaml.v3"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
)

type ServiceHost string

type Service struct {
	Name       string      `json:"name" yaml:"name"`
	Host       ServiceHost `json:"host" yaml:"host"`
	Repository string      `json:"repository" yaml:"repository"`
}

type Config struct {
	Services       []Service `json:"services" yaml:"services"`
	ConfigFile     string    `env:"TESTS_COVERAGE_CONFIG_FILE" json:"-"`
	ResultsDir     string    `env:"TESTS_COVERAGE_RESULTS_DIR" envDefault:"." json:"-" yaml:"resultsDir"`
	HTMLReportDir  string    `env:"TESTS_COVERAGE_HTML_REPORT_DIR" envDefault:"." json:"-" yaml:"htmlReportDir"`
	JSONReportDir  string    `env:"TESTS_COVERAGE_JSON_REPORT_DIR" envDefault:"." json:"-" yaml:"jsonReportDir"`
	HTMLReportFile string    `env:"TESTS_COVERAGE_HTML_REPORT_FILE" envDefault:"index.html" json:"-" yaml:"htmlReportFile"`
	JSONReportFile string    `env:"TESTS_COVERAGE_JSON_REPORT_FILE" envDefault:"coverage-report.json" json:"-" yaml:"jsonReportFile"`
}

func (h ServiceHost) String() string {
	return string(h)
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

	return cfg, nil
}

func (c *Config) GetResultsDir() string {
	return fmt.Sprintf("%s/coverage-results", c.ResultsDir)
}

func (c *Config) GetJSONReportFile() string {
	return fmt.Sprintf("%s/%s", c.JSONReportDir, c.JSONReportFile)
}

func (c *Config) PrintConfig() {
	log.Printf("Tests coverage config: %+v\n", c)
}
