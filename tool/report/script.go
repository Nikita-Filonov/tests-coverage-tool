package report

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"tests-coverage-tool/tool/config"
	"tests-coverage-tool/tool/utils"
)

func getStateJSON(conf config.Config, coverages []coverageoutput.LogicalServiceCoverage) (string, error) {
	log.Println("Starting to make HTML report state")

	state := map[string]any{
		config.StateConfigJSON:                  conf,
		config.StateLogicalServicesCoverageJSON: coverages,
	}
	data, err := json.Marshal(state)
	if err != nil {
		return "", err
	}

	log.Printf("Successfully made HTML report state, content length: %d", len(data))

	return string(data), nil
}

func getIndexHTMLFileWithState(state string) (string, error) {
	log.Println("Starting to make HTML report file with state")

	file, err := os.Open("./submodules/tests-coverage-report/build/index.html")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, file); err != nil {
		return "", err
	}
	html := buf.String()

	scriptRegex := regexp.MustCompile(`<script id="state" type="application/json">[\s\S]*?<\/script>`)
	scriptTag := fmt.Sprintf(`<script id="state" type="application/json">%s</script>`, state)

	log.Printf("Successfully made HTML report file with state")

	return scriptRegex.ReplaceAllString(html, scriptTag), nil
}

func SaveHTMLReport(config config.Config, coverages []coverageoutput.LogicalServiceCoverage) error {
	if config.HTMLReportDir == "" {
		log.Println("Env variable 'TESTS_COVERAGE_HTML_REPORT_DIR' empty, skipping")
		return nil
	}

	state, err := getStateJSON(config, coverages)
	if err != nil {
		return err
	}

	html, err := getIndexHTMLFileWithState(state)
	if err != nil {
		return err
	}

	return utils.SaveFile([]byte(html), config.HTMLReportDir, "index.html")
}
