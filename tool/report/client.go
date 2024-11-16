package report

import (
	_ "embed"
	"fmt"
	"regexp"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/logger"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
)

//go:embed templates/index.html
var indexHTML string

type CoverageReportClient struct {
	state  models.CoverageState
	config config.Config
}

func NewCoverageReportClient(conf config.Config, state models.CoverageState) CoverageReportClient {
	return CoverageReportClient{config: conf, state: state}
}

func (c CoverageReportClient) getIndexHTMLFileWithState() (string, error) {
	stateJSON, err := c.state.GetReportStateJSON()
	if err != nil {
		return "", err
	}

	scriptRegex := regexp.MustCompile(`<script id="state" type="application/json">[\s\S]*?<\/script>`)
	scriptTag := fmt.Sprintf(`<script id="state" type="application/json">%s</script>`, stateJSON)

	return scriptRegex.ReplaceAllString(indexHTML, scriptTag), nil
}

func (c CoverageReportClient) SaveHTMLReport() error {
	logger.StartMakeReport("HTML")

	if c.config.HTMLReportDir == "" {
		logger.EnvVariableEmptySkipping(config.HTMLReportDir.String())
		return nil
	}

	if c.config.HTMLReportFile == "" {
		logger.EnvVariableEmptySkipping(config.HTMLReportFile.String())
		return nil
	}

	html, err := c.getIndexHTMLFileWithState()
	if err != nil {
		return err
	}

	err = utils.SaveFile([]byte(html), c.config.HTMLReportDir, c.config.HTMLReportFile)
	if err != nil {
		logger.ErrorMakingReport("HTML")
		return err
	}

	logger.SuccessfullyMadeReport("HTML")
	return nil
}

func (c CoverageReportClient) SaveJSONReport() error {
	logger.StartMakeReport("JSON")

	if c.config.JSONReportDir == "" {
		logger.EnvVariableEmptySkipping(config.JSONReportDir.String())
		return nil
	}

	if c.config.JSONReportFile == "" {
		logger.EnvVariableEmptySkipping(config.JSONReportFile.String())
		return nil
	}

	err := utils.SaveJSONFile(c.state, c.config.JSONReportDir, c.config.JSONReportFile)
	if err != nil {
		logger.ErrorMakingReport("HTML")
		return err
	}

	logger.SuccessfullyMadeReport("JSON")
	return nil
}
