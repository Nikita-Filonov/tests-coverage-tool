package config

import (
	"github.com/spf13/cobra"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/logger"
)

func PrintConfig() {
	toolConfig, err := config.NewConfig()
	if err != nil {
		logger.FatalBuildingNewClient("config", err)
	}

	toolConfig.PrintConfig()
}

func NewPrintConfigCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "print-config",
		Short: "Prints config",
		Run:   func(_ *cobra.Command, _ []string) { PrintConfig() },
	}
}
