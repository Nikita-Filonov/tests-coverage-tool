package config

type Variable string

const (
	HistoryDir     Variable = "TESTS_COVERAGE_HISTORY_DIR"
	HistoryFile    Variable = "TESTS_COVERAGE_HISTORY_FILE"
	HTMLReportDir  Variable = "TESTS_COVERAGE_HTML_REPORT_DIR"
	JSONReportDir  Variable = "TESTS_COVERAGE_JSON_REPORT_DIR"
	HTMLReportFile Variable = "TESTS_COVERAGE_HTML_REPORT_FILE"
	JSONReportFile Variable = "TESTS_COVERAGE_JSON_REPORT_FILE"
)

func (v Variable) String() string {
	return string(v)
}
