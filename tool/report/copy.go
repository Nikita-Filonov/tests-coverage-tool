package report

import (
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
)

func CopyHTMLReport() error {
	return utils.CopyFile(sourceReportFile, destinationReportFile)
}
