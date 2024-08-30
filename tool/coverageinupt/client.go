package coverageinupt

import (
	"os"

	"tests-coverage-tool/tool/utils"

	"github.com/samber/lo"
)

type InputCoverageClient struct {
	results []InputCoverageResult
}

type ResultsFilters struct {
	Method         string
	LogicalService string
}

func (f ResultsFilters) getFilter(result InputCoverageResult) bool {
	if f.Method == "" {
		return result.LogicalService == f.LogicalService
	}

	return result.Method == f.Method && result.LogicalService == f.LogicalService
}

func NewInputCoverageClient(resultsDir string) (*InputCoverageClient, error) {
	dir, err := os.Open(resultsDir)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	results := make([]InputCoverageResult, len(files))
	for index, file := range files {
		result, err := utils.ReadJSONFile[InputCoverageResult]("%s/%s", resultsDir, file.Name())
		if err != nil {
			continue
		}

		results[index] = *result
	}

	return &InputCoverageClient{results: results}, nil
}

func (c InputCoverageClient) FilterResults(filters ResultsFilters) []InputCoverageResult {
	return lo.Filter(c.results, func(item InputCoverageResult, _ int) bool { return filters.getFilter(item) })
}

func (c InputCoverageClient) GetMethods(filters ResultsFilters) []string {
	results := c.FilterResults(filters)
	return lo.Map(results, func(item InputCoverageResult, _ int) string { return item.Method })
}

func (c InputCoverageClient) GetUniqueMethods(filters ResultsFilters) []string {
	return lo.Uniq(c.GetMethods(filters))
}

func (c InputCoverageClient) IsMethodCoveted(filters ResultsFilters) bool {
	_, found := lo.Find(c.results, func(result InputCoverageResult) bool { return filters.getFilter(result) })
	return found
}

func (c InputCoverageClient) GetRequestParameters(filters ResultsFilters) []string {
	results := c.FilterResults(filters)
	return lo.Flatten(lo.Map(results, func(item InputCoverageResult, _ int) []string { return item.RequestParameters }))
}

func (c InputCoverageClient) GetResponseParameters(filters ResultsFilters) []string {
	results := c.FilterResults(filters)
	return lo.Flatten(lo.Map(results, func(item InputCoverageResult, _ int) []string { return item.ResponseParameters }))
}
