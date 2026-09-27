package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/logger"
)

func PrintConfig() {
	if err := printConfig(); err != nil {
		logger.FatalBuildingNewClient("config", err)
	}
}

func printConfig() error {
	toolConfig, err := config.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to build config: %w", err)
	}

	toolConfig.PrintConfig()
	return nil
}

func NewPrintConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "print-config",
		Short: "Prints config",
		RunE:  func(_ *cobra.Command, _ []string) error { return printConfig() },
	}
}
