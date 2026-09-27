package report

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/report"
)

func CopyReport() {
	if err := copyReport(); err != nil {
		log.Fatal(err)
	}
}

func copyReport() error {
	if err := report.CopyHTMLReport(); err != nil {
		return fmt.Errorf("failed to copy report: %w", err)
	}
	return nil
}

func NewCopyReportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "copy-report",
		Short: "Copies a report",
		RunE:  func(_ *cobra.Command, _ []string) error { return copyReport() },
	}
}
