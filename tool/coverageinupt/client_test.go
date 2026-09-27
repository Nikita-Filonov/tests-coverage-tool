package coverageinupt

import (
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/utils"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverage"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

type inputCoverageClientTest[T any] struct {
	name    string
	want    T
	client  InputCoverageClient
	filters ResultsFilters
}

var filterResultsClient = InputCoverageClient{
	results: []models.Result{
		{Method: "service.get"},
		{Method: "service.create"},
		{Method: "service.update"},
		{Method: "service2.get"},
		{Method: "service2.get"},
	},
}

var getRequestParametersClient = InputCoverageClient{
	results: []models.Result{
		{
			Method:   "service.get",
			Request:  []models.ResultParameters{{Parameter: "a"}, {Parameter: "b"}},
			Response: []models.ResultParameters{{Parameter: "a"}, {Parameter: "b"}},
		},
		{
			Method:   "service.create",
			Request:  []models.ResultParameters{{Parameter: "c"}, {Parameter: "d"}},
			Response: []models.ResultParameters{{Parameter: "c"}, {Parameter: "d"}},
		},
	},
}

var getMergedRequestParametersClient = InputCoverageClient{
	results: []models.Result{
		{
			Method:   "service.get",
			Request:  []models.ResultParameters{{Parameter: "a"}, {Parameter: "b", Covered: true}},
			Response: []models.ResultParameters{{Parameter: "a"}, {Parameter: "b", Covered: true}},
		},
		{
			Method:   "service.get",
			Request:  []models.ResultParameters{{Parameter: "a", Covered: true}, {Parameter: "b"}},
			Response: []models.ResultParameters{{Parameter: "a", Covered: true}, {Parameter: "b"}},
		},
		{
			Method:   "service.create",
			Request:  []models.ResultParameters{{Parameter: "a"}, {Parameter: "b"}},
			Response: []models.ResultParameters{{Parameter: "a"}, {Parameter: "b"}},
		},
		{
			Method:   "service2.get",
			Request:  []models.ResultParameters{{Parameter: "a"}, {Parameter: "b"}},
			Response: []models.ResultParameters{{Parameter: "a"}, {Parameter: "b"}},
		},
	},
}

func TestInputCoverageClientFilterResults(t *testing.T) {
	tests := []inputCoverageClientTest[[]models.Result]{
		{
			name:    "Filter by full method",
			want:    []models.Result{{Method: "service.get"}},
			client:  filterResultsClient,
			filters: ResultsFilters{FilterByFullMethod: "service.get"},
		},
		{
			name: "Filter by logical service",
			want: []models.Result{
				{Method: "service.get"},
				{Method: "service.create"},
				{Method: "service.update"},
			},
			client:  filterResultsClient,
			filters: ResultsFilters{FilterByLogicalService: "service"},
		},
		{
			name:    "Empty filters",
			want:    []models.Result{},
			client:  filterResultsClient,
			filters: ResultsFilters{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.client.FilterResults(test.filters))
		})
	}
}

func TestInputCoverageClientGetMethods(t *testing.T) {
	tests := []inputCoverageClientTest[[]string]{
		{
			name:    "Filter by full method",
			want:    []string{"service.get"},
			client:  filterResultsClient,
			filters: ResultsFilters{FilterByFullMethod: "service.get"},
		},
		{
			name:    "Filter by logical service",
			want:    []string{"service.get", "service.create", "service.update"},
			client:  filterResultsClient,
			filters: ResultsFilters{FilterByLogicalService: "service"},
		},
		{
			name:    "Empty filters",
			want:    []string{},
			client:  filterResultsClient,
			filters: ResultsFilters{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.client.GetMethods(test.filters))
		})
	}
}

func TestInputCoverageClientGetUniqueMethods(t *testing.T) {
	tests := []inputCoverageClientTest[[]string]{
		{
			name:    "Filter by full method",
			want:    []string{"service.get"},
			client:  filterResultsClient,
			filters: ResultsFilters{FilterByFullMethod: "service.get"},
		},
		{
			name:    "Filter by logical service",
			want:    []string{"service2.get"},
			client:  filterResultsClient,
			filters: ResultsFilters{FilterByLogicalService: "service2"},
		},
		{
			name:    "Empty filters",
			want:    []string{},
			client:  filterResultsClient,
			filters: ResultsFilters{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.client.GetUniqueMethods(test.filters))
		})
	}
}

func TestInputCoverageClientGetRequestParameters(t *testing.T) {
	tests := []inputCoverageClientTest[[][]models.ResultParameters]{
		{
			name: "Filter by full method",
			want: [][]models.ResultParameters{
				{{Parameter: "a"}, {Parameter: "b"}},
			},
			client:  getRequestParametersClient,
			filters: ResultsFilters{FilterByFullMethod: "service.get"},
		},
		{
			name: "Filter by logical service",
			want: [][]models.ResultParameters{
				{{Parameter: "a"}, {Parameter: "b"}},
				{{Parameter: "c"}, {Parameter: "d"}},
			},
			client:  getRequestParametersClient,
			filters: ResultsFilters{FilterByLogicalService: "service"},
		},
		{
			name:    "Empty filters",
			want:    [][]models.ResultParameters{},
			client:  getRequestParametersClient,
			filters: ResultsFilters{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.client.GetRequestParameters(test.filters))
		})
	}
}

func TestInputCoverageClientGetResponseParameters(t *testing.T) {
	tests := []inputCoverageClientTest[[][]models.ResultParameters]{
		{
			name: "Filter by full method",
			want: [][]models.ResultParameters{
				{{Parameter: "a"}, {Parameter: "b"}},
			},
			client:  getRequestParametersClient,
			filters: ResultsFilters{FilterByFullMethod: "service.get"},
		},
		{
			name: "Filter by logical service",
			want: [][]models.ResultParameters{
				{{Parameter: "a"}, {Parameter: "b"}},
				{{Parameter: "c"}, {Parameter: "d"}},
			},
			client:  getRequestParametersClient,
			filters: ResultsFilters{FilterByLogicalService: "service"},
		},
		{
			name:    "Empty filters",
			want:    [][]models.ResultParameters{},
			client:  getRequestParametersClient,
			filters: ResultsFilters{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.client.GetResponseParameters(test.filters))
		})
	}
}

func TestInputCoverageClientGetMergedRequestParameters(t *testing.T) {
	tests := []inputCoverageClientTest[[]models.ResultParameters]{
		{
			name: "Filter by full method",
			want: []models.ResultParameters{
				{Covered: true, Parameter: "a", Parameters: []models.ResultParameters{}},
				{Covered: true, Parameter: "b", Parameters: []models.ResultParameters{}},
			},
			client:  getMergedRequestParametersClient,
			filters: ResultsFilters{FilterByFullMethod: "service.get"},
		},
		{
			name: "Filter by logical service",
			want: []models.ResultParameters{
				{Covered: true, Parameter: "a", Parameters: []models.ResultParameters{}},
				{Covered: true, Parameter: "b", Parameters: []models.ResultParameters{}},
			},
			client:  getMergedRequestParametersClient,
			filters: ResultsFilters{FilterByLogicalService: "service"},
		},
		{
			name:    "Empty filters",
			client:  getMergedRequestParametersClient,
			filters: ResultsFilters{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.client.GetMergedRequestParameters(test.filters)
			coverage.SortResultParameters(result)
			coverage.SortResultParameters(test.want)

			assert.Equal(t, test.want, result)
		})
	}
}

func TestInputCoverageClientGetMergedResponseParameters(t *testing.T) {
	tests := []inputCoverageClientTest[[]models.ResultParameters]{
		{
			name: "Filter by full method",
			want: []models.ResultParameters{
				{Covered: true, Parameter: "a", Parameters: []models.ResultParameters{}},
				{Covered: true, Parameter: "b", Parameters: []models.ResultParameters{}},
			},
			client:  getMergedRequestParametersClient,
			filters: ResultsFilters{FilterByFullMethod: "service.get"},
		},
		{
			name: "Filter by logical service",
			want: []models.ResultParameters{
				{Covered: true, Parameter: "a", Parameters: []models.ResultParameters{}},
				{Covered: true, Parameter: "b", Parameters: []models.ResultParameters{}},
			},
			client:  getMergedRequestParametersClient,
			filters: ResultsFilters{FilterByLogicalService: "service"},
		},
		{
			name:    "Empty filters",
			client:  getMergedRequestParametersClient,
			filters: ResultsFilters{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.client.GetMergedResponseParameters(test.filters)
			coverage.SortResultParameters(result)
			coverage.SortResultParameters(test.want)

			assert.Equal(t, test.want, result)
		})
	}
}

func TestInputCoverageClientReadsValidResults(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, utils.SaveJSONFile(models.Result{Method: "service.Check"}, dir, "valid.json"))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "broken.json"), []byte("{"), 0o600))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "directory"), 0o755))
	client, err := NewInputCoverageClient(dir)
	require.NoError(t, err)
	assert.Equal(t, []string{"service.Check"}, client.GetMethods(ResultsFilters{FilterByLogicalService: "service"}))
}

func TestInputCoverageClientReadErrors(t *testing.T) {
	_, err := NewInputCoverageClient(filepath.Join(t.TempDir(), "missing"))
	assert.ErrorIs(t, err, os.ErrNotExist)
	file := filepath.Join(t.TempDir(), "file")
	require.NoError(t, os.WriteFile(file, nil, 0o600))
	_, err = NewInputCoverageClient(file)
	assert.Error(t, err)
}
