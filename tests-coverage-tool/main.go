package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Nikita-Filonov/tests-coverage-tool/tests-coverage-tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tests-coverage-tool/report"
)

func newRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{Use: "tests-coverage-tool", SilenceUsage: true, SilenceErrors: true}

	rootCmd.AddCommand(report.NewSaveReportCommand())
	rootCmd.AddCommand(report.NewCopyReportCommand())
	rootCmd.AddCommand(config.NewPrintConfigCommand())

	return rootCmd
}

func main() {
	if err := newRootCommand().Execute(); err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}
}
