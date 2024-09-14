package coverageinupt

import (
	"log"
	"os"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverage"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"

	"github.com/samber/lo"
)

type InputCoverageClient struct {
	results []coverage.Result
}

type ResultsFilters struct {
	FilterByFullMethod     string
	FilterByLogicalService string
}

func (f ResultsFilters) getFilter(result coverage.Result) bool {
	if f.FilterByFullMethod != "" {
		return result.Method == f.FilterByFullMethod
	}

	if f.FilterByLogicalService != "" {
		return result.GetLogicalService() == f.FilterByLogicalService
	}

	return false
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

	var results []coverage.Result
	for _, file := range files {
		result, err := utils.ReadJSONFile[coverage.Result]("%s/%s", resultsDir, file.Name())
		if err != nil {
			log.Printf("Error reading input coverage result from file %s: %v", file.Name(), err)
			continue
		}

		results = append(results, *result)
	}

	return &InputCoverageClient{results: results}, nil
}

func (c InputCoverageClient) FilterResults(filters ResultsFilters) []coverage.Result {
	return lo.Filter(c.results, func(item coverage.Result, _ int) bool { return filters.getFilter(item) })
}

func (c InputCoverageClient) GetMethods(filters ResultsFilters) []string {
	results := c.FilterResults(filters)
	return lo.Map(results, func(item coverage.Result, _ int) string { return item.Method })
}

func (c InputCoverageClient) GetUniqueMethods(filters ResultsFilters) []string {
	return lo.Uniq(c.GetMethods(filters))
}

func (c InputCoverageClient) GetRequestParameters(filters ResultsFilters) [][]coverage.ResultParameters {
	results := c.FilterResults(filters)
	return lo.Map(results, func(item coverage.Result, _ int) []coverage.ResultParameters { return item.Request })
}

func (c InputCoverageClient) GetResponseParameters(filters ResultsFilters) [][]coverage.ResultParameters {
	results := c.FilterResults(filters)
	return lo.Map(results, func(item coverage.Result, _ int) []coverage.ResultParameters { return item.Response })
}

func (c InputCoverageClient) GetMergedRequestParameters(filters ResultsFilters) []coverage.ResultParameters {
	return coverage.MergeFilteredResultParameters(c.GetRequestParameters(filters))
}

func (c InputCoverageClient) GetMergedResponseParameters(filters ResultsFilters) []coverage.ResultParameters {
	return coverage.MergeFilteredResultParameters(c.GetResponseParameters(filters))
}
