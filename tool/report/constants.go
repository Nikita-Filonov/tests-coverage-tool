package report

type StateKey string

const (
	sourceReportFile      = "./submodules/tests-coverage-report/build/index.html"
	destinationReportFile = "./tool/report/templates/index.html"
)

const (
	configKey                  StateKey = "config"
	createdAtKey               StateKey = "createdAt"
	serviceCoveragesKey        StateKey = "serviceCoverages"
	logicalServiceCoveragesKey StateKey = "logicalServiceCoverages"
)
