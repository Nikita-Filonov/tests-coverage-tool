package report

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/report"
)

func CopyReport() {
	if err := report.CopyHTMLReport(); err != nil {
		log.Fatalf("Error coping report file: %v", err)
	}
}

func NewCopyReportCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "copy-report",
		Short: "Copies a report",
		Run:   func(_ *cobra.Command, _ []string) { CopyReport() },
	}
}
