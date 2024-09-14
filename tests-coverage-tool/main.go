package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Nikita-Filonov/tests-coverage-tool/tests-coverage-tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tests-coverage-tool/report"
)

func main() {
	var rootCmd = &cobra.Command{Use: "tests-coverage-tool"}

	rootCmd.AddCommand(report.NewSaveReportCommand())
	rootCmd.AddCommand(report.NewCopyReportCommand())
	rootCmd.AddCommand(config.NewPrintConfigCommand())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to run command: %v", err)
	}
}
