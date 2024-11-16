package history

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/config"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
)

type inputHistoryClientTest[T any] struct {
	name           string
	want           T
	method         string
	client         InputHistoryClient
	history        []models.CoverageHistory
	totalCoverage  float64
	logicalService string
}

type inputHistoryClientFactoryTest[T any] struct {
	key     config.ServiceKey
	name    string
	want    T
	factory InputHistoryClientFactory
}

var (
	createdAt       = time.Now()
	pastCreatedAt   = createdAt.AddDate(0, -1, 0)
	futureCreatedAt = createdAt.AddDate(0, 1, 0)
)

var defaultConfig = config.Config{HistoryRetentionLimit: 3}

var (
	singleHistory          = models.CoverageHistory{CreatedAt: createdAt, TotalCoverage: 10.11}
	defaultCoverageHistory = []models.CoverageHistory{
		{CreatedAt: createdAt, TotalCoverage: 99.99},
		{CreatedAt: createdAt, TotalCoverage: 77.77},
	}
)

var (
	emptyInputHistoryClient = InputHistoryClient{
		config:    defaultConfig,
		createdAt: createdAt,
	}
	defaultInputHistoryClient = InputHistoryClient{
		config: defaultConfig,
		history: models.History{
			Service: models.ServiceHistory{TotalCoverage: defaultCoverageHistory},
			LogicalServices: map[string]models.LogicalServiceHistory{
				"user-service": {
					Methods: map[string]models.MethodHistory{
						"GetUser": {
							RequestTotalCoverage:  defaultCoverageHistory,
							ResponseTotalCoverage: defaultCoverageHistory,
						},
						"GetUsers": {
							RequestTotalCoverage:  defaultCoverageHistory,
							ResponseTotalCoverage: defaultCoverageHistory,
						},
					},
					TotalCoverage: defaultCoverageHistory,
				},
				"account-service": {
					Methods: map[string]models.MethodHistory{
						"GetAccount": {
							RequestTotalCoverage:  defaultCoverageHistory,
							ResponseTotalCoverage: defaultCoverageHistory,
						},
						"GetAccounts": {
							RequestTotalCoverage:  defaultCoverageHistory,
							ResponseTotalCoverage: defaultCoverageHistory,
						},
					},
					TotalCoverage: defaultCoverageHistory,
				},
			},
		},
		createdAt: createdAt,
	}
)

func TestInputHistoryClientBuildCoverageHistory(t *testing.T) {
	tests := []inputHistoryClientTest[models.CoverageHistory]{
		{
			name:          "Total coverage less than 100",
			want:          models.CoverageHistory{CreatedAt: createdAt, TotalCoverage: 17},
			client:        InputHistoryClient{createdAt: createdAt},
			totalCoverage: 17,
		},
		{
			name:          "Total coverage more than 100",
			want:          models.CoverageHistory{CreatedAt: createdAt, TotalCoverage: 100},
			client:        InputHistoryClient{createdAt: createdAt},
			totalCoverage: 200,
		},
		{
			name:          "Without created at",
			want:          models.CoverageHistory{TotalCoverage: 5},
			client:        InputHistoryClient{},
			totalCoverage: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.client.buildCoverageHistory(test.totalCoverage)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestInputHistoryClientAppendCoverageHistory(t *testing.T) {
	tests := []inputHistoryClientTest[[]models.CoverageHistory]{
		{
			name:    "Empty history. History disabled",
			want:    []models.CoverageHistory{},
			client:  InputHistoryClient{isHistoryDisabled: true},
			history: []models.CoverageHistory{},
		},
		{
			name:    "Existing history. History disabled",
			want:    []models.CoverageHistory{},
			client:  InputHistoryClient{isHistoryDisabled: true},
			history: []models.CoverageHistory{{}, {}, {}},
		},
		{
			name:          "Existing history. Total coverage equals to zero",
			want:          []models.CoverageHistory{{}, {}, {}},
			client:        InputHistoryClient{},
			history:       []models.CoverageHistory{{}, {}, {}},
			totalCoverage: 0,
		},
		{
			name:          "Empty history. Total coverage is more than zero",
			want:          []models.CoverageHistory{singleHistory},
			client:        InputHistoryClient{config: defaultConfig, createdAt: createdAt},
			history:       []models.CoverageHistory{},
			totalCoverage: singleHistory.TotalCoverage,
		},
		{
			name:          "Existing history. Total coverage is more than zero",
			want:          append(defaultCoverageHistory, singleHistory),
			client:        InputHistoryClient{config: defaultConfig, createdAt: createdAt},
			history:       defaultCoverageHistory,
			totalCoverage: singleHistory.TotalCoverage,
		},
		{
			name: "Existing history. Sorting",
			want: []models.CoverageHistory{
				{CreatedAt: pastCreatedAt, TotalCoverage: 11.55},
				{CreatedAt: createdAt, TotalCoverage: 50.55},
				{CreatedAt: futureCreatedAt, TotalCoverage: 10.55},
			},
			client: InputHistoryClient{config: defaultConfig, createdAt: createdAt},
			history: []models.CoverageHistory{
				{CreatedAt: futureCreatedAt, TotalCoverage: 10.55},
				{CreatedAt: pastCreatedAt, TotalCoverage: 11.55},
			},
			totalCoverage: 50.55,
		},
		{
			name: "Existing history. History out of retention limit",
			want: []models.CoverageHistory{
				{CreatedAt: createdAt, TotalCoverage: 11.55},
				{CreatedAt: createdAt, TotalCoverage: 50.55},
				{CreatedAt: createdAt, TotalCoverage: 77.55},
			},
			client: InputHistoryClient{config: defaultConfig, createdAt: createdAt},
			history: []models.CoverageHistory{
				{CreatedAt: createdAt, TotalCoverage: 10.55},
				{CreatedAt: createdAt, TotalCoverage: 11.55},
				{CreatedAt: createdAt, TotalCoverage: 50.55},
			},
			totalCoverage: 77.55,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.client.appendCoverageHistory(test.history, test.totalCoverage)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestInputHistoryClientBuildServiceHistoryTotalCoverage(t *testing.T) {
	tests := []inputHistoryClientTest[[]models.CoverageHistory]{
		{
			name:          "Empty history",
			want:          []models.CoverageHistory{singleHistory},
			client:        emptyInputHistoryClient,
			totalCoverage: singleHistory.TotalCoverage,
		},
		{
			name:          "Existing history",
			want:          append(defaultCoverageHistory, singleHistory),
			client:        defaultInputHistoryClient,
			totalCoverage: singleHistory.TotalCoverage,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.client.BuildServiceHistoryTotalCoverage(test.totalCoverage)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestInputHistoryClientBuildLogicalServiceHistoryTotalCoverage(t *testing.T) {
	tests := []inputHistoryClientTest[[]models.CoverageHistory]{
		{
			name:           "Empty history",
			want:           []models.CoverageHistory{singleHistory},
			client:         emptyInputHistoryClient,
			totalCoverage:  singleHistory.TotalCoverage,
			logicalService: "user-service",
		},
		{
			name:           "Existing history",
			want:           append(defaultCoverageHistory, singleHistory),
			client:         defaultInputHistoryClient,
			totalCoverage:  singleHistory.TotalCoverage,
			logicalService: "user-service",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.client.BuildLogicalServiceHistoryTotalCoverage(
				test.logicalService, test.totalCoverage,
			)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestInputHistoryClientBuildMethodHistoryRequestTotalCoverage(t *testing.T) {
	tests := []inputHistoryClientTest[[]models.CoverageHistory]{
		{
			name:           "Empty history",
			want:           []models.CoverageHistory{singleHistory},
			method:         "GetUser",
			client:         emptyInputHistoryClient,
			totalCoverage:  singleHistory.TotalCoverage,
			logicalService: "user-service",
		},
		{
			name:           "Existing history",
			want:           append(defaultCoverageHistory, singleHistory),
			method:         "GetAccount",
			client:         defaultInputHistoryClient,
			totalCoverage:  singleHistory.TotalCoverage,
			logicalService: "account-service",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.client.BuildMethodHistoryRequestTotalCoverage(
				test.logicalService, test.method, test.totalCoverage,
			)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestInputHistoryClientBuildMethodHistoryResponseTotalCoverage(t *testing.T) {
	tests := []inputHistoryClientTest[[]models.CoverageHistory]{
		{
			name:           "Empty history",
			want:           []models.CoverageHistory{singleHistory},
			method:         "GetUser",
			client:         emptyInputHistoryClient,
			totalCoverage:  singleHistory.TotalCoverage,
			logicalService: "user-service",
		},
		{
			name:           "Existing history",
			want:           append(defaultCoverageHistory, singleHistory),
			method:         "GetAccount",
			client:         defaultInputHistoryClient,
			totalCoverage:  singleHistory.TotalCoverage,
			logicalService: "account-service",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.client.BuildMethodHistoryResponseTotalCoverage(
				test.logicalService, test.method, test.totalCoverage,
			)
			assert.Equal(t, test.want, result)
		})
	}
}

func TestInputHistoryClientFactoryNewClient(t *testing.T) {
	tests := []inputHistoryClientFactoryTest[*InputHistoryClient]{
		{
			key:     "any",
			name:    "Empty factory",
			want:    &InputHistoryClient{},
			factory: InputHistoryClientFactory{},
		},
		{
			key:  "any",
			name: "Filled factory",
			want: &defaultInputHistoryClient,
			factory: InputHistoryClientFactory{
				state:             models.HistoryState{"any": defaultInputHistoryClient.history},
				config:            defaultConfig,
				createdAt:         createdAt,
				isHistoryDisabled: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, test.factory.NewClient(test.key))
		})
	}
}
