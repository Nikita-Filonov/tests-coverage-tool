package coverageoutput

import (
	"fmt"

	"github.com/jhump/protoreflect/desc"
	"github.com/samber/lo"

	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverage"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/coverageinupt"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/history"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/models"
	"github.com/Nikita-Filonov/tests-coverage-tool/tool/reflection"
)

type OutputCoverageClient struct {
	reflectionClient    reflection.GRPCReflectionClient
	inputHistoryClient  history.InputHistoryClient
	inputCoverageClient coverageinupt.InputCoverageClient
}

func NewOutputCoverageClient(
	reflectionClient *reflection.GRPCReflectionClient,
	inputHistoryClient *history.InputHistoryClient,
	inputCoverageClient *coverageinupt.InputCoverageClient,
) (*OutputCoverageClient, error) {
	if reflectionClient == nil {
		return nil, fmt.Errorf("reflection client is nil")
	}

	if inputHistoryClient == nil {
		return nil, fmt.Errorf("input history client is nil")
	}

	if inputCoverageClient == nil {
		return nil, fmt.Errorf("input coverage client is nil")
	}

	return &OutputCoverageClient{
		reflectionClient:    *reflectionClient,
		inputHistoryClient:  *inputHistoryClient,
		inputCoverageClient: *inputCoverageClient,
	}, nil
}

func (c *OutputCoverageClient) getMethodCoverage(descriptor *desc.MethodDescriptor) models.MethodCoverage {
	var (
		option                     = descriptor.GetMethodOptions()
		filters                    = coverageinupt.ResultsFilters{FilterByFullMethod: descriptor.GetFullyQualifiedName()}
		results                    = c.inputCoverageClient.FilterResults(filters)
		covered                    = len(results) > 0
		totalCases                 = len(results)
		methodName                 = descriptor.GetName()
		logicalServiceName         = descriptor.GetService().GetFullyQualifiedName()
		actualRequestParameters    = c.inputCoverageClient.GetMergedRequestParameters(filters)
		actualResponseParameters   = c.inputCoverageClient.GetMergedResponseParameters(filters)
		expectedRequestParameters  = reflection.BuildExpectedResultParameters(descriptor.GetInputType().UnwrapMessage())
		expectedResponseParameters = reflection.BuildExpectedResultParameters(descriptor.GetOutputType().UnwrapMessage())
	)

	var (
		requestTotalParameters         = coverage.GetTotalResultParameters(expectedRequestParameters)
		responseTotalParameters        = coverage.GetTotalResultParameters(expectedResponseParameters)
		requestTotalCoveredParameters  = coverage.GetTotalCoveredResultParameters(actualRequestParameters)
		responseTotalCoveredParameters = coverage.GetTotalCoveredResultParameters(actualResponseParameters)
		requestTotalCoverage           = getRequestCoveragePercent(covered, requestTotalParameters, requestTotalCoveredParameters)
		responseTotalCoverage          = getRequestCoveragePercent(covered, responseTotalParameters, responseTotalCoveredParameters)
		requestTotalCoverageHistory    = c.inputHistoryClient.BuildMethodHistoryRequestTotalCoverage(logicalServiceName, methodName, requestTotalCoverage)
		responseTotalCoverageHistory   = c.inputHistoryClient.BuildMethodHistoryResponseTotalCoverage(logicalServiceName, methodName, responseTotalCoverage)
		requestParametersCoverage      []models.ResultParameters
		responseParametersCoverage     []models.ResultParameters
	)
	if covered {
		requestParametersCoverage = coverage.MergeResultParameters(expectedRequestParameters, actualRequestParameters)
		responseParametersCoverage = coverage.MergeResultParameters(expectedResponseParameters, actualResponseParameters)
	}

	coverage.EnrichSliceWithUncoveredResultParameters(requestParametersCoverage)
	coverage.EnrichSliceWithUncoveredResultParameters(responseParametersCoverage)

	return models.MethodCoverage{
		Method:     descriptor.GetName(),
		Covered:    covered,
		TotalCases: totalCases,
		Deprecated: option.GetDeprecated(),
		RequestCoverage: models.MethodRequestCoverage{
			Name:                   descriptor.GetInputType().GetFullyQualifiedName(),
			TotalCoverage:          requestTotalCoverage,
			TotalParameters:        requestTotalParameters,
			ParametersCoverage:     requestParametersCoverage,
			TotalCoverageHistory:   requestTotalCoverageHistory,
			TotalCoveredParameters: requestTotalCoveredParameters,
		},
		ResponseCoverage: models.MethodRequestCoverage{
			Name:                   descriptor.GetOutputType().GetFullyQualifiedName(),
			TotalCoverage:          responseTotalCoverage,
			TotalParameters:        responseTotalParameters,
			ParametersCoverage:     responseParametersCoverage,
			TotalCoverageHistory:   responseTotalCoverageHistory,
			TotalCoveredParameters: responseTotalCoveredParameters,
		},
	}
}

func (c *OutputCoverageClient) getLogicalServiceCoverage(logicalService string) (models.LogicalServiceCoverage, error) {
	filters := coverageinupt.ResultsFilters{FilterByLogicalService: logicalService}
	resultsMethods := c.inputCoverageClient.GetUniqueMethods(filters)

	reflectionMethodsDescriptors, err := c.reflectionClient.GetServiceMethodsDescriptors(logicalService)
	if err != nil {
		return models.LogicalServiceCoverage{}, err
	}

	reflectionMethods, err := c.reflectionClient.GetServiceMethods(logicalService)
	if err != nil {
		return models.LogicalServiceCoverage{}, err
	}

	totalCoverage := getCoveragePercent(reflectionMethods, resultsMethods)
	totalCoverageHistory := c.inputHistoryClient.BuildLogicalServiceHistoryTotalCoverage(logicalService, totalCoverage)

	return models.LogicalServiceCoverage{
		Methods:              lo.Map(reflectionMethodsDescriptors, func(item *desc.MethodDescriptor, _ int) models.MethodCoverage { return c.getMethodCoverage(item) }),
		TotalMethods:         len(reflectionMethods),
		TotalCoverage:        totalCoverage,
		LogicalService:       logicalService,
		TotalCoveredMethods:  len(resultsMethods),
		TotalCoverageHistory: totalCoverageHistory,
	}, nil
}

func (c *OutputCoverageClient) GetLogicalServiceCoverages() ([]models.LogicalServiceCoverage, error) {
	logicalServices, err := c.reflectionClient.GetServices()
	if err != nil {
		return nil, err
	}

	coverages := make([]models.LogicalServiceCoverage, len(logicalServices))
	for index, logicalService := range logicalServices {
		if serviceCoverage, err := c.getLogicalServiceCoverage(logicalService); err == nil {
			coverages[index] = serviceCoverage
		}
	}

	return coverages, nil
}

func (c *OutputCoverageClient) GetServiceCoverage() (models.ServiceCoverage, error) {
	logicalServices, err := c.reflectionClient.GetServices()
	if err != nil {
		return models.ServiceCoverage{}, err
	}

	var (
		resultsMethods    []string
		reflectionMethods []string
	)
	for _, logicalService := range logicalServices {
		filters := coverageinupt.ResultsFilters{FilterByLogicalService: logicalService}
		resultsMethods = append(resultsMethods, c.inputCoverageClient.GetUniqueMethods(filters)...)

		methods, err := c.reflectionClient.GetServiceMethods(logicalService)
		if err != nil {
			return models.ServiceCoverage{}, err
		}

		reflectionMethods = append(reflectionMethods, methods...)
	}

	totalCoverage := getCoveragePercent(reflectionMethods, resultsMethods)
	totalCoverageHistory := c.inputHistoryClient.BuildServiceHistoryTotalCoverage(totalCoverage)

	return models.ServiceCoverage{
		TotalCoverage:        totalCoverage,
		TotalCoverageHistory: totalCoverageHistory,
	}, nil
}
