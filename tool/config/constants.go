package config

type Variable string

const (
	HTMLReportDir  Variable = "TESTS_COVERAGE_HTML_REPORT_DIR"
	JSONReportDir  Variable = "TESTS_COVERAGE_JSON_REPORT_DIR"
	JSONReportFile Variable = "TESTS_COVERAGE_JSON_REPORT_FILE"
)

func (v Variable) String() string {
	return string(v)
}
